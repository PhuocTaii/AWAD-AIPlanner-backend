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
	if len(SubjectList) == 0 {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "No subject found",
		}
	}

	// var taskResponses [][]*ai.AiTask
	taskResponses := make([][]*ai.AiTask, 0)

	for _, subject := range SubjectList {
		filter := bson.M{"subject": utils.ConvertStringToObjectID(subject.ID.Hex()), "user": utils.ConvertStringToObjectID(curUser.ID.Hex())}
		tasks, _ := repository.GetTasks(c, filter)
		if len(tasks) == 0 {
			continue
		}
		var taskAI []*ai.AiTask
		for _, task := range tasks {
			tmpSubject, _ := repository.FindSubjectById(c, task.Subject.Hex())
			if tmpSubject == nil {
				return nil, &config.APIError{
					Code:    http.StatusInternalServerError,
					Message: "Failed to find subject",
				}
			}
			tmp := &ai.AiTask{
				Name:        task.Name,
				Description: task.Description,
				Subject:     tmpSubject.Name,
				Priority:    constant.PriorityToString(task.Priority),
				Status:      constant.StatusToString(task.Status),
				FocusTime:   task.FocusTime,
			}
			if task.EstimatedStartTime == nil {
				tmp.EstimatedStartTime = ""
			} else {
				tmp.EstimatedStartTime = task.EstimatedStartTime.Format("02-01-2006 15:04:05")
			}
			if task.EstimatedEndTime == nil {
				tmp.EstimatedEndTime = ""
			} else {
				tmp.EstimatedEndTime = task.EstimatedEndTime.Format("02-01-2006 15:04:05")
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

	textPrompt := "You are an expert in creating study plans, and you will evaluate the following plan and provide feedback. The focus time of task is not necesssary. Your feedback should on potential adjustments, such as:" +
		"Warning about overly tight schedules that may lead to burnout." +
		"Recommending prioritization changes for improved focus and balance." + string(jsonString)

	client, err := genai.NewClient(c, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create client",
		}
	}
	defer client.Close()
	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL"))
	resp, err := model.GenerateContent(c, genai.Text(textPrompt))
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate content",
		}
	}
	return resp, nil
}
