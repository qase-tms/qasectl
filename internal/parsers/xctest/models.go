package xctest

type TestMeta struct {
	Device        string
	Configuration string
	StartTime     *float64
	EndTime       *float64
	Suite         string
}

type Structure struct {
	Actions struct {
		Values []struct {
			ActionResult struct {
				TestsRef *struct {
					ID struct {
						Value string `json:"_value"`
					} `json:"id"`
				} `json:"testsRef,omitempty"`
			} `json:"actionResult"`
			RunDestination *RunDestination `json:"runDestination"`
			EndedTime      *struct {
				Value string `json:"_value"`
			} `json:"endedTime"`
			StartedTime *struct {
				Value string `json:"_value"`
			} `json:"startedTime"`
		} `json:"_values"`
	} `json:"actions"`
}

type RunDestination struct {
	DisplayName struct {
		Value string `json:"_value"`
	} `json:"displayName"`
}

type ActionTestPlanRunSummaries struct {
	Summaries struct {
		Values []struct {
			Name struct {
				Value string `json:"_value"`
			} `json:"name"`
			TestableSummaries struct {
				Values []struct {
					Tests struct {
						Values []struct {
							Subtests *Subtests `json:"subtests"`
						} `json:"_values"`
					} `json:"tests"`
				} `json:"_values"`
			} `json:"testableSummaries"`
		} `json:"_values"`
	} `json:"summaries"`
}

type Subtests struct {
	Values []struct {
		Type struct {
			Name string `json:"_name"`
		}
		Duration struct {
			Value string `json:"_value"`
		} `json:"duration"`
		IdentifierURL struct {
			Value string `json:"_value"`
		} `json:"identifierURL"`
		Name struct {
			Value string `json:"_value"`
		} `json:"name"`
		SummaryRef *struct {
			ID struct {
				Value string `json:"_value"`
			} `json:"id"`
		} `json:"summaryRef"`
		TestStatus *struct {
			Value string `json:"_value"`
		} `json:"testStatus"`
		Subtests *Subtests `json:"subtests"`
	} `json:"_values"`
}

type XCTest struct {
	Name      string
	Suites    []string
	Action    ActionTestSummary
	Metadata  TestMeta
	Signature string
	Duration  float64
}

type ActionTestSummary struct {
	ActivitySummaries ActivitySummaries `json:"activitySummaries"`
	TestStatus        struct {
		Value string `json:"_value"`
	} `json:"testStatus"`
	FailureSummaries *struct {
		Values []FailureSummary `json:"_values"`
	} `json:"failureSummaries,omitempty"`
}

type FailureSummary struct {
	Attachments *Attachments `json:"attachments,omitempty"`
	Message     struct {
		Value string `json:"_value"`
	} `json:"message"`
	LineNumber struct {
		Value string `json:"_value"`
	} `json:"lineNumber"`
	UUID struct {
		Value string `json:"_value"`
	} `json:"uuid"`
}

type ActivitySummaries struct {
	Values []struct {
		ActivityType struct {
			Value string `json:"_value"`
		} `json:"activityType"`
		Finish struct {
			Value string `json:"_value"`
		} `json:"finish"`
		Start struct {
			Value string `json:"_value"`
		} `json:"start"`
		Title struct {
			Value string `json:"_value"`
		} `json:"title"`
		Attachments       *Attachments       `json:"attachments,omitempty"`
		Subactivities     *ActivitySummaries `json:"subactivities,omitempty"`
		FailureSummaryIDs *struct {
			Values []struct {
				Value string `json:"_value"`
			} `json:"_values"`
		} `json:"failureSummaryIDs"`
	} `json:"_values"`
}

type Attachments struct {
	Values []Attachment `json:"_values"`
}

type Attachment struct {
	Filename struct {
		Value string `json:"_value"`
	} `json:"filename"`
	Name struct {
		Value string `json:"_value"`
	} `json:"name"`
	PayloadRef struct {
		ID struct {
			Value string `json:"_value"`
		} `json:"id"`
	} `json:"payloadRef"`
}

type QaseId struct {
	ID int64 `json:"caseId"`
}
