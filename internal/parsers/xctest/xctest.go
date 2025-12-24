package xctest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	models "github.com/qase-tms/qasectl/internal/models/result"
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
	layoutTime        = "2006-01-02T15:04:05.000-0700"
	caseId            *int64
	failures          map[string]FailureSummary
	processedFailures []string
)

const (
	internalStep   = "com.apple.dt.xctest.activity-type.internal"
	deleteStep     = "com.apple.dt.xctest.activity-type.deletedAttachment"
	attachmentStep = "com.apple.dt.xctest.activity-type.attachmentContainer"
	qaseConfig     = "Qase config"
	heicExt        = ".heic"
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

	tests := p.getXCTests(testPlanSums)

	logger.Debug("got tests", "xctest", tests)

	results := p.getResults(tests)

	logger.Debug("got results", "results", results)

	return results, nil
}

func (p *Parser) readJson(id *string) ([]byte, error) {
	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "json"}
	if id != nil {
		args = append(args, "--id", *id)
	}

	executeCommand := func(args []string) ([]byte, error) {
		out, err := exec.Command("xcrun", args...).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get XCResult: %w", err)
		}
		return out, nil
	}

	out, err := executeCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 64") {
			args = append(args, "--legacy")
			out, err = executeCommand(args)
			if err != nil {
				return nil, fmt.Errorf("failed to get XCResult with --legacy: %w", err)
			}
			return out, nil
		}
		return nil, err
	}

	return out, nil
}

func (p *Parser) readAttachment(id string) ([]byte, error) {
	const op = "xctest.Parser.readAttachment"
	logger := slog.With("op", op, "attachmentId", id)

	logger.Debug("starting to read attachment", "path", p.path)

	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "raw", "--id", id}
	logger.Debug("executing xcresulttool command", "args", args)

	executeCommand := func(args []string) ([]byte, error) {
		out, err := exec.Command("xcrun", args...).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get XCResult attachment: %w", err)
		}
		logger.Debug("successfully executed xcrun command", "outputSize", len(out))
		return out, nil
	}

	out, err := executeCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 64") {
			args = append(args, "--legacy")
			out, err = executeCommand(args)
			if err != nil {
				return nil, fmt.Errorf("failed to get XCResult attachment with --legacy: %w", err)
			}
			logger.Debug("successfully got attachment with --legacy flag", "size", len(out))
		} else {
			logger.Error("failed to get attachment", "error", err)
			return nil, err
		}
	} else {
		logger.Debug("successfully got attachment", "size", len(out))
	}

	return out, nil
}

func (p *Parser) getTestPlanSumIDs() (map[string]TestMeta, error) {
	out, err := p.readJson(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	var structure Structure
	err = json.Unmarshal(out, &structure)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal test plan summaries: %w", err)
	}

	testPlanSumIDs := make(map[string]TestMeta)

	for _, action := range structure.Actions.Values {
		if action.ActionResult.TestsRef == nil {
			continue
		}

		m := TestMeta{}

		if action.StartedTime != nil {
			st := float64(parseTime(action.StartedTime.Value).UnixMilli())
			m.StartTime = &st
		}
		if action.EndedTime != nil {
			et := float64(parseTime(action.EndedTime.Value).UnixMilli())
			m.EndTime = &et
		}
		if action.RunDestination != nil {
			m.Device = action.RunDestination.DisplayName.Value
		}

		testPlanSumIDs[action.ActionResult.TestsRef.ID.Value] = m
	}

	return testPlanSumIDs, nil
}

