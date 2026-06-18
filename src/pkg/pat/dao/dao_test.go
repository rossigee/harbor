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

package dao

import (
	"fmt"
	"testing"
	"time"

	htesting "github.com/goharbor/harbor/src/testing"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/controller/user"
	"github.com/goharbor/harbor/src/lib/errors"
	"github.com/goharbor/harbor/src/lib/q"
	"github.com/goharbor/harbor/src/pkg/pat/model"
)

type DAOTestSuite struct {
	htesting.Suite
	dao    DAO
	userCtl user.Controller
}

func (suite *DAOTestSuite) SetupSuite() {
	suite.Suite.SetupSuite()
	suite.ClearTables = []string{"personal_access_token", "harbor_user"}
	suite.dao = New()
	suite.userCtl = user.Ctl
}

func (suite *DAOTestSuite) createTestUser(username string) int {
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

func (suite *DAOTestSuite) TestCreate() {
	userID := suite.createTestUser("testuser1")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "test-token",
		Secret:      "hashed_secret",
		Salt:        "salt_value",
		Description: "Test token",
		ExpiresAt:   time.Now().AddDate(0, 0, 30).Unix(),
	}

	id, err := suite.dao.Create(suite.Context(), pat)
	suite.NoError(err)
	suite.True(id > 0)
}

func (suite *DAOTestSuite) TestCreateDuplicate() {
	userID := suite.createTestUser("testuser2")
	pat1 := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "duplicate-token",
		Secret: "secret1",
		Salt:   "salt1",
	}

	pat2 := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "duplicate-token",
		Secret: "secret2",
		Salt:   "salt2",
	}

	// First create should succeed
	id, err := suite.dao.Create(suite.Context(), pat1)
	suite.NoError(err)
	suite.True(id > 0)

	// Second create with same user_id and name should fail
	_, err = suite.dao.Create(suite.Context(), pat2)
	suite.Error(err)
	suite.True(errors.IsConflictErr(err))
}

func (suite *DAOTestSuite) TestGet() {
	userID := suite.createTestUser("testuser3")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "get-test",
		Secret: "secret",
		Salt:   "salt",
	}

	// Create first
	id, err := suite.dao.Create(suite.Context(), pat)
	suite.NoError(err)

	// Get
	retrieved, err := suite.dao.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal(pat.UserID, retrieved.UserID)
	suite.Equal(pat.Name, retrieved.Name)
}

func (suite *DAOTestSuite) TestGetNotFound() {
	_, err := suite.dao.Get(suite.Context(), 99999)
	suite.Error(err)
	suite.True(errors.IsNotFoundErr(err))
}

func (suite *DAOTestSuite) TestUpdate() {
	userID := suite.createTestUser("testuser4")
	pat := &model.PersonalAccessToken{
		UserID:      userID,
		Name:        "update-test",
		Secret:      "secret",
		Salt:        "salt",
		Description: "Original",
		Disabled:    false,
	}

	// Create
	id, err := suite.dao.Create(suite.Context(), pat)
	suite.NoError(err)

	// Update
	pat.ID = id
	pat.Description = "Updated"
	pat.Disabled = true

	err = suite.dao.Update(suite.Context(), pat, "description", "disabled")
	suite.NoError(err)

	// Verify
	updated, err := suite.dao.Get(suite.Context(), id)
	suite.NoError(err)
	suite.Equal("Updated", updated.Description)
	suite.True(updated.Disabled)
}

func (suite *DAOTestSuite) TestDelete() {
	userID := suite.createTestUser("testuser5")
	pat := &model.PersonalAccessToken{
		UserID: userID,
		Name:   "delete-test",
		Secret: "secret",
		Salt:   "salt",
	}

	// Create
	id, err := suite.dao.Create(suite.Context(), pat)
	suite.NoError(err)

	// Delete
	err = suite.dao.Delete(suite.Context(), id)
	suite.NoError(err)

	// Verify it's gone
	_, err = suite.dao.Get(suite.Context(), id)
	suite.Error(err)
	suite.True(errors.IsNotFoundErr(err))
}

func (suite *DAOTestSuite) TestList() {
	userID := suite.createTestUser("testuser6")
	for i := 1; i <= 3; i++ {
		pat := &model.PersonalAccessToken{
			UserID: userID,
			Name:   "token-" + string(rune(i)),
			Secret: "secret",
			Salt:   "salt",
		}
		_, err := suite.dao.Create(suite.Context(), pat)
		suite.NoError(err)
	}

	// List all for user
	query := q.New(q.KeyWords{"user_id": userID})
	pats, err := suite.dao.List(suite.Context(), query)
	suite.NoError(err)
	suite.Equal(3, len(pats))
}

func (suite *DAOTestSuite) TestCount() {
	userID := suite.createTestUser("testuser7")
	for i := 1; i <= 2; i++ {
		pat := &model.PersonalAccessToken{
			UserID: userID,
			Name:   "token-" + string(rune(i)),
			Secret: "secret",
			Salt:   "salt",
		}
		_, err := suite.dao.Create(suite.Context(), pat)
		suite.NoError(err)
	}

	// Count
	query := q.New(q.KeyWords{"user_id": userID})
	count, err := suite.dao.Count(suite.Context(), query)
	suite.NoError(err)
	suite.Equal(int64(2), count)
}

func TestDAOTestSuite(t *testing.T) {
	suite.Run(t, new(DAOTestSuite))
}
