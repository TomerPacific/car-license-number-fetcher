package utils

import (
	"log"
	"strings"
	"sync"
	"unicode"
)

var carManufacturerNameTranslationToEnglish = map[string]string{
	"אאודי":             "Audi",
	"אבארט":             "Abarth",
	"אווטאר":            "Avatar",
	"אוטוביאנקי":        "Autobianchi",
	"איוויס":            "Aiways",
	"אי.וי.איזי":        "A.V.EZ",
	"אופל":              "Opel",
	"אורה":              "Ora",
	"איווקו":            "Iveco",
	"אניאוס":            "Ineos",
	"איסוזו":            "Isuzu",
	"אינפיניטי":         "Infiniti",
	"אלפא רומיאו":       "Alfa Romeo",
	"אם.ג'י":            "MG",
	"אסטון מרטין":       "Aston Martin",
	"אל.אי.וי.סי":       "L.E.V.C",
	"אל.טי.איי":         "L.T.I",
	"אלפין":             "Alpine",
	"אם דאבל יו אם":     "M.W.M",
	"אקורה":             "Acura",
	"אקס אי וי":         "XEV",
	"אקספנג":            "XPeng",
	"ב.מ.וו":            "BMW",
	"בי.אי.דבאליו":      "B.I.W",
	"ביואיק":            "Buick",
	"בנטלי":             "Bentley",
	"ג'אקו":             "Jaecoo",
	"ג'י.איי.סי":        "GAC",
	"ג'י.אם.סי":         "GMC",
	"ג'ילי":             "Geely",
	"ג'נסיס":            "Genesis",
	"גופיל":             "Goupil",
	"ג'יפ":              "Jeep",
	"גיאיוואן":          "Gyon",
	"גרייט וול":         "Great Wall",
	"ג'יי.איי.סי":       "JAC",
	"דאבל יו אם מוטורס": "W.M. Motors",
	"דאצ'יה":            "Dacia",
	"דודג'":             "Dodge",
	"דונגפנג":           "Dongfeng",
	"די.אס":             "DS",
	"דייהו":             "Daewoo",
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
	"לינק אנד קו":       "Lynk & Co",
	"למבורגיני":         "Lamborghini",
	"לנד רובר":          "Land Rover",
	"לקסוס":             "Lexus",
	"לינקולן":           "Lincoln",
	"ליפמוטור":          "Leapmotor",
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
	"סאנגיונג":          "SsangYong",
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
	"פולקסווגן":         "Volkswagen",
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

var englishManufacturerToLegacySlug = map[string]string{
	"Mercedes":    "mercedes-benz",
	"Alfa Romeo":  "alfa-romeo",
	"Land Rover":  "land-rover",
	"Rolls-Royce": "rolls-royce",
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

	if isLatinManufacturerName(manufacturerName) {
		return slugifyManufacturerName(manufacturerName)
	}

	normalizedManufacturerName := normalizeManufacturerName(manufacturerName)
	if translatedName, found := carManufacturerNameTranslationToEnglish[normalizedManufacturerName]; found {
		return slugifyManufacturerName(translatedName)
	}

	unknownManufacturersMu.Lock()
	if _, seen := unknownManufacturers[manufacturerName]; !seen {
		unknownManufacturers[manufacturerName] = struct{}{}
		log.Printf("ConvertManufacturerToEnglish: unmapped Hebrew manufacturer: %q — consider adding translation support", manufacturerName)
	}
	unknownManufacturersMu.Unlock()

	return ""
}

func TranslateManufacturerNameToEnglish(manufacturerName string) string {
	manufacturerName = strings.TrimSpace(manufacturerName)
	if manufacturerName == "" {
		return ""
	}

	normalizedManufacturerName := normalizeManufacturerName(manufacturerName)
	if translatedName, found := carManufacturerNameTranslationToEnglish[normalizedManufacturerName]; found {
		return translatedName
	}

	if isLatinManufacturerName(manufacturerName) {
		return manufacturerName
	}

	return ""
}

func isLatinManufacturerName(s string) bool {
	hasLatinLetter := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r):
			if !unicode.In(r, unicode.Latin) {
				return false
			}
			hasLatinLetter = true
		case unicode.IsDigit(r):
			continue
		case unicode.IsSpace(r):
			continue
		case strings.ContainsRune("&-.'/", r):
			continue
		default:
			return false
		}
	}
	return hasLatinLetter
}

func slugifyManufacturerName(manufacturerName string) string {
	if slug, found := englishManufacturerToLegacySlug[manufacturerName]; found {
		return slug
	}

	slug := strings.ToLower(strings.TrimSpace(manufacturerName))
	slug = strings.ReplaceAll(slug, "&", " and ")
	slug = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			return r
		case unicode.IsSpace(r), r == '-', r == '.', r == '\'' || r == '/':
			return '-'
		default:
			return -1
		}
	}, slug)

	return strings.Join(strings.FieldsFunc(slug, func(r rune) bool { return r == '-' }), "-")
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
