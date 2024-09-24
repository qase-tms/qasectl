package run

type Environment struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
	Slug  string `json:"slug"`
}

type Milestone struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

type Plan struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

type Run struct {
	ID int64 `json:"id"`
}