func (p *Parser) getTestPlanSums(IDs map[string]TestMeta) (map[TestMeta][]ActionTestPlanRunSummaries, error) {
	testPlanSums := make(map[TestMeta][]ActionTestPlanRunSummaries)

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
	const op = "xctest.Parser.getActionTestSummary"
	logger := slog.With("op", op, "testId", ID)

	logger.Debug("getting action test summary")
	out, err := p.readJson(&ID)
	if err != nil {
		logger.Error("failed to read JSON", "error", err)
		return ActionTestSummary{}, fmt.Errorf("failed to get action test summary: %w", err)
	}

	logger.Debug("successfully read JSON", "size", len(out))
	var summary ActionTestSummary
	err = json.Unmarshal(out, &summary)
	if err != nil {
		logger.Error("failed to unmarshal JSON", "error", err)
		return ActionTestSummary{}, fmt.Errorf("failed to unmarshal action test summary: %w", err)
	}

	logger.Debug("successfully parsed action test summary", "testStatus", summary.TestStatus.Value, "hasFailureSummaries", summary.FailureSummaries != nil, "activitySummariesCount", len(summary.ActivitySummaries.Values))

	if summary.FailureSummaries != nil {
		logger.Debug("failure summaries details", "count", len(summary.FailureSummaries.Values))
		for i, f := range summary.FailureSummaries.Values {
			logger.Debug("failure summary", "index", i, "uuid", f.UUID.Value, "hasAttachments", f.Attachments != nil)
			if f.Attachments != nil {
				logger.Debug("failure attachments", "count", len(f.Attachments.Values))
			}
		}
	}

	return summary, nil
}

func (p *Parser) getXCTests(testPlanSums map[TestMeta][]ActionTestPlanRunSummaries) []XCTest {
	tests := make([]XCTest, 0)

	for k, v := range testPlanSums {
		for _, testPlanSum := range v {
			for _, testPlanRunSum := range testPlanSum.Summaries.Values {
				k.Configuration = testPlanRunSum.Name.Value
				for _, testSum := range testPlanRunSum.TestableSummaries.Values {
					for _, test := range testSum.Tests.Values {
						if test.Subtests != nil {
							tt := p.getXCTestsFromSubtest(*test.Subtests, make([]string, 0))
							for _, t := range tt {
								t.Metadata = k
								tests = append(tests, t)
							}
						}
					}
				}
			}
		}
	}

	return tests
}

func (p *Parser) getXCTestsFromSubtest(s Subtests, suites []string) []XCTest {
	const op = "xctest.Parser.getXCTestsFromSubtest"
	logger := slog.With("op", op)

	var tests []XCTest

	for _, v := range s.Values {
		if v.SummaryRef != nil {
			act, err := p.getActionTestSummary(v.SummaryRef.ID.Value)
			if err != nil {
				logger.Error("failed to get action test summary", "error", err, "ID", v.SummaryRef.ID.Value)
				continue
			}

			d, err := strconv.ParseFloat(v.Duration.Value, 64)
			if err != nil {
				logger.Error("failed to parse duration", "error", err)
				d = 0
			}

			test := XCTest{
				Name:      v.Name.Value,
				Action:    act,
				Suites:    make([]string, 0),
				Signature: v.IdentifierURL.Value,
				Duration:  d,
			}

			if len(suites) > 0 {
				test.Suites = append(test.Suites, suites...)
			}

			tests = append(tests, test)

			continue
		}

		if v.Subtests != nil {
			// Create a new slice to avoid suite nesting issues
			newSuites := make([]string, len(suites))
			copy(newSuites, suites)
			newSuites = append(newSuites, v.Name.Value)

			tt := p.getXCTestsFromSubtest(*v.Subtests, newSuites)
			tests = append(tests, tt...)
		}
	}

	return tests
}

