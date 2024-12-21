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

	commandPayload := job.createJobPayload().Data.Command
	if !strings.HasPrefix(commandPayload, `O:27:"App\Events\UpdateDashboards":1:{`) {
		t.Errorf("Invalid serialized data: `%s`", commandPayload)
	}
}

func TestNewBroadcastEvent(t *testing.T) {
	job, _ := NewBroadcastEvent(`App\Events\UpdateDashboards`, job{
		Parts: []string{"foo", "bar"},
	})

	commandPayload := job.createJobPayload().Data.Command
	expected := fmt.Sprintf(`O:%v:"%s":1:{s:5:"event";O:%v:"%s"`, len(broadcastJobClass), broadcastJobClass, len(`App\Events\UpdateDashboards`), `App\Events\UpdateDashboards`)
	if !strings.HasPrefix(commandPayload, expected) {
		t.Errorf("Invalid serialized data `%s`", commandPayload)
	}
}
