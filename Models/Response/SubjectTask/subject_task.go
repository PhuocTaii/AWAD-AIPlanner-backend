package subjecttask

type SubjectDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SubjectTask struct {
	Subject SubjectDto `json:"subject"`
	Tasks   int        `json:"tasks"`
}
