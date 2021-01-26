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
	ST       struct {
		Inner string
	} `json:"st"`
	TestT2
	T3 *TestT3
}

type TestT3 struct {
}

func (t *TestT3) String() string {
	panic("implement me")
}

type TestT2 struct {
	A int
	B int
	T *TestT3
}

//func (t TestT) String() string {
//	panic("implement me")
//}

func (s *TestService) GetByID(companyId, id *string) (*TestT, error) {

	return nil, nil
}

func (s *TestService) GetByIDs(companyId, id string) {

}
