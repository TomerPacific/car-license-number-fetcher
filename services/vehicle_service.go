package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	config "car-license-number-fetcher/config"
	vehicle "car-license-number-fetcher/models"
	"car-license-number-fetcher/utils"
)

// FetchVehicleDetailsByLicensePlate fetches and processes vehicle details from the API
// based on the provided license plate. Returns the processed vehicle response or an error.
func FetchVehicleDetailsByLicensePlate(licensePlate string) (vehicle.VehicleResponse, error) {
	requestUrl := fmt.Sprintf("%s%s", config.VehicleDataAPIEndpoint, licensePlate)

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
	splitManufactureCountryCharacter := utils.GetSplitCharacter(record.ManufactureCountry)
	manufacturerCountryAndName := strings.Split(record.ManufactureCountry, splitManufactureCountryCharacter)

	safetyFeaturesLevel, conversionError := utils.ParseSafetyFeaturesLevelField(record)
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
