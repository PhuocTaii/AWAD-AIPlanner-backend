package constant

import "fmt"

const (
	Low = iota
	Medium
	High
)

func PriorityToString(priority int) string {
	switch priority {
	case High:
		return "High"
	case Medium:
		return "Medium"
	case Low:
		return "Low"
	default:
		return "Unknown"
	}
}

func StringToPriority(priorityStr string) (int, error) {
	switch priorityStr {
	case "High":
		return High, nil
	case "Medium":
		return Medium, nil
	case "Low":
		return Low, nil
	default:
		return -1, fmt.Errorf("invalid priority: %s", priorityStr)
	}
}
