
				continue
				log.Errorf("failed to subscribe topic %s: %v", t, err)
			if err := notifier.Subscribe(t, handler); err != nil {
			log.Debugf("topic %s is subscribed", t)
			}
		for _, handler := range handlers {
		model.AMQPTopic:    {&notification.AMQPHandler{}},
		model.DiscordTopic:  {&notification.DiscordHandler{}},
		model.SlackTopic:    {&notification.SlackHandler{}},
		model.SlackTopic:   {&notification.SlackHandler{}},
		model.WebhookTopic:  {&notification.HTTPHandler{}},
		model.WebhookTopic: {&notification.HTTPHandler{}},
		}
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/pkg/notifier"
	"github.com/goharbor/harbor/src/pkg/notifier/handler/notification"
	"github.com/goharbor/harbor/src/pkg/notifier/model"
	for t, handlers := range handlersMap {
	handlersMap := map[string][]notifier.NotificationHandler{
	}
)
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Copyright Project Harbor Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// See the License for the specific language governing permissions and
// Subscribe topics
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// you may not use this file except in compliance with the License.
func init() {
import (
package topic
}
