package delete_project

type DeleteProjectInput struct {
	Id int32
}

type DeleteProjectOutput struct {
	Ok bool `json:"ok"`
}

type Mutation struct {
	DeleteProject *DeleteProjectOutput
}

type DeleteProjectArgs struct {
	Input DeleteProjectInput
}
