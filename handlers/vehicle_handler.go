package handlers

import (
	"errors"
	"fmt"
	"net/http"

	config "car-license-number-fetcher/config"
	serrors "car-license-number-fetcher/serrors"
	"car-license-number-fetcher/services"
	"car-license-number-fetcher/utils"

	"github.com/gin-gonic/gin"
)

func GetVehiclePlateNumber(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licenseplate := c.Param(config.LicensePlateKey)
	if licenseplate == "" {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licenseplate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licenseplate)
		return
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}

func GetTirePressure(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	licenseplate := c.Param(config.LicensePlateKey)
	if licenseplate == "" {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("license Plate was not found in request"))
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licenseplate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licenseplate)
		return
	}

	tirePressureResponse, err := services.FetchTirePressureByVehicleDetails(vehicleDetails)
	if err != nil {
		utils.RespondWithError(
			c, 
			http.StatusInternalServerError,
			fmt.Errorf("%w: %v", serrors.ErrFetchTirePressure, err))
		return
	}

	c.IndentedJSON(http.StatusOK, tirePressureResponse)
}
