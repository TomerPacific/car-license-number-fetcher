package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	config "car-license-number-fetcher/config"
	vehicle "car-license-number-fetcher/models"
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

// FetchTirePressureByVehicleDetails fetches tire pressure information from the wheel-size API
// based on the provided vehicle details. Returns the tire pressure response or an error.
func FetchTirePressureByVehicleDetails(vehicleDetails vehicle.VehicleResponse) (vehicle.TirePressureResponse, error) {
	apiKey := os.Getenv(config.WheelSizeAPIKeyEnvVar)
	if apiKey == "" {
		return vehicle.TirePressureResponse{}, fmt.Errorf("WHEEL_SIZE_KEY environment variable is not set")
	}

	// Build the request URL with query parameters
	baseURL, err := url.Parse(config.WheelSizeAPIEndpoint)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("error parsing wheel-size API endpoint: %w", err)
	}

	params := url.Values{}
	params.Add("make", vehicleDetails.ManufacturerName)
	params.Add("model", vehicleDetails.CommercialName)
	params.Add("year", fmt.Sprintf("%d", vehicleDetails.ManufacturYear))
	params.Add("region", config.WheelSizeDefaultRegion)
	params.Add("user_key", apiKey)

	// Modification is optional - only add if we have a way to determine it
	// For now, we'll make the request without it
	baseURL.RawQuery = params.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	// Execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("error fetching tire pressure: %w", err)
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return vehicle.TirePressureResponse{}, fmt.Errorf("wheel-size API returned status %d: %s", res.StatusCode, string(body))
	}

	// Read response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("error reading response: %w", err)
	}

	// Parse the response
	var wheelSizeResponse WheelSizeAPIResponse
	if err := json.Unmarshal(resBody, &wheelSizeResponse); err != nil {
		return vehicle.TirePressureResponse{}, fmt.Errorf("error parsing wheel-size API response: %w", err)
	}

	// Check if we have any data
	if len(wheelSizeResponse.Data) == 0 {
		return vehicle.TirePressureResponse{}, fmt.Errorf("no vehicle data found in wheel-size API response")
	}

	// Extract tire pressure from the first vehicle's wheels
	// Prefer stock wheels, otherwise use the first available wheel
	vehicleData := wheelSizeResponse.Data[0]
	var frontPsi, rearPsi *float64
	var foundStockWheel bool

	// First, try to find a stock wheel
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

	// If no stock wheel found, use the first wheel with tire pressure data
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

	// Build the response
	tirePressureResponse := vehicle.TirePressureResponse{
		Source:   "wheel-size.com",
		FrontPsi: frontPsi,
		RearPsi:  rearPsi,
		Unit:     "psi",
	}

	// If we found at least one pressure value, we're good
	if frontPsi != nil || rearPsi != nil {
		return tirePressureResponse, nil
	}

	// If no tire pressure found, return an error
	return vehicle.TirePressureResponse{}, fmt.Errorf("no tire pressure data found in wheel-size API response")
}
