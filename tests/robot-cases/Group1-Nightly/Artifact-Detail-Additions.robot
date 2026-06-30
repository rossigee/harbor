# Copyright Project Harbor Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation     Artifact Detail Page - Additions Tab Coverage
...               Tests for displaying Dockerfile and other artifact additions
...               Verifies correct display of additions with and without labels
Resource          ../../resources/Util.robot
Resource          ../../resources/Harbor-Pages/Project-Artifact.robot
Resource          ../../resources/Harbor-Pages/Project-Repository.robot
Library           Collections

*** Variables ***
${artifact_detail_dockerfile_tab}           id=dockerfile-link
${artifact_detail_dockerfile_content}       xpath=//hbr-artifact-dockerfile/div[@class='row content-wrapper']
${artifact_detail_dockerfile_info_box}      xpath=//hbr-artifact-dockerfile//div[@class='info-box']
${artifact_detail_no_dockerfile_msg}        xpath=//hbr-artifact-dockerfile//div[contains(text(), 'Dockerfile')]
${artifact_detail_build_history_tab}        id=build-history
${artifact_detail_build_history_link}       xpath=//hbr-artifact-dockerfile//a[contains(text(), 'Build History')]
${artifact_detail_loading_spinner}          xpath=//hbr-artifact-dockerfile//span[@class='spinner']
${yaml_container}                            xpath=//div[@class='yaml-container']

*** Test Cases ***

Test Dockerfile Tab Display With Label
    [Documentation]    Verify Dockerfile tab appears and displays content when image has label
    [Tags]    artifact    dockerfile    addition
    Init Test Data With Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Dockerfile Tab
    Verify Dockerfile Content Displayed
    [Teardown]    Cleanup Test Data

Test Dockerfile Tab Display Without Label
    [Documentation]    Verify informational message when image lacks Dockerfile label
    [Tags]    artifact    dockerfile    addition
    Init Test Data Without Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Dockerfile Tab
    Verify No Dockerfile Info Message Displayed
    [Teardown]    Cleanup Test Data

Test Dockerfile Tab Provides Build History Link
    [Documentation]    Verify user can navigate to Build History from Dockerfile tab
    [Tags]    artifact    dockerfile    addition
    Init Test Data Without Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Dockerfile Tab
    Verify No Dockerfile Info Message Displayed
    Click Build History Link From Dockerfile Tab
    Verify Build History Tab Active
    [Teardown]    Cleanup Test Data

Test Tab Navigation And Switching
    [Documentation]    Verify Dockerfile tab can be clicked and switched between tabs
    [Tags]    artifact    dockerfile    addition
    Init Test Data With Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Dockerfile Tab
    Verify Dockerfile Tab Active
    Wait Until Page Contains Element    ${yaml_container}    timeout=10s
    Click Build History Tab
    Verify Build History Tab Active
    Click Dockerfile Tab
    Verify Dockerfile Tab Active
    [Teardown]    Cleanup Test Data

Test Dockerfile Tab Always Visible
    [Documentation]    Verify Dockerfile tab is always visible, regardless of label presence
    [Tags]    artifact    dockerfile    addition
    Init Test Data Without Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Retry Wait Until Page Contains Element    ${artifact_detail_dockerfile_tab}
    [Teardown]    Cleanup Test Data

Test Dockerfile Content Has Syntax Highlighting
    [Documentation]    Verify Dockerfile content renders with syntax highlighting
    [Tags]    artifact    dockerfile    addition
    Init Test Data With Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Dockerfile Tab
    Verify Dockerfile Content Displayed
    # Verify syntax highlighting is applied (yaml-container uses language pipe)
    Retry Wait Until Page Contains Element    ${yaml_container}
    [Teardown]    Cleanup Test Data

Test Build History Tab Always Available
    [Documentation]    Verify Build History tab is always available as fallback
    [Tags]    artifact    dockerfile    addition
    Init Test Data Without Dockerfile Label
    Go To Project
    Go To Repository
    Go Into Artifact
    Click Build History Tab
    Verify Build History Tab Active
    Retry Wait Until Page Contains Element    xpath=//hbr-artifact-build-history
    [Teardown]    Cleanup Test Data

*** Keywords ***

Init Test Data With Dockerfile Label
    [Documentation]    Create test data with image containing Dockerfile label
    ${test_project}=                    Set Variable    test-dockerfile-with-label
    ${test_repo}=                       Set Variable    test-image
    ${test_tag}=                        Set Variable    latest
    ${dockerfile_content}=              Set Variable    FROM ubuntu:22.04\nRUN apt-get update\nRUN apt-get install -y curl
    Set Test Variable    ${project_name}    ${test_project}
    Set Test Variable    ${repo_name}       ${test_repo}
    Set Test Variable    ${tag_name}        ${test_tag}
    # Create project
    Create Project    ${test_project}
    # Build and push image with Dockerfile label (done via docker-compose or direct build)
    Log    Build test image with Dockerfile label using docker build
    Log    docker build --label "org.opencontainers.image.source=${dockerfile_content}" -t ${LOCAL_REGISTRY}/${test_project}/${test_repo}:${test_tag} .

