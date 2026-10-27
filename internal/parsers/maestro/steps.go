package maestro

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

var serviceCommands = map[string]bool{
	"defineVariablesCommand":    true,
	"applyConfigurationCommand": true,
}

var attemptLabels = map[string]string{
	"retryCommand":  "Attempt",
	"repeatCommand": "Iteration",
}

type node struct {
	cmd      *command
	children []*node
}

func buildTree(cmds []command) []*node {
	var roots []*node
	var stack []*node
	for i := range cmds {
		n := &node{cmd: &cmds[i]}
		for len(stack) > 0 && stack[len(stack)-1].cmd.Depth >= n.cmd.Depth {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			roots = append(roots, n)
		} else {
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, n)
		}
		stack = append(stack, n)
	}
	return roots
}

func buildSteps(flowDir string, cmds []command, showService bool) []models.Step {
	return nodesToSteps(flowDir, buildTree(cmds), showService)
}

func nodesToSteps(flowDir string, nodes []*node, showService bool) []models.Step {
	var steps []models.Step
	for _, n := range nodes {
		if !showService && serviceCommands[n.cmd.Type] {
			continue
		}
		steps = append(steps, nodeToStep(flowDir, n, showService))
	}
	return steps
}

func nodeToStep(flowDir string, n *node, showService bool) models.Step {
	start := float64(n.cmd.Timestamp)
	end := float64(n.cmd.Timestamp + n.cmd.Duration)
	dur := end - start
	status, comment := stepStatus(*n.cmd)
	step := models.Step{
		Data: models.Data{Action: describe(*n.cmd)},
		Execution: models.StepExecution{
			StartTime:   &start,
			EndTime:     &end,
			Duration:    &dur,
			Status:      status,
			Comment:     comment,
			Attachments: stepAttachments(flowDir, *n.cmd),
		},
	}

	attempts := splitAttempts(n)
	if attempts == nil {
		step.Steps = nodesToSteps(flowDir, n.children, showService)
		return step
	}
	for i, a := range attempts {
		step.Steps = append(step.Steps, attemptStep(flowDir, fmt.Sprintf("%s %d", attemptLabels[n.cmd.Type], i+1), a, showService))
	}
	return step
}

// splitAttempts returns nil when the children cannot be split reliably and must stay flat.
func splitAttempts(n *node) [][]*node {
	if _, ok := attemptLabels[n.cmd.Type]; !ok || len(n.children) == 0 {
		return nil
	}
	var a struct {
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(n.cmd.Args, &a); err != nil || len(a.Commands) == 0 {
		return nil
	}
	var attempts [][]*node
	var cur []*node
	for _, c := range n.children {
		cur = append(cur, c)
		if c.cmd.Status == "FAILED" || len(cur) == len(a.Commands) {
			attempts = append(attempts, cur)
			cur = nil
		}
	}
	if len(cur) > 0 {
		attempts = append(attempts, cur)
	}
	return attempts
}

func attemptStep(flowDir, action string, children []*node, showService bool) models.Step {
	start := children[0].cmd.Timestamp
	end := start
	status := "passed"
	for _, c := range children {
		if e := c.cmd.Timestamp + c.cmd.Duration; e > end {
			end = e
		}
		if s, _ := stepStatus(*c.cmd); s == "failed" {
			status = "failed"
		}
	}
	startF, endF := float64(start), float64(end)
	dur := endF - startF
	return models.Step{
		Data: models.Data{Action: action},
		Execution: models.StepExecution{
			StartTime: &startF,
			EndTime:   &endF,
			Duration:  &dur,
			Status:    status,
		},
		Steps: nodesToSteps(flowDir, children, showService),
	}
}

func stepStatus(c command) (status, comment string) {
	switch c.Status {
	case "COMPLETED":
		return "passed", ""
	case "FAILED":
		if c.Error != nil {
			comment = c.Error.Message
		}
		return "failed", comment
	case "SKIPPED":
		return "skipped", ""
	case "WARNED":
		return "passed", "Optional command failed"
	case "RUNNING", "PENDING":
		return "failed", "Command did not finish"
	default:
		slog.Warn("unknown Maestro command status, marking step skipped", "status", c.Status, "command", c.Type)
		return "skipped", ""
	}
}

// deepestFailure follows the last FAILED command at each level from the top, so a recovered
// retry attempt or an earlier failed attempt never wins over the failure that ended the flow.
func deepestFailure(cmds []command) *command {
	level := buildTree(cmds)
	var found *node
	for {
		var next *node
		for _, n := range level {
			if n.cmd.Status == "FAILED" {
				next = n
			}
		}
		if next == nil {
			break
		}
		found = next
		level = next.children
	}
	if found == nil {
		return nil
	}
	return found.cmd
}

func hadFailedAttempt(cmds []command) bool {
	var walk func([]*node) bool
	walk = func(nodes []*node) bool {
		for _, n := range nodes {
			if n.cmd.Type == "retryCommand" && n.cmd.Status == "COMPLETED" {
				for _, c := range n.children {
					if c.cmd.Status == "FAILED" {
						return true
					}
				}
			}
			if walk(n.children) {
				return true
			}
		}
		return false
	}
	return walk(buildTree(cmds))
}

type selectorArgs struct {
	TextRegex string `json:"textRegex"`
	IDRegex   string `json:"idRegex"`
}

func (s *selectorArgs) String() string {
	switch {
	case s == nil:
		return ""
	case s.TextRegex != "":
		return fmt.Sprintf("%q", s.TextRegex)
	case s.IDRegex != "":
		return fmt.Sprintf("id %q", s.IDRegex)
	}
	return ""
}

type conditionArgs struct {
	Visible         *selectorArgs `json:"visible"`
	NotVisible      *selectorArgs `json:"notVisible"`
	ScriptCondition string        `json:"scriptCondition"`
}

func (c *conditionArgs) String() string {
	switch {
	case c == nil:
		return ""
	case c.Visible.String() != "":
		return "visible " + c.Visible.String()
	case c.NotVisible.String() != "":
		return "not visible " + c.NotVisible.String()
	case c.ScriptCondition != "":
		return "true: " + stripExpr(c.ScriptCondition)
	}
	return ""
}

func stripExpr(s string) string {
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		return s[2 : len(s)-1]
	}
	return s
}

