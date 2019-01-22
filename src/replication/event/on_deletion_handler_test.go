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

package event

import (
	"testing"

	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/utils/test"
	"github.com/goharbor/harbor/src/replication/event/notification"
	"github.com/stretchr/testify/assert"
)

func TestHandleOfOnDeletionHandler(t *testing.T) {
	dao.DefaultDatabaseWatchItemDAO = &test.FakeWatchItemDAO{}

	handler := &OnDeletionHandler{}

	assert.NotNil(t, handler.Handle(nil))
	assert.NotNil(t, handler.Handle(map[string]string{}))
	assert.NotNil(t, handler.Handle(struct{}{}))

	assert.Nil(t, handler.Handle(notification.OnDeletionNotification{
		Image: "library/hello-world:latest",
	}))
}

func TestIsStatefulOfOnDeletionHandler(t *testing.T) {
	handler := &OnDeletionHandler{}
	assert.False(t, handler.IsStateful())
}
