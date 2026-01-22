package utils

import (
	"fmt"
	"strconv"
	"strings"

	vehicle "car-license-number-fetcher/models"
)

// GetSplitCharacter determines the split character for manufacture country
// Manufacture country can sometimes be separated by a dash or by a space
func GetSplitCharacter(country string) string {
	if strings.Contains(country, "-") {
		return "-"
	}
	return " "
}

// ParseSafetyFeaturesLevelField parses the safety features level field from a vehicle record
func ParseSafetyFeaturesLevelField(record vehicle.VehicleRecord) (int, error) {
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

// GetQuestionBasedOnLocale returns a localized question for vehicle review
func GetQuestionBasedOnLocale(language string, vehicleName string) string {
	if strings.HasPrefix(language, "en") {
		return fmt.Sprintf("Give a pros and cons list of %s", vehicleName)
	}

	return fmt.Sprintf("תן רשימה של יתרונות וחסרונות של %s", vehicleName)
}