// describe falls back to the bare command type whenever the template's main argument is missing.
func describe(c command) string {
	if c.Label != "" {
		return c.Label
	}
	var a struct {
		Selector          *selectorArgs  `json:"selector"`
		Condition         *conditionArgs `json:"condition"`
		SourceDescription string         `json:"sourceDescription"`
		AppID             string         `json:"appId"`
		Text              string         `json:"text"`
		Path              string         `json:"path"`
		ScriptString      string         `json:"scriptString"`
		Direction         string         `json:"direction"`
		MaxRetries        string         `json:"maxRetries"`
		Times             string         `json:"times"`
	}
	_ = json.Unmarshal(c.Args, &a)

	text := ""
	switch c.Type {
	case "tapOnElement":
		if s := a.Selector.String(); s != "" {
			text = "Tap on " + s
		}
	case "assertConditionCommand":
		if s := a.Condition.String(); s != "" {
			text = "Assert " + s
		}
	case "runFlowCommand":
		switch {
		case a.SourceDescription != "":
			text = "Run flow " + a.SourceDescription
		case a.Condition.String() != "":
			text = "Run flow if " + a.Condition.String()
		default:
			text = "Run flow"
		}
	case "launchAppCommand":
		if a.AppID != "" {
			text = "Launch app " + a.AppID
		}
	case "inputTextCommand":
		if a.Text != "" {
			text = fmt.Sprintf("Input text %q", a.Text)
		}
	case "takeScreenshotCommand":
		if a.Path != "" {
			text = "Take screenshot " + a.Path
		}
	case "evalScriptCommand":
		if a.ScriptString != "" {
			text = "Eval script: " + stripExpr(a.ScriptString)
		}
	case "runScriptCommand":
		text = "Run script"
	case "swipeCommand":
		if a.Direction != "" {
			text = "Swipe " + strings.ToLower(a.Direction)
		}
	case "waitForAnimationToEndCommand":
		text = "Wait for animation to end"
	case "retryCommand":
		if a.MaxRetries != "" {
			text = "Retry up to " + a.MaxRetries + " times"
		}
	case "repeatCommand":
		if a.Times != "" {
			text = "Repeat " + a.Times + " times"
		}
	}
	if text == "" {
		text = strings.TrimSuffix(c.Type, "Command")
	}
	return text
}
