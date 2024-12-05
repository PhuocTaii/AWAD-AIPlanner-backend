package utils

import "time"

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
