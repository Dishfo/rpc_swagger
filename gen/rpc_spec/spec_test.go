package rpc_spec

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type StrA string

func TestRpcServiceAnalysis_LoadServiceFunction(t *testing.T) {
	ayInst := NewRpcServiceAnalysis(&AnalysisOpt{})

	ayInst.loadLocalPackage("github.com/Dishfo/rpc_swagger/")
	t.Log("complete")

	paramNames, find := ayInst.getFunctionParamName(ServiceFunction{
		Package:  "github.com/Dishfo/rpc_swagger/gen/rpc_spec",
		Service:  "TestService",
		Function: "GetByID",
	})

	require.True(t, find)
	t.Log(paramNames)

	var a TestService
	rv := reflect.ValueOf(a)
	rt := rv.Type()
	fmt.Println(rt.String(), rt.Kind())

}

func TestRpcServiceAnalysis_AppointService(t *testing.T) {
	inst := NewRpcServiceAnalysis(&AnalysisOpt{
		ServerName: "EmptyServer",
		TargetFile: "/home/dishfo/go/src/github.com/Dishfo/rpc_swagger/test.yaml",
	})
	inst.AppointService(ServiceRegister{
		RegisterName: "test",
		ServiceInst:  &TestService{},
	})
	t.Log("complete")

	err := inst.Render()
	t.Log(err)
}
