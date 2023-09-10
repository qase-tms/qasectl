package internal

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qase-tms/qasectl/internal/xcresult"
	"github.com/qase-tms/qasectl/internal/xcresult/model"
	"github.com/qase-tms/qasectl/pkg"
)

func Import(xcPath string) error {
	json, err := xcresult.ToJson(xcPath, nil)
	if err != nil {
		return errors.Wrap(err, "xcresult to json")
	}

	record := xcresult.DecodeObject[model.ActionsInvocationRecord](json)
	tests, err := extractTests(xcPath, record)
	if err != nil {
		return err
	}

	var caseIDs []string
	for _, test := range tests {
		caseID := findCaseID(xcPath, test)
		if caseID == nil {
			continue
		}

		caseIDs = append(caseIDs, *caseID)
	}

	return nil
}

func findCaseID(xcPath string, summary model.ActionTestSummary) *string {
	config := extractQaseConfig(xcPath, summary)
	if config == nil {
		return nil
	}

	return pkg.Ptr(config["case_id"].(string))
}

func extractQaseConfig(xcPath string, summary model.ActionTestSummary) map[string]any {
	for _, activitySummary := range summary.ActivitySummaries {
		binaryData, err := extractAttachmentFromActivitySummary(xcPath, activitySummary, QaseConfigFile)
		if err != nil {
			panic(err)
		}
		if binaryData == nil {
			return nil
		}

		var m map[string]any
		err = json.Unmarshal(binaryData, &m)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func extractAttachmentFromActivitySummary(xcPath string, summary model.ActionTestActivitySummary, attachmentName string) ([]byte, error) {
	if summary.ActivityType != model.ActivityAttachmentContainer {
		return nil, nil
	}

	for _, attachment := range summary.Attachments {
		if pkg.Deref(attachment.Name) == attachmentName && attachment.PayloadRef != nil {
			return xcresult.OpenAttachment(xcPath, pkg.Ptr(attachment.PayloadRef.ObjID))
		}
	}

	return nil, nil
}

func extractIdentifiableObjects(xcPath string, d xcresult.Decoder) []model.ActionTestSummary {
	var result []model.ActionTestSummary

	switch v := d.(type) {
	case *model.ActionTestSummaryGroup:
		for _, test := range v.Subtests {
			objects := extractIdentifiableObjects(xcPath, test)
			result = append(result, objects...)
		}
	case *model.ActionTestMetadata:
		if v.SummaryRef != nil {
			json, err := xcresult.ToJson(xcPath, pkg.Ptr(v.SummaryRef.ObjID))
			if err != nil {
				panic(err)
			}

			testSummary := xcresult.DecodeObject[model.ActionTestSummary](json)
			result = append(result, testSummary)
		}
	default:
		panic(fmt.Sprintf("unexpected identifiable object (type=%T)", v))
	}

	return result
}

func extractTests(xcPath string, record model.ActionsInvocationRecord) ([]model.ActionTestSummary, error) {
	var result []model.ActionTestSummary

	for _, action := range record.Actions {
		testsRef := action.Result.TestsRef
		if testsRef == nil {
			continue
		}

		fileID := pkg.Ptr(testsRef.ObjID)
		json, err := xcresult.ToJson(xcPath, fileID)
		if err != nil {
			return nil, err
		}

		testPlanRunSummaries := xcresult.DecodeObject[model.ActionTestPlanRunSummaries](json)
		for _, testPlanRunSummary := range testPlanRunSummaries.Summaries {
			for _, testableSummary := range testPlanRunSummary.TestableSummaries {
				for _, testIdentifiableObj := range testableSummary.Tests {
					objects := extractIdentifiableObjects(xcPath, testIdentifiableObj)
					result = append(result, objects...)
				}
			}
		}
	}

	return result, nil
}
