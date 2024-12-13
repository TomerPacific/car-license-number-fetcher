package vehicle

type Vehicle struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  struct {
		IncludeTotal             bool   `json:"include_total"`
		Limit                    int    `json:"limit"`
		Q                        string `json:"q"`
		RecordsFormat            string `json:"records_format"`
		ResourceID               string `json:"resource_id"`
		TotalEstimationThreshold any    `json:"total_estimation_threshold"`
		Records                  []struct {
			ID                    int     `json:"_id"`
			LicenseNumber         int     `json:"mispar_rechev"`
			ProductionCountry     int     `json:"tozeret_cd"`
			ModelType             string  `json:"sug_degem"`
			ManufactureCountry    string  `json:"tozeret_nm"`
			ModelSerialNumber     int     `json:"degem_cd"`
			ModelCode             string  `json:"degem_nm"`
			TrimLevel             string  `json:"ramat_gimur"`
			SafetyFeaturesLevel   any     `json:"ramat_eivzur_betihuty"`
			PollutionLevel        int     `json:"kvutzat_zihum"`
			ManufacturYear        int     `json:"shnat_yitzur"`
			EngineSerialNumber    string  `json:"degem_manoa"`
			LastTestDate          string  `json:"mivchan_acharon_dt"`
			ValidDate             string  `json:"tokef_dt"`
			Ownership             string  `json:"baalut"`
			FrameNumber           string  `json:"misgeret"`
			ColorCode             int     `json:"tzeva_cd"`
			Color                 string  `json:"tzeva_rechev"`
			FrontWheel            string  `json:"zmig_kidmi"`
			RearWheel             string  `json:"zmig_ahori"`
			FuelType              string  `json:"sug_delek_nm"`
			RegisterySerialNumber int     `json:"horaat_rishum"`
			FirstOnRoadDate       string  `json:"moed_aliya_lakvish"`
			CommercialName        string  `json:"kinuy_mishari"`
			Rank                  float64 `json:"rank"`
		} `json:"records"`
		Fields []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"fields"`
		Links struct {
			Start string `json:"start"`
			Next  string `json:"next"`
		} `json:"_links"`
		Total             int  `json:"total"`
		TotalWasEstimated bool `json:"total_was_estimated"`
	} `json:"result"`
}