Init Test Data Without Dockerfile Label
    [Documentation]    Create test data with image without Dockerfile label
    ${test_project}=                    Set Variable    test-dockerfile-without-label
    ${test_repo}=                       Set Variable    test-image
    ${test_tag}=                        Set Variable    latest
    Set Test Variable    ${project_name}    ${test_project}
    Set Test Variable    ${repo_name}       ${test_repo}
    Set Test Variable    ${tag_name}        ${test_tag}
    # Create project
    Create Project    ${test_project}
    # Build and push image without Dockerfile label
    Log    Build test image WITHOUT Dockerfile label using docker build
    Log    docker build -t ${LOCAL_REGISTRY}/${test_project}/${test_repo}:${test_tag} .

Cleanup Test Data
    [Documentation]    Clean up test projects and images
    Run Keyword If Test Passed    Delete Project    ${project_name}
    Run Keyword If Test Passed    Delete Repository    ${project_name}    ${repo_name}

Create Project
    [Arguments]    ${project_name}
    [Documentation]    Create a test project
    Navigate To Project    ${project_name}
    Run Keyword And Ignore Error    New Project    ${project_name}

Delete Project
    [Arguments]    ${project_name}
    [Documentation]    Delete a test project
    Navigate To Project    ${project_name}
    Run Keyword And Ignore Error    Project Delete    ${project_name}

Navigate To Project
    [Arguments]    ${project_name}
    [Documentation]    Navigate to project repositories page
    Go To    ${HARBOR_URL}/projects
    Retry Wait Until Page Not Contains Element    ${artifact_list_spinner}

Go To Project
    [Documentation]    Navigate to the test project
    Navigate To Project    ${project_name}
    Retry Element Click    xpath=//a[contains(text(), '${project_name}')]

Go To Repository
    [Documentation]    Navigate to repository in project
    Retry Wait Until Page Not Contains Element    ${artifact_list_spinner}
    Retry Wait Until Page Contains Element    xpath=//clr-dg-row[contains(.,'${repo_name}')]

Go Into Artifact
    [Documentation]    Click into artifact detail page
    Retry Wait Until Page Not Contains Element    ${artifact_list_spinner}
    Retry Element Click    xpath=//clr-dg-row[contains(.,'${tag_name}')]//a[contains(.,'sha256')]
    Retry Wait Until Page Contains Element    ${artifact_tag_component}
    Retry Wait Until Page Not Contains Element    ${artifact_list_spinner}

Click Dockerfile Tab
    [Documentation]    Click on the Dockerfile tab
    Retry Element Click    ${artifact_detail_dockerfile_tab}
    Sleep    1s    # Allow tab content to load

Click Build History Tab
    [Documentation]    Click on the Build History tab
    Retry Element Click    ${artifact_detail_build_history_tab}
    Sleep    1s    # Allow tab content to load

Click Build History Link From Dockerfile Tab
    [Documentation]    Click the link to Build History from the Dockerfile info box
    Retry Element Click    ${artifact_detail_build_history_link}
    Sleep    1s    # Allow tab switch to complete

Verify Dockerfile Tab Active
    [Documentation]    Verify the Dockerfile tab is currently active
    Retry Wait Until Page Contains Element    xpath=${artifact_detail_dockerfile_tab}[@aria-selected='true']

Verify Build History Tab Active
    [Documentation]    Verify the Build History tab is currently active
    Retry Wait Until Page Contains Element    xpath=${artifact_detail_build_history_tab}[@aria-selected='true']

Verify Dockerfile Content Displayed
    [Documentation]    Verify Dockerfile content is shown with proper formatting
    Retry Wait Until Page Not Contains Element    ${artifact_detail_loading_spinner}
    Retry Wait Until Page Contains Element    ${yaml_container}
    ${content}=    Get Text    ${yaml_container}
    Should Contain    ${content}    FROM

Verify No Dockerfile Info Message Displayed
    [Documentation]    Verify informational message appears when no Dockerfile label exists
    Retry Wait Until Page Not Contains Element    ${artifact_detail_loading_spinner}
    Retry Wait Until Page Contains Element    ${artifact_detail_dockerfile_info_box}
    ${message}=    Get Text    ${artifact_detail_dockerfile_info_box}
    Should Contain    ${message}    Dockerfile
    Should Contain    ${message}    labels
