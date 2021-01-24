package gen

type FunctionParams struct {
	Name      string    `json:"name"`
	GenSchema GenSchema `json:"gen_schema"`
}

type FunctionResult struct {
	GenSchema GenSchema `json:"gen_schema"`
}

type FunctionGenerator struct {
	ServiceName string           `json:"service_name"`
	MethodName  string           `json:"method_name"`
	GenOpt      GenOpt           `json:"gen_opt"`
	Params      []FunctionParams `json:"params"`
	Result      *FunctionResult  `json:"result"`
}

func (m *FunctionGenerator) Generate() (err error) {

	return nil
}
