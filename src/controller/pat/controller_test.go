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

package pat

import (
	"fmt"
	"testing"
	"time"

	htesting "github.com/goharbor/harbor/src/testing"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/controller/user"
	"github.com/goharbor/harbor/src/lib/q"
	"github.com/goharbor/harbor/src/pkg/pat/model"
)

type ControllerTestSuite struct {
	htesting.Suite
	ctl     Controller
	userCtl user.Controller
}

func (suite *ControllerTestSuite) SetupSuite() {
	suite.Suite.SetupSuite()
	suite.ClearTables = []string{"personal_access_token", "harbor_user"}
	suite.ctl = NewController()
	suite.userCtl = user.Ctl
}

func (suite *ControllerTestSuite) createTestUser(username string) int {
	ctx := suite.Context()
	u := &models.User{
		Username: username,
		Email:    fmt.Sprintf("%s@example.com", username),
		Realname: username,
	}
	uid, err := suite.userCtl.Create(ctx, u)
	suite.NoError(err)
	return int(uid)
}

func (suite *ControllerTestSuite) TestCreateGeneratesSecretWithPrefix() {
	userID := suite.createTestUser("testuser1")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "test-token",
		Description: "Test token",
		ExpiresAt:   time.Now().AddDate(0, 0, 30).Unix(),
	}

	id, plaintext, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)
	suite.True(id > 0)

	// Verify plaintext secret has hbr_pat_ prefix
	suite.Contains(plaintext, "hbr_pat_")
	suite.True(len(plaintext) > 8) // hbr_pat_ = 8 chars + random
}

func (suite *ControllerTestSuite) TestCreateValidatesSecretFormat() {
	userID := suite.createTestUser("testuser2")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "test-token",
		Description: "Test",
		ExpiresAt:   -1, // Never expires
	}

	id, plaintext, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)
	suite.NotEmpty(plaintext)
	suite.True(id > 0)

	// Verify secret is secure (no obviously weak format)
	suite.NotContains(plaintext, "admin")
	suite.NotContains(plaintext, "user")
	suite.NotContains(plaintext, "password")
}

func (suite *ControllerTestSuite) TestCreateWithExpiry() {
	userID := suite.createTestUser("testuser3")
	expiryTime := time.Now().AddDate(0, 0, 90).Unix()
	pat := &model.PersonalAccessToken{
		UserID:    userID,
		Name:      "expiring-token",
		ExpiresAt: expiryTime,
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	// Retrieve and verify expiry
	retrieved, err := suite.ctl.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal(expiryTime, retrieved.ExpiresAt)
}

func (suite *ControllerTestSuite) TestCreateNeverExpiresToken() {
	userID := suite.createTestUser("testuser4")
	pat := &model.PersonalAccessToken{
		UserID:    userID,
		Name:      "never-expire",
		ExpiresAt: -1, // -1 means never
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	retrieved, err := suite.ctl.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal(int64(-1), retrieved.ExpiresAt)
}

func (suite *ControllerTestSuite) TestGetReturnsToken() {
	userID := suite.createTestUser("testuser5")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "get-test",
		Description: "Get test",
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	retrieved, err := suite.ctl.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal(pat.UserID, retrieved.UserID)
	suite.Equal(pat.Name, retrieved.Name)
	suite.Equal(pat.Description, retrieved.Description)
}

func (suite *ControllerTestSuite) TestListTokensForUser() {
	userID := suite.createTestUser("testuser6")
	for i := 1; i <= 3; i++ {
		pat := &model.PersonalAccessToken{
			UserID: userID,
			Name:   "token-" + string(rune(48+i)),
		}
		_, _, err := suite.ctl.Create(suite.Context(), pat)
		suite.NoError(err)
	}

	query := q.New(q.KeyWords{"user_id": userID})
	pats, err := suite.ctl.List(suite.Context(), query)
	suite.NoError(err)
	suite.Equal(3, len(pats))

	for _, p := range pats {
		suite.Equal(userID, p.UserID)
	}
}

func (suite *ControllerTestSuite) TestUpdateTokenMetadata() {
	userID := suite.createTestUser("testuser7")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "update-test",
		Description: "Original",
		Disabled:    false,
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	// Update
	pat.ID = id
	pat.Description = "Updated description"
	pat.Disabled = true

	err = suite.ctl.Update(suite.Context(), pat)
	suite.NoError(err)

	// Verify
	updated, err := suite.ctl.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal("Updated description", updated.Description)
	suite.True(updated.Disabled)
}

func (suite *ControllerTestSuite) TestDeleteToken() {
	userID := suite.createTestUser("testuser8")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "delete-test",
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	err = suite.ctl.Delete(suite.Context(), id)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.ctl.Get(suite.Context(), id)
	suite.Error(err)
}

func (suite *ControllerTestSuite) TestRefreshSecretGeneratesNewSecret() {
	userID := suite.createTestUser("testuser9")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "refresh-test",
	}

	id, originalSecret, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	// Refresh without providing new secret (auto-generate)
	newSecret, err := suite.ctl.RefreshSecret(suite.Context(), id, "")
	suite.NoError(err)

	// Verify new secret is different
	suite.NotEqual(originalSecret, newSecret)
	suite.Contains(newSecret, "hbr_pat_")
}

func (suite *ControllerTestSuite) TestRefreshSecretWithProvidedSecret() {
	userID := suite.createTestUser("testuser10")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "refresh-provided",
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	// Refresh with a specific secret that passes validation
	// (must have uppercase, lowercase, digit, and be 8-128 chars)
	newSecret := "TestSecret123"
	refreshedSecret, err := suite.ctl.RefreshSecret(suite.Context(), id, newSecret)
	suite.NoError(err)

	// Verify the refreshed secret contains the provided value
	suite.Contains(refreshedSecret, "TestSecret123")
}

func (suite *ControllerTestSuite) TestRefreshSecretInvalidFormat() {
	userID := suite.createTestUser("testuser11")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "invalid-refresh",
	}

	id, _, err := suite.ctl.Create(suite.Context(), pat)
	suite.NoError(err)

	// Try to refresh with invalid secret (no uppercase, lowercase, or digit)
	_, err = suite.ctl.RefreshSecret(suite.Context(), id, "invalid")
	suite.Error(err)
}

func (suite *ControllerTestSuite) TestCountTokens() {
	userID := suite.createTestUser("testuser12")
	expectedCount := 2

	for i := 0; i < expectedCount; i++ {
		pat := &model.PersonalAccessToken{
			UserID: userID,
			Name:   "token-" + string(rune(48+i)),
		}
		_, _, err := suite.ctl.Create(suite.Context(), pat)
		suite.NoError(err)
	}

	query := q.New(q.KeyWords{"user_id": userID})
	count, err := suite.ctl.Count(suite.Context(), query)
	suite.NoError(err)
	suite.Equal(int64(expectedCount), count)
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
