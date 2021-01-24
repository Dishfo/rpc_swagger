package rpc_spec

type SpecBuilder struct {
	Definitions map[string]DefinitionSpec
	Paths       []PathSpec
}

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		Definitions: map[string]DefinitionSpec{},
	}
}

func (b *SpecBuilder) ChildPathBuilder() *PathSpecBuilder {
	inst := &PathSpecBuilder{
		center: b,
	}
	return inst
}

func (b *SpecBuilder) Build() Spec {
	var definitions []DefinitionSpec
	for _, spec := range b.Definitions {
		definitions = append(definitions, spec)
	}

	return Spec{
		Paths:       b.Paths,
		Definitions: definitions,
	}
}

func (b *SpecBuilder) AppendDefinitions(spec DefinitionSpec) {
	b.Definitions[spec.ModelName] = spec
}

func (b *SpecBuilder) ExistDefinition(name string) bool {

	return false
}

func (b *SpecBuilder) AppendPath(spec PathSpec) {
	b.Paths = append(b.Paths, spec)
}

func (b *SpecBuilder) GlobalSeq() int {
	return 0
}
