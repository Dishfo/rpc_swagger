package rpc_spec

import (
	"errors"
	"reflect"
)

type StringConvert struct {
}

func (s *StringConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	ret.IsPrimitive = true
	ret.SwaggerType = "string"
	return ret, nil
}

type NumberConvert struct {
}

func (n *NumberConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	ret.IsPrimitive = true

	kind := typ.Kind()
	if kind > reflect.Uint64 {
		ret.SwaggerType = "number"
		if kind == reflect.Float32 {
			ret.Format = "float"
			ret.HasFormat = true
		}
	} else {
		ret.SwaggerType = "integer"
		if kind == reflect.Int {
			ret.Format = "int32"
			ret.HasFormat = true
		}
	}

	return ret, nil
}

type BoolConvert struct {
}

func (b *BoolConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	ret.IsPrimitive = true
	ret.SwaggerType = "bool"
	return ret, nil
}

type InterfaceConvert struct {
}

func (i *InterfaceConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	ret.IsPrimitive = true
	ret.SwaggerType = "object"
	return ret, nil
}

//next type_convert may be more complex
type StructConvert struct {
}

//TODO need handle Anonymous struct
func (s *StructConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	if typ.Name() == "time.Time" {
		ret.Format = "date-time"
		ret.HasFormat = true
		ret.SwaggerType = "string"
		return ret, nil
	}

	ret.IsReference = true
	ret.SwaggerType = TypeName(typ)

	if ctx.NeedRegisterModel && !ctx.CheckModelExist(ret.SwaggerType) {
		//parse model
		var definitionSpec DefinitionSpec
		definitionSpec.Pkg = typ.PkgPath()
		definitionSpec.GoTypeName = typ.Name()
		definitionSpec.ModelName = TypeName(typ)

		fieldNum := typ.NumField()
		for i := 0; i < fieldNum; i++ {
			fieldType := typ.Field(i)
			if NeedOmit(fieldType) {
				continue
			}
			newCtx := ctx
			newCtx.Anonymous = fieldType.Anonymous

			fieldTypeSpec, err := ctx.AnalysisProxy(newCtx, fieldType.Type)
			if err != nil {
				return TypeSpec{}, err
			}

			fieldSpec := FieldSpec{
				TypeSpec: fieldTypeSpec,
				Name:     fieldType.Name,
			}

			definitionSpec.Properties = append(definitionSpec.Properties, fieldSpec)
		}

		ret.ReferenceType = &definitionSpec
	}

	return ret, nil
}

type MapConvert struct {
}

func (m *MapConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	keyTyp := typ.Key()
	if keyTyp.Kind() != reflect.String {
		return TypeSpec{}, errors.New("map key must be string type")
	}

	elemTyp := typ.Elem()
	ret.IsPrimitive = false
	ret.IsMap = true
	ret.SwaggerType = "object"
	newCtx := ctx
	newCtx.PrefixSpace = ctx.PrefixSpace + "  "
	itemTypeSpec, err := ctx.AnalysisProxy(newCtx, elemTyp)
	if err != nil {
		return TypeSpec{}, err
	}

	ret.ItemSpec = &itemTypeSpec

	return ret, nil
}

type ArrayConvert struct {
}

func (a *ArrayConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	eleType := typ.Elem()
	if eleType.Kind() == reflect.Uint8 {
		ret.SwaggerType = "string"
		ret.Format = "byte"
		ret.HasFormat = true
		ret.IsPrimitive = true
		return ret, nil
	}

	elemTyp := typ.Elem()
	ret.IsPrimitive = true
	ret.SwaggerType = "array"
	newCtx := ctx
	newCtx.PrefixSpace = ctx.PrefixSpace + "  "

	itemTypeSpec, err := ctx.AnalysisProxy(newCtx, elemTyp)
	if err != nil {
		return TypeSpec{}, err
	}

	ret.ItemSpec = &itemTypeSpec

	return ret, nil
}
