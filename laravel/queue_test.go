package laravel

import (
	"strings"
	"testing"
)

type UpdateDashboardPayload struct {
	Parts []string `php:"parts"`
}

func TestSerializeJob(t *testing.T) {
	job, _ := NewQueueJob(`App\Events\UpdateDashboards`, UpdateDashboardPayload{
		Parts: []string{"foo", "bar"},
	})

	commandPayload := job.createJobPayload().Data.Command
	if strings.Contains(commandPayload, "::") {
		t.Error("Invalid serialized data: ", commandPayload)
	}
}
