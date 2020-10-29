package models

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int32   `json:"quantity"`
	IDC      string  `json:"idc"`
}
