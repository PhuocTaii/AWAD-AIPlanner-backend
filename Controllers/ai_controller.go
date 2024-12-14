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

	type TaskResponse struct {
		ID          int    `json:"id"`
		Day         string `json:"day"`
		Description string `json:"description"`
		Subject     string `json:"subject"`
	}

	var taskResponses []TaskResponse

	t1 := TaskResponse{
		ID:          1,
		Day:         "Monday 02/12/2024",
		Description: "Create code base",
		Subject:     "Software Engineering",
	}

	t2 := TaskResponse{
		ID:          2,
		Day:         "Tuesday 03/12/2024",
		Description: "Hosting db",
		Subject:     "Database Management",
	}

	t3 := TaskResponse{
		ID:          3,
		Day:         "Wednesday 04/12/2024",
		Description: "Create CRUD user",
		Subject:     "Software Engineering",
	}

	taskResponses = append(taskResponses, t1)
	taskResponses = append(taskResponses, t2)
	taskResponses = append(taskResponses, t3)

	//parse taskRepsonses to string
	taskResponsesStr := ""
	for _, task := range taskResponses {
		taskResponsesStr += fmt.Sprintf("ID: %d, Day: %s, Description: %s\n", task.ID, task.Day, task.Description)
	}

	fmt.Println(taskResponses)

	textPromt := "bạn là chuyên gia trong việc lên kế hoạch học tập, và bạn sẽ đánh giá kế hoạch sau đây và đưa ra nhận xét" + taskResponsesStr

	client, err := genai.NewClient(c, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL"))
	resp, err := model.GenerateContent(c, genai.Text(textPromt))
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, resp)
}
