package result

type Attachment struct {
	ID          *int64
	Name        string
	FilePath    *string
	ContentType string
	Content     *[]byte
}
