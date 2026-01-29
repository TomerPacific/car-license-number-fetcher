package serrors

import "errors"

var (
    ErrFetchLicensePlate          = errors.New("fetch license plate")
    ErrParseResponse              = errors.New("parse response")
    ErrResponseNotSuccessful      = errors.New("response not successful")
    ErrNoMatchingVehicle          = errors.New("no matching vehicle")
    ErrConvertSafetyFeaturesLevel = errors.New("convert safetyFeaturesLevel")
    ErrFetchTirePressure          = errors.New("fetch tire pressure")
    ErrNoTirePressureData         = errors.New("no tire pressure data")
    ErrInvalidVehicleDetails      = errors.New("invalid vehicle details")
)