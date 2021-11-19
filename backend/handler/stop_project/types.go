package stop_project

type StopProjectInput struct {
	Id int32
}

type StopProjectOutput struct {
	Ok bool `json:"ok"`
}

type Mutation struct {
	StopProject *StopProjectOutput
}

type StopProjectArgs struct {
	Input StopProjectInput
}
