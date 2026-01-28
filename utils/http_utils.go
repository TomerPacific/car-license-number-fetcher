package utils

import (
	"net/http"
	"errors"

	config "car-license-number-fetcher/config"
	serrors "car-license-number-fetcher/serrors"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{config.ErrorKey: err.Error()})
}

func HandleVehicleDetailsError(c *gin.Context, err error, licensePlate string) {
    switch {
    case errors.Is(err, serrors.ErrFetchLicensePlate):
        RespondWithError(c, http.StatusBadGateway, err)
    case errors.Is(err, serrors.ErrParseResponse),
         errors.Is(err, serrors.ErrConvertSafetyFeaturesLevel):
        RespondWithError(c, http.StatusInternalServerError, err)
    case errors.Is(err, serrors.ErrResponseNotSuccessful),
         errors.Is(err, serrors.ErrNoMatchingVehicle):
        RespondWithError(c, http.StatusNotFound, err)
    default:
        RespondWithError(c, http.StatusInternalServerError, err)
    }
}