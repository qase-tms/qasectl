package result

type Data struct {
	Action         string
	ExceptedResult *string
	InputData      *string
	Attachments    []Attachment
}
