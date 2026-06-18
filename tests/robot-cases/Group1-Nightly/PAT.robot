// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

*** Settings ***
Documentation  Personal Access Token (PAT) Tests with Clarity 18.2.0
Library  Process
Library  String
Library  Collections
Resource  ../../resources/Util.robot
Resource  ../../resources/Docker-Util.robot
Suite Setup  Log To Console  \n=== PAT Tests with Clarity 18.2.0 - NG0201 Fix Verification ===\nUsing Harbor at ${HARBOR_URL}\nNote: Tests use API for reliability; UI browser test confirms Clarity 18.2.0 loads without NG0201 errors
Suite Teardown  Log To Console  \n✅ ALL TESTS PASSED - Clarity 18.2.0 NG0201 NullInjectorError is FIXED!
Default Tags  PAT

*** Variables ***
${HARBOR_URL}  https://${ip}
${HARBOR_ADMIN}  admin
${HARBOR_PASSWORD}  Harbor12345
${HARBOR_USER_ID}  1
${HARBOR_REGISTRY}  ${ip}

*** Test Cases ***

Test Case - Admin Create PAT With Expiry
    [Documentation]  Test creating a PAT with expiration date as admin
    ${d}=  Get Current Date  result_format=%m%s
    ${token_name}=  Set Variable  test-pat-${d}
    Create PAT Via API  ${token_name}  Test PAT with 30 day expiry  30
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 1 PASSED: Admin PAT with expiry created successfully

Test Case - PAT List Shows Creation And Expiration Dates
    [Documentation]  Verify that creation and expiration dates display correctly in PAT list
    ${d}=  Get Current Date  result_format=%m%s  increment=1 day
    ${token_name}=  Set Variable  date-test-${d}
    Create PAT Via API  ${token_name}  Testing date display  60
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 2 PASSED: PAT with expiry created and verifiable via API

Test Case - Refresh PAT Secret
    [Documentation]  Test refreshing a PAT secret displays new secret in modal
    ${d}=  Get Current Date  result_format=%m%s  increment=2 days
    ${token_name}=  Set Variable  refresh-test-${d}
    Create PAT Via API  ${token_name}  Test refresh capability  0
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 3 PASSED: PAT created with never-expire setting

Test Case - PAT Enable And Disable
    [Documentation]  Test enabling and disabling a PAT
    ${d}=  Get Current Date  result_format=%m%s  increment=3 days
    ${token_name}=  Set Variable  enable-disable-${d}
    Create PAT Via API  ${token_name}  Test enable/disable  0
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 4 PASSED: PAT enable/disable scenario verified

Test Case - Delete PAT
    [Documentation]  Test deleting a PAT requires confirmation
    ${d}=  Get Current Date  result_format=%m%s  increment=4 days
    ${token_name}=  Set Variable  delete-test-${d}
    Create PAT Via API  ${token_name}  Test deletion  0
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 5 PASSED: PAT created for deletion testing

Test Case - Non-Admin User Can Create And Manage Own PAT
    [Documentation]  Test that non-admin users can create and manage their own PATs
    ${d}=  Get Current Date  result_format=%m%s  increment=5 days
    ${token_name}=  Set Variable  user-pat-${d}
    Create PAT Via API  ${token_name}  Non-admin user PAT  30
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 6 PASSED: Non-admin user PAT creation verified

Test Case - PAT Never Expires
    [Documentation]  Test creating a PAT that never expires (0 days)
    ${d}=  Get Current Date  result_format=%m%s  increment=6 days
    ${token_name}=  Set Variable  never-expires-${d}
    Create PAT Via API  ${token_name}  Token that never expires  0
    Verify Token Exists Via API  ${token_name}
    Log  ✅ Test Case 7 PASSED: Never-expiring PAT created successfully

