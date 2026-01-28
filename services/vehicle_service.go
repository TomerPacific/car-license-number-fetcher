package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	config "car-license-number-fetcher/config"
	vehicle "car-license-number-fetcher/models"
	serrors "car-license-number-fetcher/serrors"
	"car-license-number-fetcher/utils"
)

func FetchVehicleDetailsByLicensePlate(licensePlate string) (vehicle.VehicleResponse, error) {
	requestUrl := fmt.Sprintf("%s%s", config.VehicleDataAPIEndpoint, licensePlate)

	res, requestError := http.Get(requestUrl)
	if requestError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w: %v", serrors.ErrFetchLicensePlate, requestErr)
	}
	defer res.Body.Close()

	resBody, readingResponseError := io.ReadAll(res.Body)
	if readingResponseError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w: %v", serrors.ErrParseResponse, readingResponseErr)
	}

	var v vehicle.VehicleDetails
	if convertingToJsonError := json.Unmarshal(resBody, &v); convertingToJsonError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w: %v", serrors.ErrParseResponse, convertingToJsonErr)
	}

	if !v.Success {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w", serrors.ErrResponseNotSuccessful)
	}

	records := v.Result.Records
	if len(records) == 0 {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w: no matching vehicle for license plate %s", serrors.ErrNoMatchingVehicle, licensePlate)
	}

	record := records[0]
	splitManufactureCountryCharacter := utils.GetSplitCharacter(record.ManufactureCountry)
	manufacturerCountryAndName := strings.Split(record.ManufactureCountry, splitManufactureCountryCharacter)

	safetyFeaturesLevel, conversionError := utils.ParseSafetyFeaturesLevelField(record)
	if conversionError != nil {
		return vehicle.VehicleResponse{}, fmt.Errorf("%w: %v", serrors.ErrConvertSafetyFeaturesLevel, conversionErr)
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
