package handlers

import (
	"errors"
	"net/http"

	config "car-license-number-fetcher/config"
	"car-license-number-fetcher/services"
	"car-license-number-fetcher/utils"

	"github.com/gin-gonic/gin"
)

// GetVehiclePlateNumber handles GET /vehicle/:licensePlate
func GetVehiclePlateNumber(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licensePlate := c.Param(config.LicensePlateKey)
	if licensePlate == "" {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licensePlate)
		return
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}

// GetTirePressure handles GET /tire-pressure/:licensePlate
func GetTirePressure(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licensePlate := c.Param(config.LicensePlateKey)
	if licensePlate == "" {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licensePlate)
		return
	}

	// TODO: Use vehicleDetails to fetch tire pressure information
	_ = vehicleDetails
}
