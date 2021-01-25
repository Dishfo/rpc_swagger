package rpc_spec

import "time"

type TestService struct {
	Field int `json:"field" gorm:"xxx"`
}

type TestT struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"-"`
	phone    string
	Value    map[string]interface{} `json:"value"`
	ArrayVal []int                  `json:"array_val"`
	Time     time.Time              `json:"time"`
}

func (s *TestService) GetByID(companyId, id string) (TestT, error) {

	return TestT{}, nil
}

func (s *TestService) GetByIDs(companyId, id string) {

}
