package result

type StepExecution struct {
	StartTime   *float64     `json:"start_time"`
	EndTime     *float64     `json:"end_time"`
	Status      string       `json:"status"`
	Duration    *float64     `json:"duration"`
	Comment     string       `json:"comment"`
	Attachments []Attachment `json:"attachments"`
}