func (p *Parser) getResults(tests []XCTest) []models.Result {
	const op = "xctest.Parser.getResults"
	logger := slog.With("op", op)

	logger.Debug("starting to process tests", "testCount", len(tests))
	var results []models.Result

	for i, t := range tests {
		caseId = nil
		failures = make(map[string]FailureSummary)
		processedFailures = make([]string, 0)

		logger := logger.With("testTitle", t.Name, "testIndex", i)
		logger.Debug("processing test", "test", t)

		suites := make([]models.SuiteData, 0)
		suites = append(suites, models.SuiteData{Title: t.Metadata.Suite})
		for _, s := range t.Suites {
			suites = append(suites, models.SuiteData{
				Title: s,
			})
		}

		duration := t.Duration
		result := models.Result{
			Title:     t.Name,
			Signature: &t.Signature,
			Execution: models.Execution{
				Status:    getStatus(t.Action.TestStatus.Value),
				Duration:  &duration,
			},
			Fields:      map[string]string{},
			Attachments: make([]models.Attachment, 0),
			Params: map[string]string{
				"Device":        t.Metadata.Device,
				"Configuration": t.Metadata.Configuration,
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: suites,
				},
			},
			StepType: "text",
			Muted:    false,
		}

		if len(t.Action.ActivitySummaries.Values) > 0 {
			startTime := t.Action.ActivitySummaries.Values[0].Start.Value
			endTime := t.Action.ActivitySummaries.Values[len(t.Action.ActivitySummaries.Values)-1].Finish.Value
			startTimeFloat := float64(parseTime(startTime).UnixMilli())
			endTimeFloat := float64(parseTime(endTime).UnixMilli())
			result.Execution.StartTime = &startTimeFloat
			result.Execution.EndTime = &endTimeFloat
		}

		if t.Action.FailureSummaries != nil {
			var message string
			for _, f := range t.Action.FailureSummaries.Values {
				failures[f.UUID.Value] = f
				message += f.Message.Value + "\n"
			}
			result.Execution.StackTrace = &message
		} else {
			logger.Debug("test has no failure summaries")
		}

		steps, _, a := p.getSteps(t.Action.ActivitySummaries, 0)
		result.Attachments = append(result.Attachments, a...)

		// Search for attachments in activitySummaries at test level
		testLevelAttachments := p.getTestLevelAttachments(t.Action.ActivitySummaries)
		result.Attachments = append(result.Attachments, testLevelAttachments...)

		result.Steps = steps
		if caseId != nil {
			result.TestOpsID = caseId
		}

		for k, v := range failures {
			if isFailureProcessed(k) {
				logger.Debug("skipping already processed failure", "failureId", k)
				continue
			}

			logger.Debug("processing failure attachments", "failureId", k, "hasAttachments", v.Attachments != nil)
			if v.Attachments != nil {
				logger.Debug("found attachments in failure", "attachmentCount", len(v.Attachments.Values))
				for i, att := range v.Attachments.Values {
					logger.Debug("processing failure attachment", "index", i, "filename", att.Filename.Value, "name", att.Name.Value)
					a, err := p.getAttachment(att)
					if err != nil {
						logger.Error("failed to get attachments", "error", err, "attachmentIndex", i)
						continue
					}

					logger.Debug("successfully added failure attachment", "attachmentName", a.Name)
					result.Attachments = append(result.Attachments, a)
				}
			} else {
				logger.Debug("no attachments found in failure", "failureId", k)
			}
		}

		logger.Debug("completed test processing", "attachmentCount", len(result.Attachments), "stepCount", len(result.Steps))
		results = append(results, result)
	}

	logger.Debug("completed all tests processing", "totalResults", len(results))
	return results
}

func isFailureProcessed(f string) bool {
	for _, p := range processedFailures {
		if f == p {
			return true
		}
	}

	return false
}

// detectFileExtension determines the file extension based on file content
func (p *Parser) detectFileExtension(data []byte) string {
	if len(data) < 4 {
		return ""
	}

	// Check for common image formats
	switch {
	case len(data) >= 8 && string(data[0:8]) == "\x89PNG\r\n\x1a\n":
		return ".png"
	case len(data) >= 2 && string(data[0:2]) == "\xff\xd8":
		return ".jpg"
	case len(data) >= 6 && string(data[0:6]) == "GIF87a" || string(data[0:6]) == "GIF89a":
		return ".gif"
	case len(data) >= 4 && string(data[0:4]) == "RIFF" && len(data) >= 8 && string(data[8:12]) == "WEBP":
		return ".webp"
	case len(data) >= 4 && string(data[0:4]) == "\x00\x00\x01\x00":
		return ".ico"
	case len(data) >= 4 && string(data[0:4]) == "BM\x00\x00":
		return ".bmp"
	case len(data) >= 4 && string(data[0:4]) == "\x00\x00\x00\x20":
		return ".tiff"
	case len(data) >= 4 && string(data[0:4]) == "II*\x00":
		return ".tiff"
	case len(data) >= 4 && string(data[0:4]) == "MM\x00*":
		return ".tiff"
	// Check for HEIC/HEIF
	case len(data) >= 12 && string(data[4:12]) == "ftypheic":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftypheix":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftypheis":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevc":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevx":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevs":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevm":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftypheis":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftypheix":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftypheic":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevc":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevx":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevs":
		return ".heic"
	case len(data) >= 12 && string(data[4:12]) == "ftyphevm":
		return ".heic"
	// Check for PDF
	case len(data) >= 4 && string(data[0:4]) == "%PDF":
		return ".pdf"
	// Check for Binary Property List (bplist00)
	case len(data) >= 8 && string(data[0:8]) == "bplist00":
		return ".plist"
	// Check for text files (UTF-8, UTF-16, ASCII)
	case len(data) >= 3 && string(data[0:3]) == "\xef\xbb\xbf":
		return ".txt" // UTF-8 BOM
	case len(data) >= 2 && string(data[0:2]) == "\xff\xfe":
		return ".txt" // UTF-16 LE BOM
	case len(data) >= 2 && string(data[0:2]) == "\xfe\xff":
		return ".txt" // UTF-16 BE BOM
	// Check for JSON
	case len(data) > 0 && (data[0] == '{' || data[0] == '['):
		// Try to parse as JSON to confirm
		var temp interface{}
		if json.Unmarshal(data, &temp) == nil {
			return ".json"
		}
	// Check for XML
	case len(data) >= 5 && string(data[0:5]) == "<?xml":
		return ".xml"
	case len(data) >= 1 && data[0] == '<':
		return ".xml"
	// Check for plain text (if all bytes are printable ASCII)
	default:
		isText := true
		for _, b := range data {
			if b < 32 && b != 9 && b != 10 && b != 13 { // tab, newline, carriage return
				isText = false
				break
			}
		}
		if isText {
			return ".txt"
		}
	}

	return ""
}

