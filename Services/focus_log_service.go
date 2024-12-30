package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateFocusLog(c *gin.Context, focusLog *models.FocusLog) (*models.FocusLog, *config.APIError) {
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	focusLog.User = &curUser.ID

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
