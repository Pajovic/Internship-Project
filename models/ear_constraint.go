package models

type EarConstraint struct {
	IDEAR         string `json:"idear"`
	IDRC          string `json:"idrc"`
	IDSC          string `json:"idsc"`
	Property      string `json:"property"`
	Operator      string `json:"operator"`
	PropertyValue int    `json:"property_value"`
}
