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

	port := getPort()

	if runningServerError := router.Run(":" + port); runningServerError != nil {
		log.Fatalf("Running server encountered an error: %s", runningServerError)
	}
}

func getVehiclePlateNumber(c *gin.Context) {

	if !isRequestFromMobile(c.Request.UserAgent()) {
		respondWithError(c, http.StatusBadRequest, "request is not from a mobile device")
		return
	}

	licensePlate := c.Param(licensePlateKey)

	if licensePlate == "" {
		respondWithError(c, http.StatusBadRequest, "license Plate was not found in request")
		return
	}

	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		respondWithError(c, http.StatusBadGateway, fmt.Sprintf("error fetching license plate: %s", requestError))
		return
	}

	defer res.Body.Close()

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("error parsing response: %s", readingResponseError))
		return
	}

	var v vehicle.VehicleDetails
	if convertingToJsonError := json.Unmarshal(resBody, &v); convertingToJsonError != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("error converting response: %s", convertingToJsonError))
		return
	}

	if !v.Success {
		respondWithError(c, http.StatusNotFound, "response was not successful")
		return
	}

	records := v.Result.Records

	if len(records) == 0 {
		respondWithError(c, http.StatusNotFound, fmt.Sprintf("no matching vehicle for the license plate entered %s", licensePlate))
		return
	}

	var record = records[0]
	splitManufactureCountryCharacter := getSplitCharacter(record.ManufactureCountry)
	manufacturerCountryAndName := strings.Split(record.ManufactureCountry, splitManufactureCountryCharacter)

	safetyFeaturesLevel, conversionError := parseSafetyFeaturesLevelField(record)
	if conversionError != nil {
		respondWithError(c, http.StatusNotFound, fmt.Sprintf("error converting safetyFeaturesLevel from string to int %s", conversionError))
		return
	}

	vehicleDetails := vehicle.VehicleResponse{
		LicenseNumber:       record.LicenseNumber,
		ManufacturerCountry: manufacturerCountryAndName[1],
		TrimLevel:           record.TrimLevel,
		SafetyFeaturesLevel: safetyFeatuesLevel,
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
		respondWithError(c, http.StatusBadRequest, "request is not from a mobile device")
		return
	}

	vehicleName, error := url.QueryUnescape(c.Param(vehicleNameKey))

	if error != nil {
		respondWithError(c, http.StatusBadRequest, error.Error())
		return
	}

	if vehicleName == "" {
		respondWithError(c, http.StatusBadRequest, "vehicle name was not found in request")
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
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, completion.Choices[0].Message.Content)

}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{errorKey: message})
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
	var safetyFeatuesLevel = 0
	if record.SafetyFeaturesLevel != nil {
		safetyFeaturesLevelString, ok := record.SafetyFeaturesLevel.(string)
		if !ok {
			safetyFeatuesLevel = record.SafetyFeaturesLevel.(int)
		} else {
			convertedSafetyFeaturesLevel, conversionError := strconv.Atoi(safetyFeaturesLevelString)
			if conversionError != nil {
				return -1, conversionError
			}
			safetyFeatuesLevel = convertedSafetyFeaturesLevel
		}
	}

	return safetyFeatuesLevel, nil
}
