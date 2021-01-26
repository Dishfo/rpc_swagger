package rpc_spec

import (
	"bytes"
	"github.com/Dishfo/rpc_swagger/templates"
	"github.com/Dishfo/wire_dao/gen_dao"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type AnalysisOpt struct {
	TargetFile string
	ServerName string
	Email      string
	Contact    string

	SpecificTypeParser func(typ reflect.Type) (TypeSpec, bool, error)
}

type RpcServiceAnalysis struct {
	opt *AnalysisOpt

	//server Spec
	Spec Spec

	//parsed package
	scanners map[string]*gen_dao.Scanner
}

type ServiceRegister struct {
	RegisterName string
	ServiceInst  interface{}
}

func NewRpcServiceAnalysis(opt *AnalysisOpt) *RpcServiceAnalysis {
	return &RpcServiceAnalysis{
		opt:      opt,
		scanners: map[string]*gen_dao.Scanner{},
	}
}

func (a *RpcServiceAnalysis) parseParams(params []reflect.Type) (ret DefinitionSpec, err error) {
	for _, param := range params {
		log.Printf("params name %s", param.String())
	}

	return
}

func (a *RpcServiceAnalysis) getFunctionParamName(sf ServiceFunction) ([]string, bool) {

	var scan *gen_dao.Scanner
	//find scanner
	for pack, scanner := range a.scanners {
		if pack == sf.Package {
			scan = scanner
			break
		}

		if strings.HasPrefix(sf.Package, pack) {
			scan = scanner
			break
		}
	}

	if scan == nil {
		log.Printf("can't %s pack", sf.Package)
		return nil, false
	}

	funcDesc, err := scan.FindFunction(gen_dao.FunctionID{
		Package:  sf.Package,
		Receiver: sf.Service,
		Function: sf.Function,
	})
	if err != nil {
		log.Printf("can't find function %s because of %s", sf.String(), err.Error())
		return nil, false
	}

	var paramNames []string
	for _, param := range funcDesc.Params {
		paramNames = append(paramNames, param.Name)
	}

	return paramNames, true
}

func (a *RpcServiceAnalysis) loadLocalPackage(packName string) (err error) {

	for pack := range a.scanners {
		if pack == packName {
			return
		}

		if strings.HasPrefix(packName, pack) {
			return
		}
	}

	goPath := os.Getenv("GOPATH")

	scanner, err := gen_dao.ScanFiles(&gen_dao.ScanConfig{
		TopDir: filepath.Join(goPath, "src", packName),
	})

	if err != nil {
		log.Printf("create scanner failed :%s", err.Error())
		return err
	}

	a.scanners[packName] = scanner

	return nil
}

func (a *RpcServiceAnalysis) AppointService(services ...ServiceRegister) (err error) {
	//existTypeSpes := make(map[string]DefinitionSpec)
	specBuild := NewSpecBuilder(a.opt.SpecificTypeParser)

	for _, service := range services {
		//parse code

		rv := reflect.ValueOf(service.ServiceInst)

		rt := rv.Type()
		indirectRt := reflect.Indirect(rv).Type()

		log.Println("package ", indirectRt.PkgPath())
		err = a.loadLocalPackage(indirectRt.PkgPath())
		if err != nil {
			log.Printf("parse package failed :%s", err.Error())
			return err
		}

		methodsNum := rt.NumMethod()
		log.Println(rt.Kind(), rt.Name())
		for i := 0; i < methodsNum; i++ {
			var builder = specBuild.ChildPathBuilder()
			builder.SetServiceName(service.RegisterName)
			methodType := rt.Method(i)
			builder.SetMethod(methodType.Name)

			paramNum := methodType.Type.NumIn()
			resultNum := methodType.Type.NumOut()
			log.Println("param number ", paramNum, "result number ", resultNum)
			var paramTypes []reflect.Type
			var paramNames []string
			var resultTypes []reflect.Type
			for k := 1; k < paramNum; k++ {
				paramType := methodType.Type.In(k)
				paramTypes = append(paramTypes, paramType)
			}

			//get params name
			paramNames, _ = a.getFunctionParamName(ServiceFunction{
				Service:  indirectRt.Name(),
				Function: methodType.Name,
				Package:  indirectRt.PkgPath(),
			})

			for k, name := range paramNames {
				builder.AppendParam(name, paramTypes[k])
			}

			log.Println("params names ", paramNames)

			for k := 0; k < resultNum; k++ {
				resultType := methodType.Type.Out(k)
				resultTypes = append(resultTypes, resultType)
				log.Println(k, resultType.Name())
			}

			resultType, valid := getRealResultType(resultTypes)
			if !valid {
				continue
			}
			builder.SetResult("", resultType)

			if builder.IsValid() {
				builder.Build()
			}
		}
	}

	a.Spec = specBuild.Build()

	return nil
}

func (a *RpcServiceAnalysis) Render() (err error) {
	a.Spec.ServerName = a.opt.ServerName
	a.Spec.Email = a.opt.Email
	a.Spec.Contact = a.opt.Contact
	tmplInst := template.New("set")
	//tt.Parse()
	//tmplInst, err := templates.ParseFiles("../../templates/spec/spec.tmpl",
	//	"../../templates/spec/definition.tmpl",
	//	"../../templates/spec/path.tmpl")

	tmplInst, err = tmplInst.Parse(string(templates.MustAsset("spec/spec.tmpl")))
	if err != nil {
		log.Printf("parse templates file failed :%s", err.Error())
		return err
	}
	tmplInst, err = tmplInst.Parse(string(templates.MustAsset("spec/definition.tmpl")))
	if err != nil {
		log.Printf("parse templates file failed :%s", err.Error())
		return err
	}
	tmplInst, err = tmplInst.Parse(string(templates.MustAsset("spec/path.tmpl")))
	if err != nil {
		log.Printf("parse templates file failed :%s", err.Error())
		return err
	}

	var buffer bytes.Buffer

	err = tmplInst.Execute(&buffer, a.Spec)
	if err != nil {
		log.Printf("render templates  failed :%s", err.Error())
		return err
	}

	log.Printf("file content \n:%s", buffer.String())

	err = ioutil.WriteFile(a.opt.TargetFile, buffer.Bytes(), os.ModePerm)

	return err
}

func getRealResultType(resultTypes []reflect.Type) (reflect.Type, bool) {
	if len(resultTypes) == 0 {
		return nil, true
	}

	switch len(resultTypes) {
	case 1:
		if resultTypes[0].Name() == "error" {
			return nil, true
		}
		return indirectType(resultTypes[0]), true
	case 2:
		log.Println(resultTypes[1].Name())
		if resultTypes[1].Name() != "error" {
			return nil, false
		}
		return indirectType(resultTypes[0]), true
	default:
		return nil, false
	}

}
