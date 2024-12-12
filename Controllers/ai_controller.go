package controllers

import (
	"fmt"
	"net/http"
	"os"
	config "project/Config"
	models "project/Models"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func CreatePlan(c *gin.Context) {
	var prompt models.Prompt

	if err := c.ShouldBindJSON(&prompt); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user data",
		}
		config.HandleError(c, error)
		return
	}
	// user := models.User{Name: request.Name, Email: request.Email, Password: request.Password}
	// newUser, timerSetting, err := services.Register(c, &user)

	// if err != nil {
	// 	defer config.HandleError(c, err)
	// 	return
	// }

	// RegisterResponse := &authResponse.RegisterResponse{
	// 	User:         newUser,
	// 	TimerSetting: timerSetting,
	// }
	client, err := genai.NewClient(c, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL"))
	resp, err := model.GenerateContent(c, genai.Text(prompt.Prompt))
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, resp)
}
