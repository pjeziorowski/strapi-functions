package create_project

type CreateProjectInput struct {
	Name string
}

type CreateProjectOutput struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Mutation struct {
	CreateProject *CreateProjectOutput
}

type CreateProjectArgs struct {
	Input CreateProjectInput
}
