package result

type Data struct {
	Action         string  `json:"action"`
	ExceptedResult *string `json:"expected_result"`
	InputData      *string
	Attachments    []Attachment
}
