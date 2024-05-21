package result

import "github.com/google/uuid"

type Attachment struct {
	ID          *uuid.UUID `json:"id"`
	Name        string     `json:"file_name"`
	FilePath    *string    `json:"file_path"`
	ContentType string     `json:"mime_type"`
	Content     *[]byte    `json:"content"`
}
