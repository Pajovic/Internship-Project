package models

type Shop struct {
	ID		string 	`json:"id"`
	Name	string 	`json:"name"`
	IDC     string 	`json:"idc"`
	Lat  	float64 `json:"lat"`
	Lon   	float64 `json:"lon"`
}
