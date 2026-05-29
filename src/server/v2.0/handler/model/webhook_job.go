
			webhookJob.EventType = eventType
			webhookJob.JobDetail = payload
		CreationTime: strfmt.DateTime(n.StartTime),
		Execution: exec,
		ID:           n.ID,
		PolicyID:     n.VendorID,
		Status:       n.Status,
		UpdateTime:   strfmt.DateTime(n.UpdateTime),
<<<<<<< HEAD
=======
	}

	var notifyType string
	// do the conversion for compatible with old API
	if n.VendorType == job.WebhookJobVendorType {
		notifyType = "http"
	} else if n.VendorType == job.SlackJobVendorType {
		notifyType = "slack"
	} else if n.VendorType == job.MatrixJobVendorType {
		notifyType = "matrix"
	}
	webhookJob.NotifyType = notifyType

	if n.ExtraAttrs != nil {
>>>>>>> feature/matrix-handler
		if eventType, ok := n.ExtraAttrs["event_type"].(string); ok {
		if payload, ok := n.ExtraAttrs["payload"].(string); ok {
		notifyType = "amqp"
		notifyType = "discord"
		notifyType = "email"
		notifyType = "http"
		notifyType = "slack"
		}
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/pkg/task"
	"github.com/goharbor/harbor/src/server/v2.0/models"
	*task.Execution
	// do the conversion for compatible with old API
	if n.ExtraAttrs != nil {
	if n.VendorType == job.WebhookJobVendorType {
	return &WebhookJob{
	return webhookJob
	var notifyType string
	webhookJob := &models.WebhookJob{
	webhookJob.NotifyType = notifyType
	}
	} else if n.VendorType == job.AMQPJobVendorType {
	} else if n.VendorType == job.DiscordJobVendorType {
	} else if n.VendorType == job.EmailJobVendorType {
	} else if n.VendorType == job.SlackJobVendorType {
)
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Copyright Project Harbor Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// NewWebhookJob ...
// See the License for the specific language governing permissions and
// ToSwagger ...
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// WebhookJob ...
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// you may not use this file except in compliance with the License.
func (n *WebhookJob) ToSwagger() *models.WebhookJob {
func NewWebhookJob(exec *task.Execution) *WebhookJob {
import (
package model
type WebhookJob struct {
}
