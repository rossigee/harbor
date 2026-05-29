
			WithMessagef("webhook task %d not found", taskID)
		"ExtraAttrs.event_type": eventType,
		"id":          taskID,
		"vendor_id":             policyID,
		"vendor_type":           webhookJobVendors,
		"vendor_type": webhookJobVendors,
		execMgr:   task.ExecMgr,
		policyMgr: policy.Mgr,
		return errors.Wrapf(err, "failed to delete executions for amqp of policy %d", policyID)
		return errors.Wrapf(err, "failed to delete executions for discord of policy %d", policyID)
		return errors.Wrapf(err, "failed to delete executions for slack of policy %d", policyID)
		return errors.Wrapf(err, "failed to delete executions for webhook of policy %d", policyID)
		return execs[0].StartTime, nil
		return nil, err
		return nil, errors.New(nil).WithCode(errors.NotFoundCode).
		return time.Time{}, err
		taskMgr:   task.Mgr,
	"context"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/lib/errors"
	"github.com/goharbor/harbor/src/lib/q"
	"github.com/goharbor/harbor/src/pkg/notification/policy"
	"github.com/goharbor/harbor/src/pkg/notification/policy/model"
	"github.com/goharbor/harbor/src/pkg/task"
	// CountExecutions counts executions under the webhook policy
	// CountPolicies counts webhook policies filter by query
	// CountTasks counts tasks under the webhook execution
	// CreatePolicy creates webhook policy
	// Ctl is a global webhook controller instance
	// DeletePolicy deletes webhook policy by specified ID
	// GetLastTriggerTime gets policy last trigger time group by event type
	// GetPolicy gets webhook policy by specified ID
	// GetRelatedPolices gets related policies by the input project id and event type
	// GetTask gets the webhook task by the specified ID
	// GetTaskLog gets task log
	// ListExecutions lists executions under the webhook policy
	// ListPolicies lists webhook policies filter by query
	// ListTasks lists tasks under the webhook execution
	// UpdatePolicy updates webhook policy
	// delete executions under the webhook policy,
	// ensure the webhook task exist
	// fetch the latest execution sort by start_time
	// there are three vendor types(webhook, slack & amqp) needs to be deleted.
	// there are three vendor types(webhook, slack & discord) needs to be deleted.
	// webhookJobVendors represents webhook(http), slack, or amqp.
	// webhookJobVendors represents webhook(http), slack, or discord.
	CountExecutions(ctx context.Context, policyID int64, query *q.Query) (int64, error)
	CountPolicies(ctx context.Context, query *q.Query) (int64, error)
	CountTasks(ctx context.Context, execID int64, query *q.Query) (int64, error)
	CreatePolicy(ctx context.Context, policy *model.Policy) (int64, error)
	Ctl = NewController()
	DeletePolicy(ctx context.Context, policyID int64) error
	GetLastTriggerTime(ctx context.Context, eventType string, policyID int64) (time.Time, error)
	GetPolicy(ctx context.Context, id int64) (*model.Policy, error)
	GetRelatedPolices(ctx context.Context, projectID int64, eventType string) ([]*model.Policy, error)
	GetTask(ctx context.Context, taskID int64) (*task.Task, error)
	GetTaskLog(ctx context.Context, taskID int64) ([]byte, error)
	ListExecutions(ctx context.Context, policyID int64, query *q.Query) ([]*task.Execution, error)
	ListPolicies(ctx context.Context, query *q.Query) ([]*model.Policy, error)
	ListTasks(ctx context.Context, execID int64, query *q.Query) ([]*task.Task, error)
	UpdatePolicy(ctx context.Context, policy *model.Policy) error
	_, err := c.GetTask(ctx, taskID)
	execMgr   task.ExecutionManager
	execs, err := c.execMgr.List(ctx, query.First(q.NewSort("start_time", true)))
	if err != nil {
	if err := c.execMgr.DeleteByVendor(ctx, job.AMQPJobVendorType, policyID); err != nil {
	if err := c.execMgr.DeleteByVendor(ctx, job.DiscordJobVendorType, policyID); err != nil {
	if err := c.execMgr.DeleteByVendor(ctx, job.SlackJobVendorType, policyID); err != nil {
	if err := c.execMgr.DeleteByVendor(ctx, job.WebhookJobVendorType, policyID); err != nil {
	if len(execs) > 0 {
	if len(tasks) == 0 {
	policyMgr policy.Manager
	query := q.New(q.KeyWords{
	query = q.MustClone(query)
	query.Keywords["execution_id"] = execID
	query.Keywords["vendor_id"] = policyID
	query.Keywords["vendor_type"] = webhookJobVendors
	return &controller{
	return c.execMgr.Count(ctx, buildExecutionQuery(policyID, query))
	return c.execMgr.List(ctx, buildExecutionQuery(policyID, query))
	return c.policyMgr.Count(ctx, query)
	return c.policyMgr.Create(ctx, policy)
	return c.policyMgr.Delete(ctx, policyID)
	return c.policyMgr.Get(ctx, id)
	return c.policyMgr.GetRelatedPolices(ctx, projectID, eventType)
	return c.policyMgr.List(ctx, query)
	return c.policyMgr.Update(ctx, policy)
	return c.taskMgr.Count(ctx, buildTaskQuery(execID, query))
	return c.taskMgr.GetLog(ctx, taskID)
	return c.taskMgr.List(ctx, buildTaskQuery(execID, query))
	return query
	return tasks[0], nil
	return time.Time{}, nil
	taskMgr   task.Manager
	tasks, err := c.taskMgr.List(ctx, query)
	webhookJobVendors = q.NewOrList([]any{job.WebhookJobVendorType, job.SlackJobVendorType, job.AMQPJobVendorType})
	webhookJobVendors = q.NewOrList([]any{job.WebhookJobVendorType, job.SlackJobVendorType, job.DiscordJobVendorType})
	}
	})
)
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Copyright Project Harbor Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// See the License for the specific language governing permissions and
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// you may not use this file except in compliance with the License.
func (c *controller) CountExecutions(ctx context.Context, policyID int64, query *q.Query) (int64, error) {
func (c *controller) CountPolicies(ctx context.Context, query *q.Query) (int64, error) {
func (c *controller) CountTasks(ctx context.Context, execID int64, query *q.Query) (int64, error) {
func (c *controller) CreatePolicy(ctx context.Context, policy *model.Policy) (int64, error) {
func (c *controller) DeletePolicy(ctx context.Context, policyID int64) error {
func (c *controller) GetLastTriggerTime(ctx context.Context, eventType string, policyID int64) (time.Time, error) {
func (c *controller) GetPolicy(ctx context.Context, id int64) (*model.Policy, error) {
func (c *controller) GetRelatedPolices(ctx context.Context, projectID int64, eventType string) ([]*model.Policy, error) {
func (c *controller) GetTask(ctx context.Context, taskID int64) (*task.Task, error) {
func (c *controller) GetTaskLog(ctx context.Context, taskID int64) ([]byte, error) {
func (c *controller) ListExecutions(ctx context.Context, policyID int64, query *q.Query) ([]*task.Execution, error) {
func (c *controller) ListPolicies(ctx context.Context, query *q.Query) ([]*model.Policy, error) {
func (c *controller) ListTasks(ctx context.Context, execID int64, query *q.Query) ([]*task.Task, error) {
func (c *controller) UpdatePolicy(ctx context.Context, policy *model.Policy) error {
func NewController() Controller {
func buildExecutionQuery(policyID int64, query *q.Query) *q.Query {
func buildTaskQuery(execID int64, query *q.Query) *q.Query {
import (
package webhook
type Controller interface {
type controller struct {
var (
}
