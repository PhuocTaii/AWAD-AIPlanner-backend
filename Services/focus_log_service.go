package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	requestlog "project/Models/Request/RequestLog"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateFocusLog(c *gin.Context, request requestlog.CreateRequestLogRequest) (*models.FocusLog, *config.APIError) {
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	focusLog := &models.FocusLog{
		User:      curUser.ID,
		FocusTime: request.FocusTime,
	}

	// Insert focus log
	res, _ := repository.InsertFocusLog(c, focusLog)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to create focus log",
		}
	}

	return res, nil
}
