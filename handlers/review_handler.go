package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	config "car-license-number-fetcher/config"
	"car-license-number-fetcher/utils"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GetVehicleReview handles GET /review/:vehicleName
func GetVehicleReview(c *gin.Context) {
	if !utils.IsRequestFromMobile(c.Request.UserAgent()) {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("request is not from a mobile device"))
		return
	}

	vehicleName, err := url.QueryUnescape(c.Param(config.VehicleNameKey))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	if vehicleName == "" {
		utils.RespondWithError(c, http.StatusBadRequest, errors.New("vehicle name was not found in request"))
		return
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv(config.OpenAIAPIKeyEnvVar)))

	language := c.GetHeader("Accept-Language")
	question := utils.GetQuestionBasedOnLocale(language, vehicleName)

	completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, completion.Choices[0].Message.Content)
}
