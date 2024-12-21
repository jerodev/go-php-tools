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
	if strings.HasSuffix(commandPayload, `O:27:"App\Events\UpdateDashboards":1:{`) {
		t.Error("Invalid serialized data: ", commandPayload)
	}
}

func TestNewBroadcastEvent(t *testing.T) {
	job, _ := NewQueueJob(`App\Events\UpdateDashboards`, job{
		Parts: []string{"foo", "bar"},
	})

	commandPayload := job.createJobPayload().Data.Command
	expected := fmt.Sprintf(`O:%v:"%s":1:{s:5:"Event";O:%v:"%s"`, len(broadcastJobClass), broadcastJobClass, len(`App\Events\UpdateDashboards`), `App\Events\UpdateDashboards`)
	if strings.HasSuffix(commandPayload, expected) {
		t.Error("Invalid serialized data: ", commandPayload)
	}
}
