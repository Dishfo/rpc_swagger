package rpc_spec

import (
	"errors"
	"log"
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
	ret.SwaggerType = "boolean"
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

	if typ.String() == "time.Time" {
		ret.Format = "date-time"
		ret.HasFormat = true
		ret.SwaggerType = "string"
		ret.IsPrimitive = true
		return ret, nil
	}

	var isStringer = isStringer(typ)

	if isStringer && ctx.InOrOut == serializeOut {
		ret.SwaggerType = "string"
		ret.IsPrimitive = true
		return ret, nil
	}

	isMarshal := isMarshal(typ)
	isUnMarshal := isUnMarshal(typ)

	ret.CustomUnMarshal = isUnMarshal
	ret.CustomMarshal = isMarshal

	ret.IsReference = true
	if typ.Name() == "" {
		log.Println("Anonymous")
		ret.SwaggerType = ctx.NamePrefix
	} else {
		ret.SwaggerType = TypeName(typ)
	}

	if ctx.NeedRegisterModel && !ctx.CheckModelExist(ret.SwaggerType) {
		//parse model
		var definitionSpec DefinitionSpec
		definitionSpec.Pkg = typ.PkgPath()
		definitionSpec.GoTypeName = typ.Name()
		if typ.Name() == "" {
			log.Println("Anonymous")
			definitionSpec.ModelName = ctx.NamePrefix
		} else {
			definitionSpec.ModelName = TypeName(typ)
		}

		fieldNum := typ.NumField()
		for i := 0; i < fieldNum; i++ {
			fieldType := typ.Field(i)
			if NeedOmit(fieldType) {
				continue
			}
			newCtx := ctx
			newCtx.Anonymous = fieldType.Anonymous
			if !newCtx.Anonymous {
				newCtx.NamePrefix = definitionSpec.ModelName + "." + fieldType.Name
			}

			fieldTypeSpec, err := ctx.AnalysisProxy(newCtx, fieldType.Type)
			if err != nil {
				return TypeSpec{}, err
			}

			if !fieldType.Anonymous {
				fieldSpec := FieldSpec{
					TypeSpec: fieldTypeSpec,
					Name:     fieldType.Name,
				}
				definitionSpec.Properties = append(definitionSpec.Properties, fieldSpec)
			} else {
				definitionSpec.Properties = append(definitionSpec.Properties, getAnonymousFields(fieldTypeSpec)...)
			}
		}

		ret.ReferenceType = &definitionSpec
	} else {
		ret.ReferenceType = ctx.GetModel(ret.SwaggerType)
	}

	if ret.IsReference {
		if ret.ReferenceType == nil || len(ret.ReferenceType.Properties) == 0 {
			ret.IsReference = false
			ret.IsPrimitive = true
			ret.SwaggerType = "object"
		}
	}

	return ret, nil
}

func getAnonymousFields(fieldSpec TypeSpec) []FieldSpec {
	if fieldSpec.IsReference && fieldSpec.ReferenceType != nil {
		return fieldSpec.ReferenceType.Properties
	}
	return []FieldSpec{}
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

type PtrConvert struct {
}

func (p *PtrConvert) GetSwaggerType(ctx AnalysisContext, typ reflect.Type) (TypeSpec, error) {
	var ret TypeSpec
	var isStringer = isStringer(typ)

	if isStringer && ctx.InOrOut == serializeOut {
		ret.SwaggerType = "string"
		ret.IsPrimitive = true
		return ret, nil
	}

	realType := indirectType(typ)
	ret, err := ctx.AnalysisProxy(ctx, realType)
	log.Println("is a pointer ")
	if err != nil {
		log.Printf("parse ptr type %s failed :%s", typ.String(), err.Error())
		return TypeSpec{}, err
	}
	ret.Pointer = true
	return ret, nil
}
