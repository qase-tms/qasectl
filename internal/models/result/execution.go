package result

type Execution struct {
	StartTime  *float64 `json:"start_time"`
	EndTime    *float64 `json:"end_time"`
	Status     string   `json:"status"`
	Duration   *float64 `json:"duration"`
	StackTrace *string  `json:"stack_trace"`
	Thread     *string  `json:"thread"`
}
