package xctest

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
			RunDestination struct {
				DisplayName struct {
					Value string `json:"_value"`
				} `json:"displayName"`
			} `json:"runDestination"`
		} `json:"_values"`
	} `json:"actions"`
}

type ActionTestPlanRunSummaries struct {
	Summaries struct {
		Values []struct {
			TestableSummaries struct {
				Values []struct {
					Tests struct {
						Values []struct {
							Subtests struct {
								Values []struct {
									Name struct {
										Value string `json:"_value"`
									} `json:"name"`
									Subtests struct {
										Values []struct {
											Name struct {
												Value string `json:"_value"`
											} `json:"name"`
											Subtests struct {
												Values []struct {
													Duration struct {
														Value string `json:"_value"`
													} `json:"duration"`
													IdentifierURL struct {
														Value string `json:"_value"`
													} `json:"identifierURL"`
													Name struct {
														Value string `json:"_value"`
													} `json:"name"`
													SummaryRef struct {
														ID struct {
															Value string `json:"_value"`
														} `json:"id"`
													} `json:"summaryRef"`
													TestStatus struct {
														Value string `json:"_value"`
													} `json:"testStatus"`
												} `json:"_values"`
											} `json:"subtests"`
										} `json:"_values"`
									} `json:"subtests"`
								} `json:"_values"`
							} `json:"subtests"`
						} `json:"_values"`
					} `json:"tests"`
				} `json:"_values"`
			} `json:"testableSummaries"`
		} `json:"_values"`
	} `json:"summaries"`
}

type ActionTestSummary struct {
	ActivitySummaries ActivitySummaries `json:"activitySummaries"`
	TestStatus        struct {
		Value string `json:"_value"`
	} `json:"testStatus"`
	FailureSummaries *struct {
		Values []struct {
			Attachments *Attachments `json:"attachments,omitempty"`
			Message     struct {
				Value string `json:"_value"`
			} `json:"message"`
		} `json:"_values"`
	} `json:"failureSummaries,omitempty"`
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
		Attachments   *Attachments   `json:"attachments,omitempty"`
		Subactivities *Subactivities `json:"subactivities,omitempty"`
	} `json:"_values"`
}

type Subactivities struct {
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
		Attachments   *Attachments   `json:"attachments,omitempty"`
		Subactivities *Subactivities `json:"subactivities"`
		Title         struct {
			Value string `json:"_value"`
		} `json:"title"`
	} `json:"_values"`
}

type Attachments struct {
	Values []struct {
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
	} `json:"_values"`
}

type QaseId struct {
	ID int64 `json:"caseId"`
}
