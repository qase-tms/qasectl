package xctest

import (
	"encoding/json"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Parser is a parser for XCTest files
type Parser struct {
	path  string
	level stepLevel
}

// NewParser creates a new Parser
func NewParser(path, level string) (*Parser, error) {
	if !strings.HasSuffix(path, ".xcresult") {
		return nil, fmt.Errorf("unsupported format: %s", path)
	}

	return &Parser{
		path:  path,
		level: parseStepLevel(level),
	}, nil
}

var (
	layoutTime = "2006-01-02T15:04:05.000-0700"
	caseId     *int64
)

const (
	internalStep   = "com.apple.dt.xctest.activity-type.internal"
	deleteStep     = "com.apple.dt.xctest.activity-type.deletedAttachment"
	attachmentStep = "com.apple.dt.xctest.activity-type.attachmentContainer"
	qaseConfig     = "Qase config"
)

// Parse parses the XCTest file and returns the results
func (p *Parser) Parse() ([]models.Result, error) {
	const op = "xctest.Parser.Parse"
	logger := slog.With("op", op)

	logger.Debug("parsing XCTest file", "path", p.path)

	testPlanSumIDs, err := p.getTestPlanSumIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	logger.Debug("got test plan summary IDs", "testPlanSumIDs", testPlanSumIDs)

	testPlanSums, err := p.getTestPlanSums(testPlanSumIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	logger.Debug("got test plan summaries", "testPlanSums", testPlanSums)

	results := p.getTests(testPlanSums)

	return results, nil
}

func (p *Parser) readJson(id *string) ([]byte, error) {
	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "json"}
	if id != nil {
		args = append(args, "--id", *id)
	}
	out, err := exec.Command("xcrun", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get XCResult: %w", err)
	}

	return out, nil
}

func (p *Parser) getTestPlanSumIDs() (map[string]string, error) {
	out, err := p.readJson(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	var structure Structure
	err = json.Unmarshal(out, &structure)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal test plan summaries: %w", err)
	}

	testPlanSumIDs := make(map[string]string)

	for _, action := range structure.Actions.Values {
		if action.ActionResult.TestsRef != nil {
			testPlanSumIDs[action.ActionResult.TestsRef.ID.Value] = action.RunDestination.DisplayName.Value
		}
	}

	return testPlanSumIDs, nil
}

func (p *Parser) getTestPlanSums(IDs map[string]string) (map[string][]ActionTestPlanRunSummaries, error) {
	testPlanSums := make(map[string][]ActionTestPlanRunSummaries)

	for k, v := range IDs {
		out, err := p.readJson(&k)
		if err != nil {
			return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
		}
		var sum ActionTestPlanRunSummaries
		err = json.Unmarshal(out, &sum)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal test plan summaries: %w", err)
		}

		testPlanSums[v] = append(testPlanSums[v], sum)
	}

	return testPlanSums, nil
}

func (p *Parser) getActionTestSummary(ID string) (ActionTestSummary, error) {
	out, err := p.readJson(&ID)
	if err != nil {
		return ActionTestSummary{}, fmt.Errorf("failed to get action test summary: %w", err)
	}

	var summary ActionTestSummary
	err = json.Unmarshal(out, &summary)
	if err != nil {
		return ActionTestSummary{}, fmt.Errorf("failed to unmarshal action test summary: %w", err)
	}

	return summary, nil
}

func (p *Parser) getTests(testPlanSums map[string][]ActionTestPlanRunSummaries) []models.Result {
	const op = "xctest.Parser.getTests"
	logger := slog.With("op", op)

	var results []models.Result

	for k, v := range testPlanSums {
		for _, testPlanSum := range v {
			for _, testPlanRunSum := range testPlanSum.Summaries.Values {
				for _, testSum := range testPlanRunSum.TestableSummaries.Values {
					for _, test := range testSum.Tests.Values {
						for _, subTest := range test.Subtests.Values {
							for _, subTest2 := range subTest.Subtests.Values {
								for _, subTest3 := range subTest2.Subtests.Values {
									caseId = nil

									logger := logger.With("testTitle", subTest3.Name.Value)
									logger.Debug("processing test", "test", subTest3)

									act, err := p.getActionTestSummary(subTest3.SummaryRef.ID.Value)
									if err != nil {
										logger.Error("failed to get action test summary", "error", err)
										continue
									}

									d, err := strconv.ParseFloat(subTest3.Duration.Value, 64)
									if err != nil {
										logger.Error("failed to parse duration", "error", err)
										d = 0
									}

									result := models.Result{
										Title:     subTest3.Name.Value,
										Signature: &subTest3.IdentifierURL.Value,
										Execution: models.Execution{
											Status:   getStatus(act.TestStatus.Value),
											Duration: time.Duration(d * float64(time.Second)),
										},
										Fields:      map[string]string{},
										Attachments: make([]models.Attachment, 0),
										Params: map[string]string{
											"Device": k,
										},
										Relations: models.Relation{
											Suite: models.Suite{
												Data: []models.SuiteData{
													{
														Title: subTest.Name.Value,
													},
													{
														Title: subTest2.Name.Value,
													},
												},
											},
										},
										StepType: "text",
										Muted:    false,
									}

									if act.FailureSummaries != nil {
										var message string
										for _, f := range act.FailureSummaries.Values {
											message += f.Message.Value + "\n"
										}

										result.Execution.StackTrace = &message

										attachments := make([]models.Attachment, 0)
										for _, a := range act.FailureSummaries.Values {
											if a.Attachments != nil {
												for _, att := range a.Attachments.Values {
													logger.Debug("processing attachment", "attachment", att)
													c, err := p.getAttachment(att.PayloadRef.ID.Value)
													if err != nil {
														logger.Error("failed to get attachments", "error", err)
													}

													attachments = append(attachments, models.Attachment{
														Name:    att.Filename.Value,
														Content: &c,
													})
												}
											}
										}

										if len(attachments) > 0 {
											result.Attachments = append(result.Attachments, attachments...)
										}
									}

									steps := p.getSteps(act.ActivitySummaries)

									result.Steps = steps
									if caseId != nil {
										result.TestOpsID = caseId
									}

									results = append(results, result)
								}
							}
						}
					}

				}
			}
		}
	}

	return results
}

