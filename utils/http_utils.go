package utils

import (
	"net/http"
	"strings"

	config "car-license-number-fetcher/config"

	"github.com/gin-gonic/gin"
)

// RespondWithError sends an error response with the specified status code
func RespondWithError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{config.ErrorKey: err.Error()})
}

// HandleVehicleDetailsError maps vehicle details fetch errors to appropriate HTTP responses
func HandleVehicleDetailsError(c *gin.Context, err error, licensePlate string) {
	errMsg := err.Error()

	if strings.Contains(errMsg, "error fetching license plate") {
		RespondWithError(c, http.StatusBadGateway, err)
	} else if strings.Contains(errMsg, "error parsing response") || strings.Contains(errMsg, "error converting response") {
		RespondWithError(c, http.StatusInternalServerError, err)
	} else if strings.Contains(errMsg, "response was not successful") || strings.Contains(errMsg, "no matching vehicle") {
		RespondWithError(c, http.StatusNotFound, err)
	} else if strings.Contains(errMsg, "error converting safetyFeaturesLevel") {
		RespondWithError(c, http.StatusNotFound, err)
	} else {
		RespondWithError(c, http.StatusInternalServerError, err)
	}
}
