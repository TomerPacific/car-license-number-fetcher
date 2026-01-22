package config

const (
	VehicleDataAPIEndpoint = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="
	WheelSizeAPIEndpoint   = "https://api.wheel-size.com/v2/search/by_model/"
	LicensePlateKey       = "licensePlate"
	VehicleNameKey        = "vehicleName"
	DefaultPort           = "8080"
	OpenAIAPIKeyEnvVar    = "OPENAPI_KEY"
	WheelSizeAPIKeyEnvVar = "WHEEL_SIZE_KEY"
	MobileUserAgent       = "Ktor client"
	ErrorKey              = "error"
)
