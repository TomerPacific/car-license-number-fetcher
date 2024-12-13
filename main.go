package main

import (
	vehicle "car-license-number-fetcher/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const endpoint = "https://data.gov.il/api/3/action/datastore_search?resource_id=053cea08-09bc-40ec-8f7a-156f0677aff3&limit=1&q="

func main() {
	router := gin.Default()
	router.GET("/vehicle/:licensePlate", getVehiclePlateNumber)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func getVehiclePlateNumber(c *gin.Context) {
	licensePlate := c.Param("licensePlate")
	requestUrl := fmt.Sprintf("%s%s", endpoint, licensePlate)

	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Printf("Error fetching license plate: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error parsing response: %s\n", err)
		os.Exit(1)
	}

	var v vehicle.Vehicle
	er := json.Unmarshal(resBody, &v)

	if er != nil {
		fmt.Printf("error converting response: %s\n", er)
		os.Exit(1)
	}

	if !v.Success {
		fmt.Printf("Response failure")
		os.Exit(1)
	}

	var records = v.Result.Records

	c.IndentedJSON(http.StatusOK, records[0])
}
