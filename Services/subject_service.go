package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	subject "project/Models/Request/Subject"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateSubject(c *gin.Context, request subject.CreateSubjectRequest) (*models.Subject, *config.APIError) {
	//Get current user
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}
	// Create task
	subject := &models.Subject{
		Name: request.Name,
		User: *curUser,
	}

	// Insert task
	res, _ := repository.InsertSubject(c, subject)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to create subject",
		}
	}

	return res, nil
}

func FindSubjectById(c *gin.Context, id string) (*models.Subject, *config.APIError) {
	// Convert id to object id
	// Find subject by id
	subject, err := repository.FindSubjectById(c, id)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Subject not found",
		}
	}
	return subject, nil
}
