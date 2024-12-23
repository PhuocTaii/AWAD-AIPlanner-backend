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

func GetSubjects(c *gin.Context) ([]models.Subject, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}
	subjects, err := repository.FindAllUserSubject(c, curUser.ID.Hex())
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Subjects not found",
		}
	}
	return subjects, nil
}

func CreateSubject(c *gin.Context, subject *models.Subject) (*models.Subject, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	// Check if subject name is empty
	if subject.Name == "" {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Subject name is required",
		}
	}

	// Check if subject name is already exist
	IsSubjectExisted := repository.IsSubjectExisted(c, subject.Name, curUser.ID.Hex())
	if IsSubjectExisted {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Subject name is already exist",
		}
	}

	// Create task
	createSubject := &models.Subject{
		Name: subject.Name,
		User: &curUser.ID,
	}

	// Insert task
	res, _ := repository.InsertSubject(c, createSubject)
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

func FindSubjectByIdAndUserId(c *gin.Context, id, userId string) (*models.Subject, *config.APIError) {

	subject, err := repository.FindSubjectByIdAndUserId(c, id, userId)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Subject not found",
		}
	}
	return subject, nil
}

func ModifySubject(c *gin.Context, id string, request subject.ModifySubjectRequest) (*models.Subject, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	subject, _ := repository.FindSubjectByIdAndUserId(c, id, curUser.ID.Hex())
	if subject == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	subject.Name = request.Name

	// Update task
	res, _ := repository.UpdateSubject(c, subject)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to update subject",
		}
	}

	return res, nil
}

func DeleteSubject(c *gin.Context, id string) (*models.Subject, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	subject, _ := repository.FindSubjectByIdAndUserId(c, id, curUser.ID.Hex())
	if subject == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Subject not found",
		}
	}

	// Update subject is_deleted to true
	subject.IsDeleted = true

	res, _ := repository.DeleteSubject(c, subject)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete subject",
		}
	}

	// Update task
	go func() {
		if err := ModifyDeletedSubjectTasks(c, id); err != nil {
		}
	}()

	return res, nil
}
