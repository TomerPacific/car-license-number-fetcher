package main

import (
	vehicle "car-license-number-fetcher/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const endpoint = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="
const licensePlateKey = "licensePlate"
const certFile = "server.crt"
const keyFile = "server.key"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/vehicle/:licensePlate", getVehiclePlateNumber)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if runningServerError := router.RunTLS(port, certFile, keyFile); runningServerError != nil {
		log.Panicf("Running server encountered an error: %s", runningServerError)
	}
}

func getVehiclePlateNumber(c *gin.Context) {

	if !isRequestFromMobile(c.Request.UserAgent()) {
		c.JSON(http.StatusBadRequest, "Request is not from a mobile device")
		return
	}

	licensePlate := c.Param(licensePlateKey)

	if licensePlate == "" {
		c.JSON(http.StatusBadRequest, "License Plate was not found in request")
		return
	}

	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		c.JSON(http.StatusBadGateway,
			fmt.Sprintf("Error fetching license plate: %s", requestError))
		return
	}

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Error parsing response: %s", readingResponseError))
		return
	}

	var v vehicle.VehicleDetails
	convertingToJsonError := json.Unmarshal(resBody, &v)

	if convertingToJsonError != nil {
		c.JSON(http.StatusInternalServerError,
			fmt.Sprintf("Error converting response: %s", convertingToJsonError))
		return
	}

	if !v.Success {
		c.JSON(http.StatusNotFound, "Response was not successful")
		return
	}

	records := v.Result.Records

	if len(records) == 0 {
		c.JSON(http.StatusNotFound,
			fmt.Sprintf("No matching vehicle for the license plate enntered %s", licensePlate))
		return
	}

	var record = v.Result.Records[0]

	var manufacturerCountryAndName = strings.Split(record.ManufactureCountry, " ")

	vehicleDetails := vehicle.VehicleResponse{
		LicenseNumber:       record.LicenseNumber,
		ManufacturerCountry: manufacturerCountryAndName[1],
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
		ManufacturerName:    manufacturerCountryAndName[0],
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}

func isRequestFromMobile(userAgent string) bool {
	match, _ := regexp.MatchString("Android", userAgent)
	return match
}
