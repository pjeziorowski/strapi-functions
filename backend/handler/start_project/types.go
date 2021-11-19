package start_project

type StartProjectInput struct {
	Id int32
}

type StartProjectOutput struct {
	Ok bool `json:"ok"`
}

type Mutation struct {
	StartProject *StartProjectOutput
}

type StartProjectArgs struct {
	Input StartProjectInput
}
