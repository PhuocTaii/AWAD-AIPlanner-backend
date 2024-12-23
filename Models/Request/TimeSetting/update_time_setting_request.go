package timesetting

type UpdateTimeSettingRequest struct {
	FocusTime  int `json:"focus_time" binding:"required"`
	ShortBreak int `json:"short_break_time" binding:"required"`
	LongBreak  int `json:"long_break_time" binding:"required"`
	Interval   int `json:"interval" binding:"required"`
}
