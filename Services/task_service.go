package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	constant "project/Models/Constant"
	task "project/Models/Request/Task"
	responseTask "project/Models/Response/Task"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateTask(c *gin.Context, request task.CreateTaskRequest) (*responseTask.GetTaskResponse, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}
	if (request.EstimatedEndTime != nil) && (request.EstimatedStartTime != nil) {
		if request.EstimatedStartTime.Unix() > request.EstimatedEndTime.Unix() {
			return nil, &config.APIError{
				Code:    http.StatusBadRequest,
				Message: "Invalid estimated start time and estimated end time",
			}
		}
	}

	subject, _ := FindSubjectByIdAndUserId(c, request.SubjectId, curUser.ID.Hex())

	priority, _ := constant.StringToPriority(request.Priority)
	if priority == -1 {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid priority",
		}
	}

	status, _ := constant.StringToStatus(request.Status)
	if status == -1 {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid status",
		}
	}
	// Create task
	task := &models.Task{
		Name:               request.Name,
		Description:        request.Description,
		User:               &curUser.ID,
		Priority:           priority,
		Status:             status,
		EstimatedStartTime: request.EstimatedStartTime,
		EstimatedEndTime:   request.EstimatedEndTime,
	}

	if subject != nil {
		task.Subject = &subject.ID
	} else {
		task.Subject = nil
	}

	// Insert task
	res, _ := repository.InsertTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to create task",
		}
	}

	resUser, _ := repository.FindUserById(c, res.User.Hex())

	var resSubject *models.Subject = nil
	if res.Subject != nil {
		resSubject, _ = repository.FindSubjectById(c, res.Subject.Hex())
	}

	response := &responseTask.GetTaskResponse{
		ID:                 res.ID,
		Name:               res.Name,
		Description:        res.Description,
		Subject:            resSubject,
		User:               resUser,
		Priority:           constant.PriorityToString(res.Priority),
		Status:             constant.StatusToString(res.Status),
		EstimatedStartTime: res.EstimatedStartTime,
		EstimatedEndTime:   res.EstimatedEndTime,
		ActualStartTime:    res.ActualStartTime,
		ActualEndTime:      res.ActualEndTime,
		FocusTime:          res.FocusTime,
		IsDeleted:          res.IsDeleted,
		CreatedAt:          res.CreatedAt,
		UpdatedAt:          res.UpdatedAt,
	}

	return response, nil
}

