package handlers

import (
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
		utils.RespondWithError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("%w: request is not from a mobile device", serrors.ErrInvalidVehicleDetails),
		)
		return
	}

	licensePlate := c.Param(config.LicensePlateKey)
	if licensePlate == "" {
		utils.RespondWithError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("%w: license plate missing from request", serrors.ErrInvalidVehicleDetails),
		)
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licensePlate)
		return
	}

	c.IndentedJSON(http.StatusOK, vehicleDetails)
}

func GetTirePressure(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("%w: request is not from a mobile device", serrors.ErrInvalidVehicleDetails),
		)
		return
	}

	licensePlate := c.Param(config.LicensePlateKey)
	if licensePlate == "" {
		utils.RespondWithError(
			c,
			http.StatusBadRequest,
			fmt.Errorf("%w: license plate missing from request", serrors.ErrInvalidVehicleDetails),
		)
		return
	}

	vehicleDetails, err := services.FetchVehicleDetailsByLicensePlate(licensePlate)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licensePlate)
		return
	}

	tirePressureResponse, err := services.FetchTirePressureByVehicleDetails(vehicleDetails)
	if err != nil {
		utils.HandleVehicleDetailsError(c, err, licensePlate)
		return
	}

	c.IndentedJSON(http.StatusOK, tirePressureResponse)
}
