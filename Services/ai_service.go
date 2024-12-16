package services

import (
	"encoding/json"
	"net/http"
	"os"
	config "project/Config"
	constant "project/Models/Constant"
	ai "project/Models/Response/AI"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

func Feedback(c *gin.Context) (*genai.GenerateContentResponse, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	SubjectList, _ := repository.FindAllUserSubject(c, curUser.ID.Hex())

	var taskResponses [][]*ai.AiTask

	for _, subject := range SubjectList {
		filter := bson.M{"subject._id": utils.ConvertStringToObjectID(subject.ID.Hex()), "user._id": utils.ConvertStringToObjectID(curUser.ID.Hex())}
		tasks, _ := repository.GetTasks(c, filter)
		var taskAI []*ai.AiTask
		for _, task := range tasks {
			tmp := &ai.AiTask{
				Name:               task.Name,
				Description:        task.Description,
				Subject:            task.Subject.Name,
				User:               task.User.Name,
				Priority:           constant.PriorityToString(task.Priority),
				Status:             constant.StatusToString(task.Status),
				FocusTime:          task.FocusTime,
				EstimatedStartTime: task.EstimatedStartTime,
				EstimatedEndTime:   task.EstimatedEndTime,
			}
			taskAI = append(taskAI, tmp)
		}
		taskResponses = append(taskResponses, taskAI)
	}

	//Create a json string from taskResponses
	jsonString, err := json.Marshal(taskResponses)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create JSON string",
		}
	}

	textPromt := "You are an expert in creating study plans, and you will evaluate the following plan and provide feedback." + string(jsonString)

	client, err := genai.NewClient(c, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create client",
		}
	}
	defer client.Close()
	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL"))
	resp, err := model.GenerateContent(c, genai.Text(textPromt))
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate content",
		}
	}
	return resp, nil
}
