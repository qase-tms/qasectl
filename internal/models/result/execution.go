package result

import "time"

type Execution struct {
	StartTime  *time.Time
	EndTime    *time.Time
	Status     string
	Duration   time.Duration
	StackTrace *string
	Thread     *string
}
