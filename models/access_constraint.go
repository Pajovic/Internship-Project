package models

type AccessConstraint struct {
	ID            string  `json:"id"`
	IDEAR         string  `json:"idear"`
	OperatorID    int32   `json:"operatorId"`
	PropertyID    int64   `json:"propertyId"`
	PropertyValue float64 `json:"propertyValue"`
}