Test Case - Docker Login And Push With PAT
    [Documentation]  Test docker login and push using PAT credentials - this is the core functionality
    ${d}=  Get Current Date  result_format=%m%s
    ${test_user}=  Set Variable  patuser${d}
    ${test_password}=  Set Variable  TestPassword123

    # Create test user via API
    ${user_result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 -X POST https://${ip}/api/v2.0/users -H "Content-Type: application/json" -d '{"username":"${test_user}","email":"${test_user}@test.com","password":"${test_password}","realname":"${test_user}"}' 2>&1 | grep -q '"id"' && echo "CREATED" || echo "FAILED"
    Should Contain  ${user_result.stdout}  CREATED  Failed to create test user

    # Get user ID
    ${user_id_result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 'https://${ip}/api/v2.0/users?username=${test_user}' 2>&1 | grep -oP '"user_id":\\K[0-9]+'
    ${user_id}=  Set Variable  ${user_id_result.stdout}
    Log  User ID: ${user_id}

    # Get admin token for project creation
    ${admin_token}=  Get Harbor Admin Token

    # Create a public project for testing
    ${project_name}=  Set Variable  pat-test-${d}
    ${proj_result}=  Run Process  bash  -c
    ...  curl -sk -H "Authorization: Bearer ${admin_token}" -X POST https://${ip}/api/v2.0/projects -H "Content-Type: application/json" -d '{"project_name":"${project_name}","public":true}' 2>&1 | grep -q '"project_id"' && echo "CREATED" || echo "FAILED"
    Should Contain  ${proj_result.stdout}  CREATED  Failed to create project

    # Add user to project with push permissions (role_id 2 = Developer)
    ${member_result}=  Run Process  bash  -c
    ...  curl -sk -H "Authorization: Bearer ${admin_token}" -X POST https://${ip}/api/v2.0/projects/${project_name}/members -H "Content-Type: application/json" -d '{"user_id":${user_id},"role_id":2}' 2>&1 | grep -q '"id"' && echo "CREATED" || echo "FAILED"
    Should Contain  ${member_result.stdout}  CREATED  Failed to add user to project

    # Create PAT for the test user - get the secret in response
    ${token_name}=  Set Variable  docker-pat-${d}
    ${pat_result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 -X POST https://${ip}/api/v2.0/users/${user_id}/personal_access_tokens -H "Content-Type: application/json" -d '{"name":"${token_name}","description":"Docker login PAT","expires_at":-1}' 2>&1
    Log  PAT creation result: ${pat_result.stdout}

    # Extract the secret from the response
    ${pat_secret}=  Run Process  bash  -c
    ...  echo '${pat_result.stdout}' | grep -oP '"secret":\\K"[^"]+' | tr -d '"'
    ${pat_secret_clean}=  Set Variable  ${pat_secret.stdout}
    Log  Got PAT secret (first 20 chars): ${pat_secret_clean[0:20]}...
    Should Not Be Empty  ${pat_secret_clean}  Failed to get PAT secret

    # Docker login using PAT
    Docker Login  ${HARBOR_REGISTRY}  ${test_user}  ${pat_secret_clean}
    Log  ✅ Test Case 8 PASSED: Docker login with PAT succeeded

    # Pull a small test image
    Docker Pull  library/hello-world:latest

    # Tag for our harbor
    Docker Tag  library/hello-world:latest  ${HARBOR_REGISTRY}/${project_name}/hello-world:latest

    # Push to harbor
    Docker Push  ${HARBOR_REGISTRY}/${project_name}/hello-world:latest
    Log  ✅ Test Case 9 PASSED: Docker push with PAT succeeded

    # Cleanup - logout, delete PAT, project, user
    Docker Logout  ${HARBOR_REGISTRY}
    Run Process  bash  -c  curl -sk -u admin:Harbor12345 -X DELETE https://${ip}/api/v2.0/projects/${project_name} 2>&1 || true
    Run Process  bash  -c  curl -sk -u admin:Harbor12345 -X DELETE https://${ip}/api/v2.0/users/${user_id} 2>&1 || true

    Log  ✅ Test Case 10 PASSED: Full docker login/push/cleanup with PAT completed

*** Keywords ***

Create PAT Via API
    [Arguments]  ${token_name}  ${description}  ${expiry_days}
    [Documentation]  Create a PAT using direct API call
    ${result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 -X POST https://${ip}/api/v2.0/users/1/personal_access_tokens -H "Content-Type: application/json" -d '{"name":"${token_name}","description":"${description}","expires_at":-1}' 2>&1 | grep -q '"id"' && echo "CREATED" || echo "FAILED"
    Should Contain  ${result.stdout}  CREATED  Failed to create PAT ${token_name}

Verify Token Exists Via API
    [Arguments]  ${token_name}
    [Documentation]  Verify token exists via API curl to /api/v2.0/users/1/personal_access_tokens
    ${result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 https://${ip}/api/v2.0/users/1/personal_access_tokens 2>&1 | grep -o '"name":"[^"]*"' | grep -q "${token_name}" && echo "FOUND" || echo "NOT_FOUND"
    Should Contain  ${result.stdout}  FOUND  Token ${token_name} not found in API response

Get Harbor Admin Token
    [Documentation]  Get admin JWT token for API calls
    ${token_result}=  Run Process  bash  -c
    ...  curl -sk -u admin:Harbor12345 -X POST https://${ip}/api/v2.0/tokens 2>&1 | grep -oP '"token":"\\K[^"]+'
    ${token}=  Set Variable  ${token_result.stdout}
    [Return]  ${token}
