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

package trigger

import (
	"testing"

	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/utils/test"
	"github.com/goharbor/harbor/src/replication"
	"github.com/goharbor/harbor/src/replication/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTrigger(t *testing.T) {
	// nil policy
	_, err := createTrigger(nil)
	require.NotNil(t, err)

	// nil trigger
	_, err = createTrigger(&models.ReplicationPolicy{})
	require.NotNil(t, err)

	// schedule trigger
	trigger, err := createTrigger(&models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindSchedule,
			ScheduleParam: &models.ScheduleParam{
				Type:    replication.TriggerScheduleWeekly,
				Weekday: 1,
				Offtime: 1,
			},
		},
	})
	require.Nil(t, err)
	assert.NotNil(t, trigger)

	// immediate trigger
	trigger, err = createTrigger(&models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindImmediate,
		},
	})
	require.Nil(t, err)
	assert.NotNil(t, trigger)

	// manual trigger
	trigger, err = createTrigger(&models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindManual,
		},
	})
	require.Nil(t, err)
	assert.Nil(t, trigger)
}

func TestSetupTrigger(t *testing.T) {
	dao.DefaultDatabaseWatchItemDAO = &test.FakeWatchItemDAO{}

	mgr := NewManager(1)

	err := mgr.SetupTrigger(&models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindImmediate,
		},
	})
	assert.Nil(t, err)
}

func TestUnsetTrigger(t *testing.T) {
	dao.DefaultDatabaseWatchItemDAO = &test.FakeWatchItemDAO{}

	mgr := NewManager(1)

	err := mgr.UnsetTrigger(&models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindImmediate,
		},
	})
	assert.Nil(t, err)
}
