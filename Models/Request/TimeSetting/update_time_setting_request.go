package timesetting

type UpdateTimeSettingRequest struct {
	FocusTime  int `json:"focus_time"`
	ShortBreak int `json:"short_break_time"`
	LongBreak  int `json:"long_break_time"`
	Interval   int `json:"interval"`
}