func ModifyTask(c *gin.Context, id string, request task.ModifyTaskRequest) (*responseTask.GetTaskResponse, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	task, _ := repository.FindTaskByIdAndUserId(c, id, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	//cannot modify task to expired
	if request.Status == "Expired" && task.Status != constant.Expired {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Cannot modify task to expired",
		}
	}

	if (request.Name == "") && (request.Description == "") && (request.SubjectId == "") && (request.Priority == "") && (request.Status == "") && (request.EstimatedStartTime == nil) && (request.EstimatedEndTime == nil) {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
	}

	if request.EstimatedEndTime != nil && request.EstimatedStartTime != nil {
		if request.EstimatedStartTime.Unix() > request.EstimatedEndTime.Unix() {
			return nil, &config.APIError{
				Code:    http.StatusBadRequest,
				Message: "Invalid estimated start time and estimated end time",
			}
		}
	}

	// var subject *models.Subject

	subject, _ := FindSubjectByIdAndUserId(c, request.SubjectId, curUser.ID.Hex())

	priority, _ := constant.StringToPriority(request.Priority)
	if priority == -1 {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid priority",
		}
	}

	status, _ := constant.StringToStatus(request.Status)
	if status == -1 {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid status",
		}
	}

	task.Name = request.Name
	task.Description = request.Description
	task.Priority = priority
	// task.Subject = subject.ID
	task.EstimatedStartTime = request.EstimatedStartTime
	task.EstimatedEndTime = request.EstimatedEndTime

	// modify status
	if task.Status == constant.Completed && task.EstimatedEndTime.Before(*utils.GetCurrent()) && status == constant.ToDo {
		// task.Status = status
		// set estimated end time to current time + 12 hours
		newTime := utils.GetCurrent().Add(12 * 60 * 60 * 1000000000)
		task.EstimatedEndTime = &newTime
	} else if task.Status == constant.Expired {
		if task.ActualStartTime == nil || task.ActualStartTime.After(*utils.GetCurrent()) {
			task.Status = constant.ToDo
		} else {
			task.Status = constant.InProgress
		}
	} else {
		//if change status to completed, set actual end time
		if status == constant.Completed {
			task.ActualEndTime = utils.GetCurrent()
			//if change status from to do to completed, set actual start time to current time
			if task.Status == constant.ToDo {
				task.ActualStartTime = utils.GetCurrent()
			}
		}
		//if change status to to do, set actual start time and end time to nil
		if status == constant.ToDo {
			task.ActualStartTime = nil
			task.ActualEndTime = nil
		}
		//if change status to in progress, set actual start time to current time
		if status == constant.InProgress {
			task.ActualStartTime = utils.GetCurrent()
			if task.Status == constant.Completed {
				task.ActualEndTime = nil
			}
		}
	}

	task.Status = status

	if task.EstimatedEndTime != nil && task.EstimatedEndTime.Before(*utils.GetCurrent()) {
		task.Status = constant.Expired
	}

	if subject != nil {
		task.Subject = &subject.ID
	} else {
		task.Subject = nil
	}

	// Update task
	res, _ := repository.UpdateTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to update task",
		}
	}

	resUser, _ := repository.FindUserById(c, res.User.Hex())
	var resSubject *models.Subject = nil
	if res.Subject != nil {
		resSubject, _ = repository.FindSubjectById(c, res.Subject.Hex())
	}

	response := &responseTask.GetTaskResponse{
		ID:                 res.ID,
		Name:               res.Name,
		Description:        res.Description,
		Subject:            resSubject,
		User:               resUser,
		Priority:           constant.PriorityToString(res.Priority),
		Status:             constant.StatusToString(res.Status),
		EstimatedStartTime: res.EstimatedStartTime,
		EstimatedEndTime:   res.EstimatedEndTime,
		ActualStartTime:    res.ActualStartTime,
		ActualEndTime:      res.ActualEndTime,
		FocusTime:          res.FocusTime,
		IsDeleted:          res.IsDeleted,
		CreatedAt:          res.CreatedAt,
		UpdatedAt:          res.UpdatedAt,
	}

	return response, nil
}

func ModifyTaskStatus(c *gin.Context, id string, request task.ModifyTaskStatusRequest) (*responseTask.GetTaskResponse, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	task, _ := repository.FindTaskByIdAndUserId(c, id, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	//cannot modify task to expired
	if request.Status == "Expired" && task.Status != constant.Expired {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Cannot modify task to expired",
		}
	}

	status, _ := constant.StringToStatus(request.Status)
	if status == -1 {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid status",
		}
	}

	// modify status
	if task.Status == constant.Completed && task.EstimatedEndTime.Before(*utils.GetCurrent()) && status == constant.ToDo {
		// task.Status = status
		// set estimated end time to current time + 12 hours
		newTime := utils.GetCurrent().Add(12 * 60 * 60 * 1000000000)
		task.EstimatedEndTime = &newTime
	} else if task.Status == constant.Expired {
		if task.ActualStartTime == nil || task.ActualStartTime.After(*utils.GetCurrent()) {
			task.Status = constant.ToDo
		} else {
			task.Status = constant.InProgress
		}
	} else {
		//if change status to completed, set actual end time
		if status == constant.Completed {
			task.ActualEndTime = utils.GetCurrent()
			//if change status from to do to completed, set actual start time to current time
			if task.Status == constant.ToDo {
				task.ActualStartTime = utils.GetCurrent()
			}
		}
		//if change status to to do, set actual start time and end time to nil
		if status == constant.ToDo {
			task.ActualStartTime = nil
			task.ActualEndTime = nil
		}
		//if change status to in progress, set actual start time to current time
		if status == constant.InProgress {
			task.ActualStartTime = utils.GetCurrent()
			if task.Status == constant.Completed {
				task.ActualEndTime = nil
			}
		}
	}

	task.Status = status

	res, _ := repository.UpdateTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to update task",
		}
	}

	resUser, _ := repository.FindUserById(c, res.User.Hex())
	var resSubject *models.Subject = nil
	if res.Subject != nil {
		resSubject, _ = repository.FindSubjectById(c, res.Subject.Hex())
	}

	response := &responseTask.GetTaskResponse{
		ID:                 res.ID,
		Name:               res.Name,
		Description:        res.Description,
		Subject:            resSubject,
		User:               resUser,
		Priority:           constant.PriorityToString(res.Priority),
		Status:             constant.StatusToString(res.Status),
		EstimatedStartTime: res.EstimatedStartTime,
		EstimatedEndTime:   res.EstimatedEndTime,
		ActualStartTime:    res.ActualStartTime,
		ActualEndTime:      res.ActualEndTime,
		FocusTime:          res.FocusTime,
		IsDeleted:          res.IsDeleted,
		CreatedAt:          res.CreatedAt,
		UpdatedAt:          res.UpdatedAt,
	}

	return response, nil
}

