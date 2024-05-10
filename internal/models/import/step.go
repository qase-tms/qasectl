package _import

type Step struct {
	Data      Data
	Execution StepExecution
	Steps     []Step
}
