package rpc_spec

import "reflect"

type SpecBuilder struct {
	Definitions        map[string]DefinitionSpec
	Paths              []PathSpec
	specificTypeParser func(typ reflect.Type) (TypeSpec, bool, error)
}

func NewSpecBuilder(SpecificTypeParser func(typ reflect.Type) (TypeSpec, bool, error)) *SpecBuilder {
	return &SpecBuilder{
		Definitions:        map[string]DefinitionSpec{},
		specificTypeParser: SpecificTypeParser,
	}
}

func (b *SpecBuilder) ChildPathBuilder() *PathSpecBuilder {
	inst := &PathSpecBuilder{
		center:             b,
		specificTypeParser: b.specificTypeParser,
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
	_, exist := b.Definitions[name]
	return exist
}

func (b *SpecBuilder) GetDefinition(name string) *DefinitionSpec {
	def, exist := b.Definitions[name]
	if !exist {
		return nil
	}
	return &def
}

func (b *SpecBuilder) AppendPath(spec PathSpec) {
	b.Paths = append(b.Paths, spec)
}

func (b *SpecBuilder) GlobalSeq() int {
	return 0
}
