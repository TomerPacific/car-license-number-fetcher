package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Vehicle struct {
	ID               int       `json:"id"`
	Vehicle_Number   int       `json:"mispar_rechev"`
	Origin_Id        int       `json:"tozeret_cd"`
	Model_Type       int       `json:"sug_degem"`
	Origin_Country   string    `json:"tozeret_nm"`
	Model_Id         int       `json:"degem_cd"`
	Model_Number     int       `json:"degem_nm"`
	Trim_Level       string    `json:"ramat_gimur"`
	Safety_Level     int       `json:"ramat_eivzur_betihuty"`
	Pollution_Level  int       `json:"kvutzat_zihum"`
	Production_Year  int       `json:"shnat_yitzur"`
	Engine_Model     string    `json:"degem_manoa"`
	Last_Test_Date   time.Time `json:"mivchan_acharon_dt"`
	Expiration       time.Time `json:"tokef_dt"`
	Ownership        string    `json:"baalut"`
	Frame            string    `json:"misgeret"`
	Color_Id         string    `json:"tzeva_cd"`
	Color_Name       string    `json:"tzeva_rechev"`
	Front_Tire_Code  string    `json:"zmig_kidmi"`
	Rear_Tire_Code   string    `json:"zmig_ahori"`
	Fuel_Type        string    `json:"sug_delek_nm"`
	Registry_Command int       `json:"horaat_rishum"`
	First_On_Road    string    `json:"moed_aliya_lakvish"`
	Commerical_Name  string    `json:"kinuy_mishari"`
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
		fmt.Printf("error fetching license plate: %s\n", err)
		os.Exit(1)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Printf("Received response")
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("error parsing response: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Response body: %s", resBody)
	}
}
