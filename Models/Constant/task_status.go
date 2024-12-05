package constant

import "fmt"

const (
	ToDo = iota
	InProgress
	Completed
	Expired
)

func StatusToString(status int) string {
	switch status {
	case ToDo:
		return "ToDo"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	case Expired:
		return "Expired"
	default:
		return "Unknown"
	}
}

func StringToStatus(statusStr string) (int, error) {
	switch statusStr {
	case "ToDo":
		return ToDo, nil
	case "InProgress":
		return InProgress, nil
	case "Completed":
		return Completed, nil
	case "Expired":
		return Expired, nil
	default:
		return -1, fmt.Errorf("invalid status: %s", statusStr)
	}
}
