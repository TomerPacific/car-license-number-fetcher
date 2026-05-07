package utils

import (
	"log"
	"strings"
	"sync"
	"unicode"
)

var CarManufacturerNameTranslationToEnglish = map[string]string{
	"אאודי":             "Audi",
	"אבארט":             "Abarth",
	"אווטאר":            "Avatar",
	"אוטוביאנקי":         "Autobianchi",
	"איוויס":            "Aiways",
	"אי.וי.איזי":         "A.V.EZ",
	"אופל":              "Opel",
	"אורה":              "Ora",
	"איווקו":            "Iveco",
	"אניאוס":            "Ineos",
	"איסוזו":            "Isuzu",
	"אינפיניטי":         "Infiniti",
	"אלפא רומיאו":        "Alfa Romeo",
	"אם.ג'י":            "MG",
	"אסטון מרטין":        "Aston Martin",
	"אל.אי.וי.סי":        "L.E.V.C",
	"אל.טי.איי":          "L.T.I",
	"אלפין":             "Alpine",
	"אם דאבל יו אם":      "M.W.M",
	"אקורה":             "Acura",
	"אקס אי וי":          "XEV",
	"אקספנג":            "XPeng",
	"ב.מ.וו":            "BMW",
	"בי.אי.דבאליו":       "B.I.W",
	"ביואיק":            "Buick",
	"בנטלי":             "Bentley",
	"ג'אקו":             "Jaecoo",
	"ג'י.איי.סי":         "GAC",
	"ג'י.אם.סי":          "GMC",
	"ג'ילי":             "Geely",
	"ג'נסיס":            "Genesis",
	"גופיל":             "Goupil",
	"ג'יפ":              "Jeep",
	"גיאיוואן":          "Gyon",
	"גרייט וול":          "Great Wall",
	"ג'יי.איי.סי":        "JAC",
	"דאבל יו אם מוטורס":  "W.M. Motors",
	"דאצ'יה":            "Dacia",
	"דודג'":             "Dodge",
	"דונגפנג":           "Dongfeng",
	"די.אס":             "DS",
	"דייהו":             "Daewood",
	"דייהטסו":           "Daihatsu",
	"דיפאל":             "Deepal",
	"האמר":              "Hummer",
	"הונגצ'י":           "Hongqi",
	"הונדה":             "Honda",
	"וויה":              "Voyah",
	"ווי":               "Wey",
	"וולוו":             "Volvo",
	"זיקר":              "Zeekr",
	"טאטא":              "Tata",
	"טויוטה":            "Toyota",
	"טסלה":              "Tesla",
	"יגואר":             "Jaguar",
	"יונדאי":            "Hyundai",
	"יודו":              "Yudo",
	"לאדה":              "Lada",
	"לוטוס":             "Lotus",
	"לינק אנד קו":        "Lynk & Co",
	"למבורגיני":          "Lamborghini",
	"לנד רובר":          "Land Rover",
	"לקסוס":             "Lexus",
	"לינקולן":           "Lincoln",
	"ליפמוטור":          "Leafmotor",
	"לנצ'יה":            "Lancia",
	"מזדה":              "Mazda",
	"מאן":               "Man",
	"מורגן":             "Morgan",
	"מזראטי":            "Maserati",
	"מיני":              "Mini",
	"מיצובישי":          "Mitsubishi",
	"מקלארן":            "McLaren",
	"מקסוס":             "Maxus",
	"מרצדס":             "Mercedes",
	"נטע":               "Neta",
	"ניאו":              "Nio",
	"ניסאן":             "Nissan",
	"ננג'ינג":           "Nanjing",
	"סאאב":              "Saab",
	"סאנגיונג":          "Sangyong",
	"סאנשיין":           "Sunshine",
	"סובארו":            "Subaru",
	"סוזוקי":            "Suzuki",
	"סיטרואן":           "Citroen",
	"סיאט":              "Seat",
	"סקודה":             "Skoda",
	"סרס":               "Seres",
	"סמארט":             "Smart",
	"סנטרו":             "Santro",
	"סקיוול":            "Skywell",
	"פוטון":             "Foton",
	"פיאט":              "Fiat",
	"פיג'ו":             "Peugeot",
	"פולסטאר":           "Polestar",
	"פולקסווגן":          "Volkswagen",
	"פורד":              "Ford",
	"פורשה":             "Porsche",
	"פרארי":             "Ferrari",
	"פורתינג":           "Forthing",
	"פיאג'ו":            "Piaggio",
	"צ'רי":              "Chery",
	"קאדילק":            "Cadillac",
	"קארמה":             "Karma",
	"קיה":               "Kia",
	"קופרה":             "Cupra",
	"קרייזלר":           "Chrysler",
	"ראם":               "Ram",
	"רובר":              "Rover",
	"רנו":               "Renault",
	"ריהיי":             "Reyee",
	"רולס רויס":         "Rolls-Royce",
	"שברולט":            "Chevrolet",
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

	normalizedManufacturerName := normalizeManufacturerName(manufacturerName)
	switch normalizedManufacturerName {
	case "פורד":
		return "ford"
	case "טויוטה":
		return "toyota"
	case "הונדה":
		return "honda"
	case "ניסאן":
		return "nissan"
	case "מיצובישי":
		return "mitsubishi"
	case "BMW":
		return "bmw"
	case "מרצדס":
		return "mercedes-benz"
	case "אאודי":
		return "audi"
	case "פולקסווגן":
		return "volkswagen"
	case "יונדאי":
		return "hyundai"
	case "קיה":
		return "kia"
	case "מזדה":
		return "mazda"
	case "סובארו":
		return "subaru"
	case "לקסוס":
		return "lexus"
	case "אינפיניטי":
		return "infiniti"
	case "וולוו":
		return "volvo"
	case "פיאט":
		return "fiat"
	case "אלפא רומיאו":
		return "alfa-romeo"
	case "פיג'ו":
		return "peugeot"
	case "רנו":
		return "renault"
	case "סיטרואן":
		return "citroen"
	case "סקודה":
		return "skoda"
	case "סיאט":
		return "seat"
	case "לנד רובר":
		return "land-rover"
	case "ג'יפ":
		return "jeep"
	case "דודג'":
		return "dodge"
	case "שברולט":
		return "chevrolet"
	case "קאדילק":
		return "cadillac"
	case "לינקולן":
		return "lincoln"
	case "פורשה":
		return "porsche"
	case "מיני":
		return "mini"
	case "יגואר":
		return "jaguar"
	case "בנטלי":
		return "bentley"
	case "רולס רויס":
		return "rolls-royce"
	case "מזראטי":
		return "maserati"
	case "למבורגיני":
		return "lamborghini"
	case "פרארי":
		return "ferrari"
	case "אופל":
		return "opel"
	case "דאצ'יה":
		return "dacia"
	}

	unknownManufacturersMu.Lock()
	if _, seen := unknownManufacturers[manufacturerName]; !seen {
		unknownManufacturers[manufacturerName] = struct{}{}
		log.Printf("ConvertManufacturerToEnglish: unmapped Hebrew manufacturer: %q — consider adding translation support", manufacturerName)
	}
	unknownManufacturersMu.Unlock()

	return strings.ToLower(manufacturerName)
}

func TranslateManufacturerNameToEnglish(manufacturerName string) string {
	manufacturerName = strings.TrimSpace(manufacturerName)
	if manufacturerName == "" {
		return ""
	}

	normalizedManufacturerName := normalizeManufacturerName(manufacturerName)
	if translatedName, found := CarManufacturerNameTranslationToEnglish[normalizedManufacturerName]; found {
		return translatedName
	}

	if isEnglish(manufacturerName) {
		return manufacturerName
	}

	return manufacturerName
}

func isEnglish(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func normalizeManufacturerName(manufacturerName string) string {
	switch manufacturerName {
	case "ניסן":
		return "ניסאן"
	case "מיצובישי-פוג'ו":
		return "מיצובישי"
	case "מרצדס-בנץ":
		return "מרצדס"
	case "קדילאק":
		return "קאדילק"
	case "דאציה":
		return "דאצ'יה"
	default:
		return manufacturerName
	}
}
