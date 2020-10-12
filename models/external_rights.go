package models

//ExternalRights .
type ExternalRights struct {
	ID       string `json:"id"`
	Read     bool   `json:"r"`
	Update   bool   `json:"u"`
	Delete   bool   `json:"d"`
	Approved bool   `json:"approved"`
	IDSC     string `json:"idsc"`
	IDRC     string `json:"idrc"`
}
