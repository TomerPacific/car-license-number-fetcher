package main

import (
	"fmt"
	"log"
	"net/http"
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

const endpoint = "https://data.gov.il/api/3/action/datastore_search"

func main() {
	http.HandleFunc("/vehicle/", vehiclePlateNumberHandler)
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

}
