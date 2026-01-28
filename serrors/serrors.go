package serrors

import "errors"

// Sentinel errors used across the vehicle-details flow.
// Keep messages short and stable (they are not the user-facing text).
var (
    ErrFetchLicensePlate          = errors.New("fetch license plate")
    ErrParseResponse              = errors.New("parse response")
    ErrResponseNotSuccessful      = errors.New("response not successful")
    ErrNoMatchingVehicle          = errors.New("no matching vehicle")
    ErrConvertSafetyFeaturesLevel = errors.New("convert safetyFeaturesLevel")
    ErrFetchTirePressure          = errors.New("fetch tire pressure")
)