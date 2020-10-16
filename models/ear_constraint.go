package models

type EarConstraint struct {
	Idear         string `json:"idear"`
	Idrc          string `json:"idrc"`
	Idsc          string `json:"idsc"`
	Property      string `json:"property"`
	Operator      string `json:"operator"`
	PropertyValue int    `json:"property_value"`
}
