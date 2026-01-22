package main

import (
	vehicle "car-license-number-fetcher/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	errorKey           = "error"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/vehicle/:licensePlate", getVehiclePlateNumber)
	router.GET("/review/:vehicleName", getVehicleReview)
	router.GET("/tire-pressure/:licensePlate", getTirePressure)
	port := getPort()

	if runningServerError := router.Run(":" + port); runningServerError != nil {
		log.Fatalf("Running server encountered an error: %s", runningServerError)
	}
}

func getVehiclePlateNumber(c *gin.Context) {
	if !isRequestFromMobile(c.Request.UserAgent()) {
		respondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licensePlate := c.Param(licensePlateKey)
	if licensePlate == "" {
		respondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := fetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		handleVehicleDetailsError(c, err, licensePlate)
		return
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}

func getVehicleReview(c *gin.Context) {
	if !isRequestFromMobile(c.Request.UserAgent()) {
		respondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	vehicleName, error := url.QueryUnescape(c.Param(vehicleNameKey))

	if error != nil {
		respondWithError(c, http.StatusBadRequest, error)
		return
	}

	if vehicleName == "" {
		respondWithError(c, http.StatusBadRequest, errors.New("vehicle name was not found in request"))
		return
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv(openAIAPIKeyEnvVar)))

	question := getQuestionBasedOnLocale(c, vehicleName)

	completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})

	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, completion.Choices[0].Message.Content)

}

func getTirePressure(c *gin.Context) {
	if !isRequestFromMobile(c.Request.UserAgent()) {
		respondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licensePlate := c.Param(licensePlateKey)
	if licensePlate == "" {
		respondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := fetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		handleVehicleDetailsError(c, err, licensePlate)
		return
	}

	
	_ = vehicleDetails
}

func respondWithError(c *gin.Context, code int, error error) {
	c.JSON(code, gin.H{errorKey: error.Error()})
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

/**
 * Manufacture country can sometimes be separated by a dash or by a space
 */
func getSplitCharacter(country string) string {
	if strings.Contains(country, "-") {
		return "-"
	}
	return " "
}

func isRequestFromMobile(userAgent string) bool {
	return strings.Contains(userAgent, mobileUserAgent)
}

func parseSafetyFeaturesLevelField(record vehicle.VehicleRecord) (int, error) {
	var safetyFeaturesLevel = 0
	if record.SafetyFeaturesLevel != nil {
		switch v := record.SafetyFeaturesLevel.(type) {
		case string:
			convertedSafetyFeaturesLevel, conversionError := strconv.Atoi(v)
			if conversionError != nil {
				return -1, conversionError
			}
			safetyFeaturesLevel = convertedSafetyFeaturesLevel
		}
	}

	return safetyFeaturesLevel, nil
}

func getQuestionBasedOnLocale(c *gin.Context, vehicleName string) string {
	language := c.GetHeader("Accept-Language")
	if strings.HasPrefix(language, "en") {
		return fmt.Sprintf("Give a pros and cons list of %s", vehicleName)
	}

	return fmt.Sprintf("תן רשימה של יתרונות וחסרונות של %s", vehicleName)
}


func fetchVehicleDetailsByLicensePlate(licensePlate string) (vehicle.VehicleResponse, error) {
	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("error fetching license plate: %w", requestError)
	}
	defer res.Body.Close()

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("error parsing response: %w", readingResponseError)
	}

	var v vehicle.VehicleDetails
	if convertingToJsonError := json.Unmarshal(resBody, &v); convertingToJsonError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("error converting response: %w", convertingToJsonError)
	}

	if !v.Success {
		return vehicle.VehicleResponse{}, errors.New("response was not successful")
	}

	records := v.Result.Records
	if len(records) == 0 {
		return vehicle.VehicleResponse{}, fmt.Errorf("no matching vehicle for the license plate entered %s", licensePlate)
	}

	record := records[0]
	splitManufactureCountryCharacter := getSplitCharacter(record.ManufactureCountry)
	manufacturerCountryAndName := strings.Split(record.ManufactureCountry, splitManufactureCountryCharacter)

	safetyFeaturesLevel, conversionError := parseSafetyFeaturesLevelField(record)
	if conversionError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("error converting safetyFeaturesLevel from string to int: %w", conversionError)
	}

	vehicleDetails := vehicle.VehicleResponse{
		LicenseNumber:       record.LicenseNumber,
		ManufacturerCountry: manufacturerCountryAndName[1],
		TrimLevel:           record.TrimLevel,
		SafetyFeaturesLevel: safetyFeaturesLevel,
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

	return vehicleDetails, nil
}

// handleVehicleDetailsError maps vehicle details fetch errors to appropriate HTTP responses
func handleVehicleDetailsError(c *gin.Context, err error, licensePlate string) {
	errMsg := err.Error()
	
	if strings.Contains(errMsg, "error fetching license plate") {
		respondWithError(c, http.StatusBadGateway, err)
	} else if strings.Contains(errMsg, "error parsing response") || strings.Contains(errMsg, "error converting response") {
		respondWithError(c, http.StatusInternalServerError, err)
	} else if strings.Contains(errMsg, "response was not successful") || strings.Contains(errMsg, "no matching vehicle") {
		respondWithError(c, http.StatusNotFound, err)
	} else if strings.Contains(errMsg, "error converting safetyFeaturesLevel") {
		respondWithError(c, http.StatusNotFound, err)
	} else {
		respondWithError(c, http.StatusInternalServerError, err)
	}
}
