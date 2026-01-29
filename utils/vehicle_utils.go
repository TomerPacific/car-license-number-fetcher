package utils

import (
	"fmt"
	"strconv"
	"strings"

	vehicle "car-license-number-fetcher/models"
)


// GetSplitCharacter returns the separator to use when splitting the manufacturer
// country+name string. If multiple separators are present (for example "-" and " "),
// the separator that appears first in the trimmed string is returned.
// Returns a single-character string separator (default is a single space).
func GetSplitCharacter(country string) string {
	country = strings.TrimSpace(country)
	if country == "" {
		return " "
	}

	separators := []string{"-", "–", "—", "-", " "}

	firstIdx := -1
	chosen := " "

	for _, sep := range separators {
		if idx := strings.Index(country, sep); idx >= 0 {
			if firstIdx == -1 || idx < firstIdx {
				firstIdx = idx
				chosen = sep
			}
		}
	}

	return chosen
}

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

func GetQuestionBasedOnLocale(language string, vehicleName string) string {
	if strings.HasPrefix(language, "en") {
		return fmt.Sprintf("Give a pros and cons list of %s", vehicleName)
	}

	return fmt.Sprintf("תן רשימה של יתרונות וחסרונות של %s", vehicleName)
}
