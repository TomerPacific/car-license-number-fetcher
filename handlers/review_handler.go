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

	apiKey := os.Getenv(config.OpenAIAPIKeyEnvVar)
	if apiKey == "" {
		utils.RespondWithError(c, http.StatusInternalServerError, errors.New("OpenAI API key environment variable is not set"))
		return
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey))

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
		utils.RespondWithError(c, http.StatusBadGateway, err)
		return
	}

	if len(completion.Choices) == 0 {
		utils.RespondWithError(c, http.StatusInternalServerError, errors.New("no completion choices returned from OpenAI"))
		return
	}
	c.IndentedJSON(http.StatusOK, completion.Choices[0].Message.Content)
}
