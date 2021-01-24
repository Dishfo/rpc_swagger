package gen

import "github.com/go-openapi/loads"

type targetDescription struct {
	UseModel bool `json:"use_model"`
}

type GenOpt struct {
	ModelPath  string `json:"model_path"`
	ModelPkg   string `json:"model_pkg"`
	TargetPath string `json:"target_path"`
	Type       string `json:"type"` //server/client
}

type codeBuilder struct {
	GenOpt
	targetDescription

	doc  *loads.Document
	file string `json:"file"`
}

func Generate(filePath string, opt *GenOpt) (err error) {
	cb := codeBuilder{
		GenOpt: *opt,
		file:   filePath,
	}

	err = cb.analysisSpec()
	if err != nil {
		return err
	}

	return nil
}

func (m *codeBuilder) analysisSpec() (err error) {
	return nil
}
