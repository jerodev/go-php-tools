package laravel

import (
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
