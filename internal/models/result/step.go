package result

type Step struct {
	Data      Data          `json:"data"`
	Execution StepExecution `json:"execution"`
	Steps     []Step        `json:"steps"`
}
