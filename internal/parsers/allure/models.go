package allure

type Test struct {
	UUID            string        `json:"uuid"`
	HistoryID       string        `json:"historyId"`
	Status          string        `json:"status"`
	StatusDetails   StatusDetails `json:"statusDetails"`
	Stage           string        `json:"stage"`
	Steps           []TestStep    `json:"steps"`
	Attachments     []Attachment  `json:"attachments"`
	Links           []Link        `json:"links"`
	Start           float64       `json:"start"`
	Name            string        `json:"name"`
	FullName        string        `json:"fullName"`
	Labels          []Label       `json:"labels"`
	Stop            float64       `json:"stop"`
	Description     *string       `json:"description"`
	DescriptionHTML *string       `json:"descriptionHtml"`
	Params          []Parameter   `json:"parameters,omitempty"`
}

type TestStep struct {
	Name          string  `json:"name"`
	Start         float64 `json:"start"`
	Stop          float64 `json:"stop"`
	Stage         string  `json:"stage"`
	Status        string  `json:"status"`
	StatusDetails struct {
	} `json:"statusDetails"`
	Attachments []Attachment `json:"attachments"`
	Parameters  []any        `json:"parameters"`
	Steps       []TestStep   `json:"steps"`
	Description string       `json:"description"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Link struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	URL  *string `json:"url"`
}

type Attachment struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type StatusDetails struct {
	Message *string `json:"message"`
	Trace   *string `json:"trace"`
}
