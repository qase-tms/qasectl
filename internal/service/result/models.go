package result

type UploadParams struct {
	RunID       int64
	Title       string
	Description string
	Batch       int64
	Project     string
	Suite       string
}