func (p *Parser) getAttachment(a Attachment) (models.Attachment, error) {
	const op = "xctest.Parser.getAttachment"
	logger := slog.With("op", op, "filename", a.Filename.Value, "name", a.Name.Value, "payloadId", a.PayloadRef.ID.Value)

	logger.Debug("processing attachment", "originalFilename", a.Filename.Value)

	// Handle HEIC conversion
	a.Filename.Value = strings.Replace(a.Filename.Value, heicExt, ".jpeg", 1)

	att, err := p.readAttachment(a.PayloadRef.ID.Value)
	if err != nil {
		logger.Error("failed to read attachment content", "error", err)
		return models.Attachment{}, fmt.Errorf("failed to get attachment: %w", err)
	}

	// Determine file extension based on content if filename doesn't have one
	filename := a.Filename.Value
	if !strings.Contains(filename, ".") || strings.HasSuffix(filename, ".") {
		extension := p.detectFileExtension(att)
		if extension != "" {
			filename = filename + extension
			logger.Debug("added file extension based on content", "filename", filename, "extension", extension)
		}
	}

	return models.Attachment{
		Name:    filename,
		Content: &att,
	}, nil
}

func (p *Parser) getSteps(as ActivitySummaries, level int) ([]models.Step, bool, []models.Attachment) {
	const op = "xctest.Parser.getSteps"
	logger := slog.With("op", op, "level", level)

	logger.Debug("processing activity summaries", "activityCount", len(as.Values))
	steps := make([]models.Step, 0, len(as.Values))
	attachments := make([]models.Attachment, 0)
	isFailedStep := false

	for i, v := range as.Values {
		logger.Debug("processing activity", "index", i, "activityType", v.ActivityType.Value, "title", v.Title.Value, "hasAttachments", v.Attachments != nil, "hasSubactivities", v.Subactivities != nil, "hasFailureSummaryIDs", v.FailureSummaryIDs != nil)

		if v.Attachments != nil {
			logger.Debug("activity has attachments", "attachmentCount", len(v.Attachments.Values))
			for j, att := range v.Attachments.Values {
				logger.Debug("attachment details", "index", j, "filename", att.Filename.Value, "name", att.Name.Value, "payloadId", att.PayloadRef.ID.Value)
			}
		}

		stepAttachments := make([]models.Attachment, 0)
		stepComment := ""
		isStepFailed := false
		childSteps := make([]models.Step, 0)
		if v.ActivityType.Value == attachmentStep && v.Attachments != nil &&
			len(v.Attachments.Values) == 1 && v.Attachments.Values[0].Name.Value == qaseConfig {
			logger.Debug("processing Qase config attachment", "filename", v.Attachments.Values[0].Filename.Value)
			att, err := p.getAttachment(v.Attachments.Values[0])
			if err != nil {
				logger.Error("failed to get Qase config attachment", "error", err)
				continue
			}

			logger.Debug("parsing Qase ID from attachment", "contentSize", len(*att.Content))
			var ID QaseId
			err = json.Unmarshal(*att.Content, &ID)
			if err != nil {
				logger.Error("failed to unmarshal Qase ID", "error", err)
				continue
			}

			logger.Debug("successfully parsed Qase ID", "caseId", ID.ID)
			if caseId == nil {
				caseId = &ID.ID
				logger.Debug("set case ID", "caseId", *caseId)
			} else {
				logger.Debug("case ID already set", "existingCaseId", *caseId, "newCaseId", ID.ID)
			}

			continue
		}

		if v.FailureSummaryIDs != nil {
			for _, f := range v.FailureSummaryIDs.Values {
				if fm, ok := failures[f.Value]; ok {
					processedFailures = append(processedFailures, f.Value)
					isStepFailed = true
					stepComment = fmt.Sprintf("Line: %s\n%s", fm.LineNumber.Value, fm.Message.Value)
					if fm.Attachments != nil {
						for _, a := range fm.Attachments.Values {
							att, err := p.getAttachment(a)
							if err != nil {
								logger.Error("failed to get failure summary attachment", "error", err, "filename", a.Filename.Value)
								continue
							}
							stepAttachments = append(stepAttachments, att)
						}
					}
				}
			}

			if v.Attachments != nil {
				for _, a := range v.Attachments.Values {
					att, err := p.getAttachment(a)
					if err != nil {
						logger.Error("failed to get step attachment", "error", err, "filename", a.Filename.Value)
						continue
					}
					stepAttachments = append(stepAttachments, att)
				}
			}
		}

		if v.Subactivities != nil {
			cs, f, a := p.getSteps(*v.Subactivities, level+1)
			if f {
				isStepFailed = true
			}
			childSteps = append(childSteps, cs...)

			if len(cs) == 0 {
				stepAttachments = append(stepAttachments, a...)
			}
		}

		isActivityTypeInternalOrDelete := v.ActivityType.Value == internalStep || v.ActivityType.Value == deleteStep
		isLevelNotAll := p.level != All
		isLevelZeroOrAboveFirst := level == 0 || (level > 0 && p.level != FirstLevel)

		if isActivityTypeInternalOrDelete && isLevelNotAll && isLevelZeroOrAboveFirst {
			attachments = append(attachments, stepAttachments...)

			if isStepFailed {
				isFailedStep = true
			}

			continue
		}

		logger := logger.With("stepAction", v.Title.Value)
		logger.Debug("processing step", "step", v)

		step := models.Step{
			Data: models.Data{
				Action: v.Title.Value,
			},
			Execution: models.StepExecution{
				Status:      "passed",
				Attachments: stepAttachments,
				Comment:     stepComment,
			},
			Steps: childSteps,
		}

		if v.Start.Value != "" {
			st := float64(parseTime(v.Start.Value).Unix())
			step.Execution.StartTime = &st
		}

		if v.Finish.Value != "" {
			et := float64(parseTime(v.Finish.Value).Unix())
			step.Execution.EndTime = &et
		}

		if isStepFailed {
			step.Execution.Status = "failed"
			isFailedStep = true
		}

		steps = append(steps, step)
	}

	logger.Debug("completed steps processing", "stepCount", len(steps), "attachmentCount", len(attachments), "hasFailedStep", isFailedStep)
	return steps, isFailedStep, attachments
}

// getTestLevelAttachments searches for attachments in activitySummaries at the test level
func (p *Parser) getTestLevelAttachments(as ActivitySummaries) []models.Attachment {
	const op = "xctest.Parser.getTestLevelAttachments"
	logger := slog.With("op", op)

	attachments := make([]models.Attachment, 0)

	for _, activity := range as.Values {
		if activity.Attachments != nil {
			for _, att := range activity.Attachments.Values {
				attachment, err := p.getAttachment(att)
				if err != nil {
					logger.Error("failed to get test-level attachment", "error", err, "filename", att.Filename.Value)
					continue
				}

				attachments = append(attachments, attachment)
			}
		}
	}

	return attachments
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

func parseTime(s string) time.Time {
	t, err := time.Parse(layoutTime, s)
	if err != nil {
		slog.Error("failed to parse time", "error", err, "time", s)
		return time.Time{}
	}

	return t
}
