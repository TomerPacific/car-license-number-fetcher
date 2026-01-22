package main

import (
	"log"

	"car-license-number-fetcher/handlers"
	"car-license-number-fetcher/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Register routes
	router.GET("/vehicle/:licensePlate", handlers.GetVehiclePlateNumber)
	router.GET("/review/:vehicleName", handlers.GetVehicleReview)
	router.GET("/tire-pressure/:licensePlate", handlers.GetTirePressure)

	port := utils.GetPort()

	if runningServerError := router.Run(":" + port); runningServerError != nil {
		log.Fatalf("Running server encountered an error: %s", runningServerError)
	}
}
