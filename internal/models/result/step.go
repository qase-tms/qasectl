package result

type Step struct {
	Data      Data
	Execution StepExecution
	Steps     []Step
}
