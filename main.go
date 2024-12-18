package main

import (
	vehicle "car-license-number-fetcher/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const endpoint = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="
const licensePlateKey = "licensePlate"

func main() {
	router := gin.Default()
	router.GET("/vehicle/:licensePlate", getVehiclePlateNumber)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if runningServerError := router.Run(":" + port); runningServerError != nil {
		log.Panicf("Running server encountered an error: %s", runningServerError)
	}
}

func getVehiclePlateNumber(c *gin.Context) {
	licensePlate := c.Param(licensePlateKey)

	if licensePlate == "" {
		fmt.Printf("License Plate was not found in request.")
		os.Exit(1)
	}

	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		fmt.Printf("Error fetching license plate: %s\n", requestError)
		os.Exit(1)
	}

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		fmt.Printf("Error parsing response: %s\n", readingResponseError)
		os.Exit(1)
	}

	var v vehicle.VehicleDetails
	convertingToJsonError := json.Unmarshal(resBody, &v)

	if convertingToJsonError != nil {
		fmt.Printf("Error converting response: %s\n", convertingToJsonError)
		os.Exit(1)
	}

	if !v.Success {
		fmt.Printf("Response failure")
		os.Exit(1)
	}

	records := v.Result.Records

	if len(records) == 0 {
		fmt.Printf("No matching vehicle for the license plate enntered")
		os.Exit(1)
	}

	var record = v.Result.Records[0]

	vehicleDetails := vehicle.VehicleResponse{
		LicenseNumber:       record.LicenseNumber,
		ManufactureCountry:  record.ManufactureCountry,
		TrimLevel:           record.TrimLevel,
		SafetyFeaturesLevel: record.SafetyFeaturesLevel,
		PollutionLevel:      record.PollutionLevel,
		ManufacturYear:      record.ManufacturYear,
		LastTestDate:        record.LastTestDate,
		ValidDate:           record.ValidDate,
		Ownership:           record.Ownership,
		FrameNumber:         record.FrameNumber,
		Color:               record.Color,
		FrontWheel:          record.FrontWheel,
		RearWheel:           record.RearWheel,
		FuelType:            record.FuelType,
		FirstOnRoadDate:     record.FirstOnRoadDate,
		CommercialName:      record.CommercialName,
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}
