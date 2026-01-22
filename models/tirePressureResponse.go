package vehicle

type TirePressureResponse struct {
	Source   string      `json:"source"`
	FrontPsi *float64    `json:"frontPsi,omitempty"`
	RearPsi  *float64    `json:"rearPsi,omitempty"`
	Unit     string      `json:"unit,omitempty"`
	Note     string      `json:"note,omitempty"`
	Raw      interface{} `json:"raw,omitempty"`
}
