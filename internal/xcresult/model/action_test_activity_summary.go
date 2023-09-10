package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
	"slices"
)

type ActivityType string

const (
	ActivityUserCreated         ActivityType = "com.apple.dt.xctest.activity-type.userCreated"
	ActivityFailure             ActivityType = "com.apple.dt.xctest.activity-type.testAssertionFailure"
	ActivityAttachmentContainer ActivityType = "com.apple.dt.xctest.activity-type.attachmentContainer"
	ActivityNone                ActivityType = "none"
)

var activityTypes = []ActivityType{
	ActivityUserCreated,
	ActivityFailure,
	ActivityAttachmentContainer,
	ActivityNone,
}

type ActionTestActivitySummary struct {
	Title         string
	ActivityType  ActivityType
	Attachments   []ActionTestAttachment
	Subactivities []ActionTestActivitySummary
}

func (a *ActionTestActivitySummary) TypeName() string {
	return "ActionTestActivitySummary"
}

func (a *ActionTestActivitySummary) Decode(m map[string]any) {
	a.Title = xcresult.DecodeString(m["title"].(map[string]any))

	activityType := xcresult.DecodeString(m["activityType"].(map[string]any))
	a.ActivityType = a.activityTypeFromString(activityType)

	if attachments, ok := m["attachments"].(map[string]any); ok {
		a.Attachments = xcresult.DecodeArray[ActionTestAttachment](attachments)
	}

	if subactivities, ok := m["subactivities"].(map[string]any); ok {
		a.Subactivities = xcresult.DecodeArray[ActionTestActivitySummary](subactivities)
	}
}

func (a *ActionTestActivitySummary) activityTypeFromString(s string) ActivityType {
	if slices.Contains(activityTypes, ActivityType(s)) {
		return ActivityType(s)
	}

	return ActivityNone
}
