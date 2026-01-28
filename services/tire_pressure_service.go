package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	config "car-license-number-fetcher/config"
	vehicle "car-license-number-fetcher/models"
	serrors "car-license-number-fetcher/serrors"
	"car-license-number-fetcher/utils"
)

// WheelSizeAPIResponse represents the structure of the wheel-size API response
type WheelSizeAPIResponse struct {
	Data []WheelSizeVehicleData `json:"data"`
	Meta struct {
		Count   int `json:"count"`
		Regions map[string]int `json:"regions"`
	} `json:"meta"`
}

// WheelSizeVehicleData represents a vehicle data entry in the API response
type WheelSizeVehicleData struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Wheels  []WheelSizeWheel `json:"wheels"`
}

// WheelSizeWheel represents a wheel configuration
type WheelSizeWheel struct {
	IsStock bool `json:"is_stock"`
	Front   WheelSizeTireData `json:"front"`
	Rear    WheelSizeTireData `json:"rear"`
}

// WheelSizeTireData represents tire data for front or rear
type WheelSizeTireData struct {
	TirePressure *WheelSizeTirePressure `json:"tire_pressure"`
}

// WheelSizeTirePressure represents tire pressure values
type WheelSizeTirePressure struct {
	Bar float64 `json:"bar"`
	Psi float64 `json:"psi"`
	KPa float64 `json:"kPa"`
}

func FetchTirePressureByVehicleDetails(vehicleDetails vehicle.VehicleResponse) (vehicle.TirePressureResponse, error) {
	apiKey := os.Getenv(config.WheelSizeAPIKeyEnvVar)
	if apiKey == "" {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: %s environment variable is not set", serrors.ErrInvalidVehicleDetails, config.WheelSizeAPIKeyEnvVar)

	commercial := strings.TrimSpace(vehicleDetails.CommercialName)
	if commercial == "" {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: commercial name (model) is empty", serrors.ErrInvalidVehicleDetails)
	}
	if vehicleDetails.ManufacturYear <= 0 {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: invalid manufacture year: %d", serrors.ErrInvalidVehicleDetails, vehicleDetails.ManufacturYear)
	}


	baseURL, err := url.Parse(config.WheelSizeAPIEndpoint)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: error parsing wheel-size API endpoint: %v", serrors.ErrFetchTirePressure, err)
	}

	englishManufacturer := utils.ConvertManufacturerToEnglish(vehicleDetails.ManufacturerName)

	if englishManufacturer == "" {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: manufacturer empty or not mapped: %q", serrors.ErrInvalidVehicleDetails, vehicleDetails.ManufacturerName)
	}
	
	params := url.Values{}
	params.Add("make", englishManufacturer)
	params.Add("model", commercial)
	params.Add("year", fmt.Sprintf("%d", vehicleDetails.ManufacturYear))
	params.Add("region", config.WheelSizeDefaultRegion)
	params.Add("user_key", apiKey)


	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: error creating request: %v", serrors.ErrFetchTirePressure, err)
	}

	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: %v", serrors.ErrFetchTirePressure, err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: error reading response body: %v", serrors.ErrParseResponse, err)
	}

	if res.StatusCode != http.StatusOK {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: status %d", serrors.ErrResponseNotSuccessful, res.StatusCode)
	}

	var wheelSizeResponse WheelSizeAPIResponse
	if err := json.Unmarshal(resBody, &wheelSizeResponse); err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: %v", serrors.ErrParseResponse, err)
	}

	if len(wheelSizeResponse.Data) == 0 {
		return vehicle.TirePressureResponse{}, fmt.Errorf("%w: no vehicle data found", serrors.ErrNoTirePressureData)

	}


	vehicleData := wheelSizeResponse.Data[0]
	var frontPsi, rearPsi *float64
	var foundStockWheel bool

	for _, wheel := range vehicleData.Wheels {
		if wheel.IsStock {
			if wheel.Front.TirePressure != nil {
				psi := wheel.Front.TirePressure.Psi
				frontPsi = &psi
			}
			if wheel.Rear.TirePressure != nil {
				psi := wheel.Rear.TirePressure.Psi
				rearPsi = &psi
			}
			foundStockWheel = true
			break
		}
	}

	if !foundStockWheel && len(vehicleData.Wheels) > 0 {
		firstWheel := vehicleData.Wheels[0]
		if firstWheel.Front.TirePressure != nil {
			psi := firstWheel.Front.TirePressure.Psi
			frontPsi = &psi
		}
		if firstWheel.Rear.TirePressure != nil {
			psi := firstWheel.Rear.TirePressure.Psi
			rearPsi = &psi
		}
	}

	tirePressureResponse := vehicle.TirePressureResponse{
		Source:   "wheel-size.com",
		FrontPsi: frontPsi,
		RearPsi:  rearPsi,
		Unit:     "psi",
	}

	if frontPsi != nil || rearPsi != nil {
		return tirePressureResponse, nil
	}

	return vehicle.TirePressureResponse{}, fmt.Errorf("%w: no tire pressure values present", serrors.ErrNoTirePressureData)
}