func GetPagingTask(c *gin.Context, limit, page int, filter, sort bson.M) ([]*responseTask.GetTaskResponse, int, int, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, 0, 0, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	filter["user"] = utils.ConvertStringToObjectID(curUser.ID.Hex())

	paginateConfig := config.NewPagingConfig(c, limit, page)

	tasks, totalPages, totalItems, err := config.PaginatedFind[*models.Task](c, config.TaskCollection, paginateConfig, filter, sort)

	if err != nil {
		return nil, 0, 0, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get tasks",
		}
	}

	response := make([]*responseTask.GetTaskResponse, len(tasks))

	for index, task := range tasks {
		resUser, _ := repository.FindUserById(c, task.User.Hex())
		var resSubject *models.Subject = nil
		if task.Subject != nil {
			resSubject, _ = repository.FindSubjectById(c, task.Subject.Hex())
		}
		response[index] = &responseTask.GetTaskResponse{
			ID:                 task.ID,
			Name:               task.Name,
			Description:        task.Description,
			Subject:            resSubject,
			User:               resUser,
			Priority:           constant.PriorityToString(task.Priority),
			Status:             constant.StatusToString(task.Status),
			EstimatedStartTime: task.EstimatedStartTime,
			EstimatedEndTime:   task.EstimatedEndTime,
			ActualStartTime:    task.ActualStartTime,
			ActualEndTime:      task.ActualEndTime,
			FocusTime:          task.FocusTime,
			IsDeleted:          task.IsDeleted,
			CreatedAt:          task.CreatedAt,
			UpdatedAt:          task.UpdatedAt,
		}
	}

	return response, totalPages, totalItems, nil
}

func DeleteTask(c *gin.Context, id string) (*models.Task, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	task, _ := repository.FindTaskByIdAndUserId(c, id, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	// Update task is_deleted to true
	task.IsDeleted = true

	res, _ := repository.DeleteTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete task",
		}
	}
	return res, nil
}

func ModifyDeletedSubjectTasks(c *gin.Context, subjectId string) *config.APIError {
	err := repository.ModifyDeletedSubjectTasks(c, subjectId)
	if err != nil {
		return &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to modify tasks",
		}
	}
	return nil
}

func ModifyTaskFocus(c *gin.Context, taskId string, request task.UpdateTaskFocusRequest) (*models.Task, *config.APIError) {
	//Get current user
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	task, _ := repository.FindTaskByIdAndUserId(c, taskId, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	if task.Status != constant.InProgress {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Cannot update focus time for task that is not in progress",
		}
	}

	if request.FocusTime < 0 {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid focus time",
		}
	}

	task.FocusTime = request.FocusTime

	res, _ := repository.UpdateTaskFocus(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to update task",
		}
	}

	return res, nil
}