func (p *Parser) getAttachment(ID string) ([]byte, error) {
	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "raw", "--id", ID}

	out, err := exec.Command("xcrun", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments from XCResult: %w", err)
	}

	return out, nil
}

func (p *Parser) getSteps(as ActivitySummaries) []models.Step {
	const op = "xctest.Parser.getSteps"
	logger := slog.With("op", op)

	steps := make([]models.Step, 0, len(as.Values))

	for _, v := range as.Values {
		if p.level != All && (v.ActivityType.Value == internalStep || v.ActivityType.Value == deleteStep) {
			continue
		}

		logger := logger.With("stepTitle", v.Title.Value)
		logger.Debug("processing step", "step", v)

		if v.ActivityType.Value == attachmentStep && len(v.Attachments.Values) == 1 &&
			v.Attachments.Values[0].Name.Value == qaseConfig {
			att, err := p.getAttachment(v.Attachments.Values[0].PayloadRef.ID.Value)
			if err != nil {
				logger.Error("failed to get attachments", "error", err)
			}

			var ID QaseId
			err = json.Unmarshal(att, &ID)
			if err != nil {
				logger.Error("failed to unmarshal Qase ID", "error", err)
			}

			if caseId == nil {
				caseId = &ID.ID
			}

			continue
		}

		step := models.Step{
			Data: models.Data{
				Action: v.Title.Value,
			},
			Execution: models.StepExecution{
				Status:      "passed",
				Attachments: make([]models.Attachment, 0),
			},
		}

		if v.Start.Value != "" {
			st, err := time.Parse(layoutTime, v.Start.Value)
			if err != nil {
				logger.Error("failed to parse start time", "error", err)
			}
			step.Execution.StartTime = &st
		}

		if v.Finish.Value != "" {
			et, err := time.Parse(layoutTime, v.Finish.Value)
			if err != nil {
				logger.Error("failed to parse finish time", "error", err)
			}
			step.Execution.EndTime = &et
		}

		if v.Attachments != nil {
			for _, a := range v.Attachments.Values {
				att, err := p.getAttachment(a.PayloadRef.ID.Value)
				if err != nil {
					logger.Error("failed to get attachments", "error", err)
				}

				step.Execution.Attachments = append(step.Execution.Attachments, models.Attachment{
					Name:    a.Filename.Value,
					Content: &att,
				})

				if a.Name.Value == "Failed Image" {
					step.Execution.Status = "failed"
				}
			}
		}

		if v.Subactivities != nil {
			cs := p.getChildSteps(*v.Subactivities)

			step.Steps = cs
		}

		steps = append(steps, step)
	}

	return steps
}

func (p *Parser) getChildSteps(s Subactivities) []models.Step {
	const op = "xctest.Parser.getChildSteps"
	logger := slog.With("op", op)

	steps := make([]models.Step, 0, len(s.Values))
	for _, v := range s.Values {
		if p.level != All && p.level != FirstLevel &&
			(v.ActivityType.Value == internalStep || v.ActivityType.Value == deleteStep) {
			continue
		}

		logger := logger.With("stepTitle", v.Title.Value)
		logger.Debug("processing step", "step", v)

		step := models.Step{
			Data: models.Data{
				Action: v.Title.Value,
			},
			Execution: models.StepExecution{
				Status: "passed",
			},
		}

		if v.Start.Value != "" {
			st, err := time.Parse(layoutTime, v.Start.Value)
			if err != nil {
				logger.Error("failed to parse start time", "error", err)
			}
			step.Execution.StartTime = &st
		}

		if v.Finish.Value != "" {
			et, err := time.Parse(layoutTime, v.Finish.Value)
			if err != nil {
				logger.Error("failed to parse finish time", "error", err)
			}
			step.Execution.EndTime = &et
		}

		if v.Subactivities != nil {
			cs := p.getChildSteps(*v.Subactivities)

			step.Steps = cs
		}

		if v.Attachments != nil {
			for _, a := range v.Attachments.Values {
				logger.Debug("processing attachment", "attachment", a)

				att, err := p.getAttachment(a.PayloadRef.ID.Value)
				if err != nil {
					logger.Error("failed to get attachments", "error", err, "ID", a.PayloadRef.ID.Value)
				}

				step.Execution.Attachments = append(step.Execution.Attachments, models.Attachment{
					Name:    a.Filename.Value,
					Content: &att,
				})

				if a.Name.Value == "Failed Image" {
					step.Execution.Status = "failed"
				}
			}
		}

		steps = append(steps, step)
	}

	return steps
}

func getStatus(s string) string {
	switch s {
	case "Success":
		return "passed"
	case "Failure":
		return "failed"
	case "Error":
		return "invalid"
	case "Skipped":
		return "skipped"
	default:
		return "passed"
	}
}
