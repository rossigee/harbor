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

package v2token

import (
	"context"
	"testing"

	"github.com/docker/distribution/registry/auth/token"
	"github.com/stretchr/testify/assert"

	"github.com/goharbor/harbor/src/common/rbac"
	rbac_project "github.com/goharbor/harbor/src/common/rbac/project"
	"github.com/goharbor/harbor/src/pkg/permission/types"
	"github.com/goharbor/harbor/src/pkg/project/models"
	"github.com/goharbor/harbor/src/testing/controller/project"
)

func TestAll(t *testing.T) {
	ctx := context.TODO()

	ctl := &project.Controller{}
	ctl.On("Get", ctx, int64(1)).Return(&models.Project{ProjectID: 1, Name: "library"}, nil)
	ctl.On("Get", ctx, int64(2)).Return(&models.Project{ProjectID: 2, Name: "test"}, nil)
	ctl.On("Get", ctx, int64(3)).Return(&models.Project{ProjectID: 3, Name: "rossgolderltd"}, nil)
	ctl.On("Get", ctx, int64(4)).Return(&models.Project{ProjectID: 4, Name: "development"}, nil)

	access := []*token.ResourceActions{
		{
			Type: "repository",
			Name: "library/ubuntu",
			Actions: []string{
				"pull",
				"push",
				"scanner-pull",
			},
		},
		{
			Type: "repository",
			Name: "test/golang",
			Actions: []string{
				"pull",
				"*",
			},
		},
		{
			Type: "cnab",
			Name: "development/cnab",
			Actions: []string{
				"pull",
				"push",
			},
		},
	}
	sc := New(context.Background(), "jack", access)
	tsc := sc.(*tokenSecurityCtx)
	tsc.ctl = ctl

	cases := []struct {
		resource types.Resource
		action   types.Action
		expect   bool
	}{
		{
			resource: rbac_project.NewNamespace(1).Resource(rbac.ResourceRepository),
			action:   rbac.ActionPush,
			expect:   true,
		},
		{
			resource: rbac_project.NewNamespace(1).Resource(rbac.ResourceRepository),
			action:   rbac.ActionScannerPull,
			expect:   true,
		},
		{
			resource: rbac_project.NewNamespace(2).Resource(rbac.ResourceRepository),
			action:   rbac.ActionPush,
			expect:   true,
		},
		{
			resource: rbac_project.NewNamespace(2).Resource(rbac.ResourceRepository),
			action:   rbac.ActionDelete,
			expect:   true,
		},
		{
			resource: rbac_project.NewNamespace(2).Resource(rbac.ResourceRepository),
			action:   rbac.ActionScannerPull,
			expect:   false,
		},
		{
			resource: rbac_project.NewNamespace(4).Resource(rbac.ResourceRepository),
			action:   rbac.ActionPush,
			expect:   false,
		},
		{
			resource: rbac_project.NewNamespace(2).Resource(rbac.ResourceArtifact),
			action:   rbac.ActionPush,
			expect:   false,
		},
		{
			resource: rbac_project.NewNamespace(1).Resource(rbac.ResourceRepository),
			action:   rbac.ActionCreate,
			expect:   false,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expect, sc.Can(ctx, c.action, c.resource))
	}
}

func TestRobotTokenAccess(t *testing.T) {
	// This test specifically validates the robot token access scenario
	// that was fixed - robot with push/pull on a specific project
	ctx := context.TODO()

	ctl := &project.Controller{}
	// Simulate the rossgolderltd project that robot-cicd-rossg has access to
	ctl.On("Get", ctx, int64(3)).Return(&models.Project{ProjectID: 3, Name: "rossgolderltd"}, nil)

	// Robot token access for rossgolderltd/test with push and pull
	access := []*token.ResourceActions{
		{
			Type:    "repository",
			Name:    "rossgolderltd/test",
			Actions: []string{"push", "pull"},
		},
	}

	sc := New(context.Background(), "robot-cicd-rossg", access)
	tsc := sc.(*tokenSecurityCtx)
	tsc.ctl = ctl

	// Test push permission on the repository
	resource := rbac_project.NewNamespace(3).Resource(rbac.ResourceRepository)
	assert.True(t, sc.Can(ctx, rbac.ActionPush, resource), "Robot should have push permission")
	assert.True(t, sc.Can(ctx, rbac.ActionPull, resource), "Robot should have pull permission")
	assert.False(t, sc.Can(ctx, rbac.ActionDelete, resource), "Robot should not have delete permission")

	// Verify the security context properties
	assert.Equal(t, "robot-cicd-rossg", sc.GetUsername())
	assert.True(t, sc.IsAuthenticated())
	assert.False(t, sc.IsSysAdmin())
}

func TestTokenWithWildcardAction(t *testing.T) {
	// Test that wildcard '*' action grants all permissions
	ctx := context.TODO()

	ctl := &project.Controller{}
	ctl.On("Get", ctx, int64(1)).Return(&models.Project{ProjectID: 1, Name: "library"}, nil)

	access := []*token.ResourceActions{
		{
			Type:    "repository",
			Name:    "library/ubuntu",
			Actions: []string{"*"},
		},
	}

	sc := New(context.Background(), "admin", access)
	tsc := sc.(*tokenSecurityCtx)
	tsc.ctl = ctl

	resource := rbac_project.NewNamespace(1).Resource(rbac.ResourceRepository)
	assert.True(t, sc.Can(ctx, rbac.ActionPush, resource), "Wildcard should include push")
	assert.True(t, sc.Can(ctx, rbac.ActionPull, resource), "Wildcard should include pull")
	assert.True(t, sc.Can(ctx, rbac.ActionDelete, resource), "Wildcard should include delete")
}

func TestTokenWithUnsupportedType(t *testing.T) {
	// Test that non-repository types are ignored
	ctx := context.TODO()

	ctl := &project.Controller{}
	ctl.On("Get", ctx, int64(1)).Return(&models.Project{ProjectID: 1, Name: "library"}, nil)

	access := []*token.ResourceActions{
		{
			Type:    "repository",
			Name:    "library/ubuntu",
			Actions: []string{"push", "pull"},
		},
		{
			Type:    "helm-chart",
			Name:    "library/mychart",
			Actions: []string{"pull"},
		},
	}

	sc := New(context.Background(), "user", access)
	tsc := sc.(*tokenSecurityCtx)
	tsc.ctl = ctl

	// Repository access should work
	resource := rbac_project.NewNamespace(1).Resource(rbac.ResourceRepository)
	assert.True(t, sc.Can(ctx, rbac.ActionPush, resource))
}

func TestTokenNameExtraction(t *testing.T) {
	// Test that project name is correctly extracted from repository name
	ctx := context.TODO()

	ctl := &project.Controller{}
	ctl.On("Get", ctx, int64(5)).Return(&models.Project{ProjectID: 5, Name: "myproject"}, nil)

	access := []*token.ResourceActions{
		{
			Type:    "repository",
			Name:    "myproject/some/image",
			Actions: []string{"push", "pull"},
		},
	}

	sc := New(context.Background(), "robot", access)
	tsc := sc.(*tokenSecurityCtx)
	tsc.ctl = ctl

	resource := rbac_project.NewNamespace(5).Resource(rbac.ResourceRepository)
	assert.True(t, sc.Can(ctx, rbac.ActionPush, resource), "Should have access to myproject/some/image")
}
