package utils

import (
	"log"
	"strings"
	"sync"
	"unicode"
)

var HebrewToEnglishManufacturerMap = map[string]string{
	"פורד":           "ford",
	"טויוטה":         "toyota",
	"הונדה":          "honda",
	"ניסאן":          "nissan",
	"ניסן":           "nissan",
	"מיצובישי":        "mitsubishi",
	"מיצובישי-פוג'ו":   "mitsubishi",
	"BMW":            "bmw",
	"מרצדס":          "mercedes-benz",
	"מרצדס-בנץ":       "mercedes-benz",
	"אאודי":          "audi",
	"פולקסווגן":       "volkswagen",
	"יונדאי":         "hyundai",
	"קיה":            "kia",
	"מזדה":           "mazda",
	"סובארו":         "subaru",
	"לקסוס":          "lexus",
	"אינפיניטי":       "infiniti",
	"וולוו":          "volvo",
	"פיאט":           "fiat",
	"אלפא רומיאו":     "alfa-romeo",
	"פיג'ו":          "peugeot",
	"רנו":            "renault",
	"סיטרואן":        "citroen",
	"סקודה":          "skoda",
	"סיאט":           "seat",
	"לנד רובר":       "land-rover",
	"ג'יפ":           "jeep",
	"דודג'":          "dodge",
	"שברולט":         "chevrolet",
	"קדילאק":         "cadillac",
	"לינקולן":        "lincoln",
	"פורשה":          "porsche",
	"מיני":           "mini",
	"יגואר":          "jaguar",
	"בנטלי":          "bentley",
	"רולס רויס":      "rolls-royce",
	"מזראטי":         "maserati",
	"למבורגיני":      "lamborghini",
	"פרארי":          "ferrari",
	"אופל":           "opel",
	"דאציה":          "dacia",
}

var (
	unknownManufacturersMu sync.Mutex
	unknownManufacturers   = map[string]struct{}{}
)

func ConvertManufacturerToEnglish(manufacturerName string) string {
	manufacturerName = strings.TrimSpace(manufacturerName)
	if manufacturerName == "" {
		log.Printf("ConvertManufacturerToEnglish: empty manufacturer name")
		return ""
	}

	if isEnglish(manufacturerName) {
		return strings.ToLower(manufacturerName)
	}

	if englishName, found := HebrewToEnglishManufacturerMap[manufacturerName]; found {
		return strings.ToLower(englishName)
	}

	unknownManufacturersMu.Lock()
	if _, seen := unknownManufacturers[manufacturerName]; !seen {
		unknownManufacturers[manufacturerName] = struct{}{}
		log.Printf("ConvertManufacturerToEnglish: unmapped Hebrew manufacturer: %q — consider adding to HebrewToEnglishManufacturerMap", manufacturerName)
	}
	unknownManufacturersMu.Unlock()

	return strings.ToLower(manufacturerName)
}

func isEnglish(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
