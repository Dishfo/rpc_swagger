package rpc_spec

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type PathSpecBuilder struct {
	center             *SpecBuilder
	specificTypeParser func(typ reflect.Type) (TypeSpec, bool, error)

	serviceName string
	methodName  string
	paramsTypes []reflect.Type
	paramsNames []string
	resultType  reflect.Type

	paramFields []FieldSpec
	resultSpec  TypeSpec
	description string

	errs []error

	result PathSpec
}

func (b *PathSpecBuilder) Build() PathSpec {
	b.result.Method = b.methodName
	b.result.ServiceName = b.serviceName
	b.result.RpcServiceMethod = fmt.Sprintf("%s_%s", b.serviceName, parselizeMethodName(b.methodName))
	b.result.RpcPath = fmt.Sprintf("%s/%s", b.serviceName, b.methodName)

	b.result.ParamSpec.Properties = b.paramFields

	if len(b.paramFields) > 0 {
		b.result.ParamSpec.ModelName = fmt.Sprintf("%s%sParamList", b.serviceName, b.methodName)
		b.center.AppendDefinitions(b.result.ParamSpec)
	}
	if b.resultType != nil {
		b.result.ResultList = &b.resultSpec
		b.result.HasRes = true
	}
	b.result.Description = b.description
	b.center.AppendPath(b.result)

	return b.result
}

func (b *PathSpecBuilder) IsValid() bool {
	if len(b.serviceName) == 0 {
		return false
	}

	if len(b.methodName) == 0 {
		return false
	}

	if b.methodName[0] < 'A' || b.methodName[0] > 'Z' {
		return false
	}
	log.Println(b.errs)
	return len(b.errs) == 0
}

func (b *PathSpecBuilder) SetServiceName(service string) *PathSpecBuilder {
	b.serviceName = service
	return b
}

func (b *PathSpecBuilder) SetMethod(method string) *PathSpecBuilder {
	b.methodName = method
	return b
}

func (b *PathSpecBuilder) AppendParam(name string, typ reflect.Type) *PathSpecBuilder {
	if typ.String() == "context.Context" && len(b.paramsNames) == 0 {
		return b
	}

	b.paramsNames = append(b.paramsNames, name)
	b.paramsTypes = append(b.paramsTypes, typ)

	//parse current params,get type of it ,if
	// it's struct should create definition spec and register to center

	b.analysisParamType(name, typ)
	return b
}

func (b *PathSpecBuilder) SetDescription(comments []string) *PathSpecBuilder {
	b.description = PackingDescription(comments)

	return b
}

func (b *PathSpecBuilder) SetResult(name string, typ reflect.Type) *PathSpecBuilder {
	b.resultType = typ
	log.Println("value nil ", b.resultType)
	if typ == nil {
		return b
	}
	ctx := AnalysisContext{
		Ctx:               context.TODO(),
		IsPointer:         false,
		InOrOut:           serializeOut,
		NeedRegisterModel: true,
		CheckModelExist:   b.checkDefinitionExist,
		GetModel:          b.getModelDefinition,
		AnalysisProxy:     b.analysisType,
		occurType:         map[string]bool{},
	}

	typeSpec, err := b.analysisType(ctx, typ)
	if err != nil {
		b.addErr(err)
		return b
	}
	b.resultSpec = typeSpec
	return b
}

func (b *PathSpecBuilder) addErr(err error) {
	if err != nil {
		b.errs = append(b.errs, err)
	}

}

const (
	off1 = 'a' - byte('A')
)

func parselizeMethodName(method string) string {
	var builder strings.Builder
	builder.WriteByte(method[0] + byte(off1))
	builder.WriteString(method[1:])
	return builder.String()
}

type serializeType int

const (
	serializeIn serializeType = iota + 1
	serializeOut
)

type AnalysisContext struct {
	Ctx        context.Context
	Anonymous  bool
	occurType  map[string]bool
	NamePrefix string
	IsPointer  bool
	InOrOut    serializeType //1 ==in 2 == out

	NeedRegisterModel bool
	CheckModelExist   func(string) bool //can't be nil if NeedRegisterModel equal true
	GetModel          func(string) *DefinitionSpec

	//for more complex inner model
	AnalysisProxy func(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error)

	PrefixSpace string
}

type TypeConvert interface {
	GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error)
}

func (b *PathSpecBuilder) checkDefinitionExist(name string) bool {
	return b.center.ExistDefinition(name)
}

func (b *PathSpecBuilder) getModelDefinition(modelName string) *DefinitionSpec {
	return b.center.GetDefinition(modelName)
}

func (b *PathSpecBuilder) analysisType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {

	var realType reflect.Type = typ
	if b.specificTypeParser != nil {
		ret, hit, err := b.specificTypeParser(typ)
		if err != nil {
			return TypeSpec{}, err
		}
		if hit {
			return ret, nil
		}
	}

	kindType := typeMerge(realType.Kind())
	conv := builtinTypeConvert[kindType]
	if conv == nil {
		return TypeSpec{}, errors.New(fmt.Sprintf("can't handle this type : %s", typ.String()))
	}

	typeSpec, err := conv.GetSwaggerType(ctx, realType)

	if err != nil {
		log.Printf("can't handle this case %s %s", err.Error(), typ.String())
		return TypeSpec{}, err
	}

	if typeSpec.IsReference && ctx.NeedRegisterModel && !ctx.Anonymous {
		if typeSpec.ReferenceType != nil && !b.center.ExistDefinition(typeSpec.SwaggerType) {
			b.center.AppendDefinitions(*typeSpec.ReferenceType)
		}
	}
	typeSpec.PrefixSpace = ctx.PrefixSpace

	return typeSpec, nil
}

func (b *PathSpecBuilder) analysisParamType(name string, typ reflect.Type) {

	//if typ is umMarshal need set this type as x-marshal

	ctx := AnalysisContext{
		Ctx:       context.TODO(),
		IsPointer: false,
		InOrOut:   serializeIn,

		NeedRegisterModel: true,
		CheckModelExist:   b.checkDefinitionExist,
		GetModel:          b.getModelDefinition,
		AnalysisProxy:     b.analysisType,
		occurType:         map[string]bool{},
	}

	typeSpec, err := b.analysisType(ctx, typ)
	if err != nil {
		b.errs = append(b.errs, err)
		return
	}
	fieldSpec := FieldSpec{
		TypeSpec: typeSpec,
		Name:     name,
	}
	b.paramFields = append(b.paramFields, fieldSpec)
}

var (
	builtinTypeConvert = map[int]TypeConvert{
		0: &NumberConvert{},    //integer
		1: &StringConvert{},    //string
		2: &ArrayConvert{},     //array
		3: &MapConvert{},       //map
		4: &InterfaceConvert{}, //interface
		5: &StructConvert{},    //struct
		6: &BoolConvert{},      //bool
		7: &PtrConvert{},       //ptr
	}
)

func typeMerge(kind reflect.Kind) int {
	switch {
	case reflect.String == kind:
		return 1
	case kind >= reflect.Int && kind <= reflect.Float64:
		return 0
	case reflect.Bool == kind:
		return 6
	case reflect.Array == kind || reflect.Slice == kind:
		return 2
	case reflect.Map == kind:
		return 3
	case reflect.Interface == kind:
		return 4
	case reflect.Struct == kind:
		return 5
	case reflect.Ptr == kind:
		return 7
	}
	return -1
}
