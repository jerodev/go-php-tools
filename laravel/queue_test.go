package laravel

import (
	"fmt"
	"strings"
	"testing"
)

type job struct {
	Parts []string `php:"parts"`
}

func TestSerializeJob(t *testing.T) {
	job, _ := NewQueueJob(`App\Events\UpdateDashboards`, job{
		Parts: []string{"foo", "bar"},
	})

	commandPayload, err := job.createJobPayload()
	if err != nil {
		t.Error("Error in job payload:", err)
	}

	if !strings.HasPrefix(commandPayload.Data.Command, `O:27:"App\Events\UpdateDashboards":1:{`) {
		t.Errorf("Invalid serialized data: `%s`", commandPayload.Data.Command)
	}
}

func TestNewBroadcastEvent(t *testing.T) {
	job, _ := NewBroadcastEvent(`App\Events\UpdateDashboards`, job{
		Parts: []string{"foo", "bar"},
	})

	commandPayload, err := job.createJobPayload()
	if err != nil {
		t.Error("Error in job payload:", err)
	}

	expected := fmt.Sprintf(`O:%v:"%s":1:{s:5:"event";O:%v:"%s"`, len(broadcastJobClass), broadcastJobClass, len(`App\Events\UpdateDashboards`), `App\Events\UpdateDashboards`)
	if !strings.HasPrefix(commandPayload.Data.Command, expected) {
		t.Errorf("Invalid serialized data `%s`", commandPayload.Data.Command)
	}
}
