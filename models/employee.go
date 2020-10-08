package models

// Employee basic model
type Employee struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CompanyID string `json:"companyId"`
	C         bool   `json:"c"`
	R         bool   `json:"r"`
	U         bool   `json:"u"`
	D         bool   `json:"d"`
}
