package auth

import models "project/Models"

type RegisterResponse struct {
	User         *models.User         `json:"user"`
	TimerSetting *models.TimerSetting `json:"timer_setting"`
}
