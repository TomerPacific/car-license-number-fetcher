package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
			ID                  int     `json:"_id"`
			MisparRechev        int     `json:"mispar_rechev"`
			TozeretCd           int     `json:"tozeret_cd"`
			SugDegem            string  `json:"sug_degem"`
			TozeretNm           string  `json:"tozeret_nm"`
			DegemCd             int     `json:"degem_cd"`
			DegemNm             string  `json:"degem_nm"`
			RamatGimur          string  `json:"ramat_gimur"`
			RamatEivzurBetihuty any     `json:"ramat_eivzur_betihuty"`
			KvutzatZihum        int     `json:"kvutzat_zihum"`
			ShnatYitzur         int     `json:"shnat_yitzur"`
			DegemManoa          string  `json:"degem_manoa"`
			MivchanAcharonDt    string  `json:"mivchan_acharon_dt"`
			TokefDt             string  `json:"tokef_dt"`
			Baalut              string  `json:"baalut"`
			Misgeret            string  `json:"misgeret"`
			TzevaCd             int     `json:"tzeva_cd"`
			TzevaRechev         string  `json:"tzeva_rechev"`
			ZmigKidmi           string  `json:"zmig_kidmi"`
			ZmigAhori           string  `json:"zmig_ahori"`
			SugDelekNm          string  `json:"sug_delek_nm"`
			HoraatRishum        int     `json:"horaat_rishum"`
			MoedAliyaLakvish    string  `json:"moed_aliya_lakvish"`
			KinuyMishari        string  `json:"kinuy_mishari"`
			Rank                float64 `json:"rank"`
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

const endpoint = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="

func main() {
	http.HandleFunc("/vehicle/", vehiclePlateNumberHandler)
	fmt.Println("Server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func vehiclePlateNumberHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetLicensePlate(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLicensePlate(w http.ResponseWriter, r *http.Request) {
	licensePlate := r.URL.Path[len("/vehicle/"):]
	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)
	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Printf("Error fetching license plate: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error parsing response: %s\n", err)
		os.Exit(1)
	}

	var v Vehicle
	er := json.Unmarshal(resBody, &v)

	if er != nil {
		fmt.Printf("error converting response: %s\n", er)
		os.Exit(1)
	}

	if !v.Success {
		fmt.Printf("Response failure")
		os.Exit(1)
	}

	var records = v.Result.Records

	for _, record := range records {
		fmt.Printf("Vehicle License Plate: %d\n", record.MisparRechev)
	}

}
