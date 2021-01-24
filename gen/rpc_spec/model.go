package rpc_spec

import "fmt"

type PathSpec struct {
	ServiceName      string
	Method           string
	RpcPath          string
	RpcServiceMethod string
	ParamSpec        DefinitionSpec
	ResultList       TypeSpec
}

type DefinitionSpec struct {
	Pkg        string
	GoTypeName string

	ModelName  string
	Properties []FieldSpec
}

type FieldSpec struct {
	TypeSpec
	Name string
}

type TypeSpec struct {
	IsPrimitive bool
	IsMap       bool
	IsReference bool
	Pointer     bool

	SwaggerType string
	HasFormat   bool
	Format      string
	ItemSpec    *TypeSpec

	ReferenceType *DefinitionSpec
}

type Spec struct {
	Paths       []PathSpec
	Definitions []DefinitionSpec
	ServerName  string
}

type ServiceFunction struct {
	Package  string
	Service  string
	Function string
}

func (m *ServiceFunction) String() string {
	return fmt.Sprintf("%s.%s.%s", m.Package, m.Service, m.Function)
}
