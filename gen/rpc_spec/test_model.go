package rpc_spec

type TestService struct {
	Field int `json:"field" gorm:"xxx"`
}

type TestT struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"-"`
	phone string
}

func (s *TestService) GetByID(companyId, id string) (TestT, error) {

	return TestT{}, nil
}
