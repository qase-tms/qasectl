package result

import "time"

type StepExecution struct {
	StartTime   *time.Time
	EndTime     *time.Time
	Status      string
	Duration    *time.Duration
	Comment     string
	Attachments []Attachment
}
