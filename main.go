package main

import (
	vehicle "car-license-number-fetcher/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	endpoint           = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="
	licensePlateKey    = "licensePlate"
	vehicleNameKey     = "vehicleName"
	defaultPort        = "8080"
	openAIAPIKeyEnvVar = "OPENAPI_KEY"
	mobileUserAgent    = "Ktor client"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/vehicle/:licensePlate", getVehiclePlateNumber)
	router.GET("/review/:vehicleName", getVehicleReview)

	port := getPort()
	if port == "" {
		port = defaultPort
	}

	if runningServerError := router.Run(":" + port); runningServerError != nil {
		log.Panicf("Running server encountered an error: %s", runningServerError)
	}
}

func getVehiclePlateNumber(c *gin.Context) {

	if !isRequestFromMobile(c.Request.UserAgent()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request is not from a mobile device"})
		return
	}

	licensePlate := c.Param(licensePlateKey)

	if licensePlate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "license Plate was not found in request"})
		return
	}

	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		c.JSON(http.StatusBadGateway,
			gin.H{"error": fmt.Sprintf("error fetching license plate: %s", requestError)})
		return
	}

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("error parsing response: %s", readingResponseError)})
		return
	}

	var v vehicle.VehicleDetails
	convertingToJsonError := json.Unmarshal(resBody, &v)

	if convertingToJsonError != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("error converting response: %s", convertingToJsonError)})
		return
	}

	if !v.Success {
		c.JSON(http.StatusNotFound, gin.H{"error": "response was not successful"})
		return
	}

	records := v.Result.Records

	if len(records) == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("no matching vehicle for the license plate enntered %s", licensePlate)})
		return
	}

	var record = v.Result.Records[0]
	var splitManufactureCountryCharacter = " "

	if strings.Contains(record.ManufactureCountry, "-") {
		splitManufactureCountryCharacter = "-"
	}

	var manufacturerCountryAndName = strings.Split(record.ManufactureCountry, splitManufactureCountryCharacter)

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

func getVehicleReview(c *gin.Context) {
	if !isRequestFromMobile(c.Request.UserAgent()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request is not from a mobile device"})
		return
	}

	vehicleName, error := url.QueryUnescape(c.Param(vehicleNameKey))

	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	if vehicleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle name was not found in request"})
		return
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv(openAIAPIKeyEnvVar)))

	question := fmt.Sprintf("תן רשימה של יתרונות וחסרונות של %s", vehicleName)

	completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, completion.Choices[0].Message.Content)

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

func isRequestFromMobile(userAgent string) bool {
	return strings.Contains(userAgent, mobileUserAgent)
}
