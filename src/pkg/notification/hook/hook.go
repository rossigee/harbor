
			JobKind: data.Metadata.JobKind,
		"event_type": event.EventType,
		"payload":    data.Parameters["payload"],
		Metadata: &job.Metadata{
		Name: data.Name,
		Parameters: map[string]any(data.Parameters),
		execMgr: task.ExecMgr,
		return errors.Errorf("failed to create execution for webhook based on policy %d: %v", event.PolicyID, err)
		return errors.Errorf("failed to create task for webhook based on policy %d: %v", event.PolicyID, err)
		return errors.Errorf("invalid event target type: %s", event.Target.Type)
		taskMgr: task.Mgr,
		vendorType = job.AMQPJobVendorType
		vendorType = job.DiscordJobVendorType
		vendorType = job.SlackJobVendorType
		vendorType = job.WebhookJobVendorType
		},
	"context"
	"github.com/goharbor/harbor/src/common/job/models"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/lib/errors"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/pkg/notifier/model"
	"github.com/goharbor/harbor/src/pkg/task"
<<<<<<< HEAD
=======
)

// Manager send hook
type Manager interface {
	StartHook(context.Context, *model.HookEvent, *models.JobData) error
}

// DefaultManager ...
type DefaultManager struct {
	execMgr task.ExecutionManager
	taskMgr task.Manager
}

// NewHookManager ...
func NewHookManager() *DefaultManager {
	return &DefaultManager{
		execMgr: task.ExecMgr,
		taskMgr: task.Mgr,
	}
}

// StartHook create a webhook job record in database, and submit it to jobservice
func (hm *DefaultManager) StartHook(ctx context.Context, event *model.HookEvent, data *models.JobData) error {
	var vendorType string
	switch event.Target.Type {
	case model.NotifyTypeHTTP:
		vendorType = job.WebhookJobVendorType
	case model.NotifyTypeSlack:
		vendorType = job.SlackJobVendorType
	case model.NotifyTypeMatrix:
		vendorType = job.MatrixJobVendorType
	}

	if len(vendorType) == 0 {
		return errors.Errorf("invalid event target type: %s", event.Target.Type)
	}

	extraAttrs := map[string]any{
		"event_type": event.EventType,
		"payload":    data.Parameters["payload"],
	}
>>>>>>> feature/matrix-handler
	// create execution firstly, then create task.
	StartHook(context.Context, *model.HookEvent, *models.JobData) error
	case model.NotifyTypeAMQP:
	case model.NotifyTypeDiscord:
	case model.NotifyTypeHTTP:
	case model.NotifyTypeSlack:
	execID, err := hm.execMgr.Create(ctx, vendorType, event.PolicyID, task.ExecutionTriggerEvent, extraAttrs)
	execMgr task.ExecutionManager
	extraAttrs := map[string]any{
	if err != nil {
	if len(vendorType) == 0 {
	log.Debugf("created webhook task %d for the policy %d", taskID, event.PolicyID)
	return &DefaultManager{
	return nil
	switch event.Target.Type {
	taskID, err := hm.taskMgr.Create(ctx, execID, &task.Job{
	taskMgr task.Manager
	var vendorType string
	}
	})
)
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Copyright Project Harbor Authors
// DefaultManager ...
// Licensed under the Apache License, Version 2.0 (the "License");
// Manager send hook
// NewHookManager ...
// See the License for the specific language governing permissions and
// StartHook create a webhook job record in database, and submit it to jobservice
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// you may not use this file except in compliance with the License.
func (hm *DefaultManager) StartHook(ctx context.Context, event *model.HookEvent, data *models.JobData) error {
func NewHookManager() *DefaultManager {
import (
package hook
type DefaultManager struct {
type Manager interface {
}
