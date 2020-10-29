package models

type Company struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsMain bool   `json:"isMain"`
}
