package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-tasks-api/internal/enum"
	"go-tasks-api/internal/model"
	"go-tasks-api/internal/repository"
	"go-tasks-api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Task struct {
	taskRepo repository.TaskConnector
}

// NewTaskHandler creates a new Task handler
func NewTaskHandler(t repository.TaskConnector) *Task {
	return &Task{
		taskRepo: t,
	}
}

func (a *Task) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := a.taskRepo.List(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to list tasks",
			Details: err.Error(),
		})

		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}

func (a *Task) Create(w http.ResponseWriter, r *http.Request) {
	var req model.TaskCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to decode request body",
			Details: err.Error(),
		})

		return
	}

	vErr := req.Validate()
	if len(vErr) > 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.ErrorDescription{
			Code:    validationError,
			Status:  http.StatusBadRequest,
			Title:   failedToCreateTask,
			Details: "failed to validate request body",
		}, vErr...)

		return
	}

	task := model.Task{
		ID:          uuid.New(),
		Title:       utils.TrimString(req.Title),
		Description: utils.TrimString(req.Description),
		CreatedAt:   time.Now(),
	}
	if err := a.taskRepo.Create(r.Context(), task); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   failedToCreateTask,
			Details: err.Error(),
		})

		return
	}

	utils.WriteJSON(w, http.StatusCreated, model.TaskCreateResponse{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		Status:      enum.Status_Todo,
	})
}

func (a *Task) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
			Status:  http.StatusNotFound,
			Code:    notFound,
			Title:   taskNotFound,
			Details: "path param 'id' cannot be empty",
		})

		return
	}
	task, err := a.taskRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
				Status:  http.StatusNotFound,
				Code:    notFound,
				Title:   taskNotFound,
				Details: err.Error(),
			})

			return
		}

		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to get task",
			Details: err.Error(),
		})

		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

func (a *Task) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
			Status:  http.StatusNotFound,
			Code:    notFound,
			Title:   taskNotFound,
			Details: "path param 'id' cannot be empty",
		})

		return
	}

	var req model.TaskUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to decode request body",
			Details: err.Error(),
		})

		return
	}

	vErr := req.Validate()
	if len(vErr) > 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.ErrorDescription{
			Code:    validationError,
			Status:  http.StatusBadRequest,
			Title:   failedToCreateTask,
			Details: "failed to validate request body",
		}, vErr...)

		return
	}

	task, err := a.taskRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
				Status:  http.StatusNotFound,
				Code:    notFound,
				Title:   taskNotFound,
				Details: err.Error(),
			})

			return
		}

		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to get task",
			Details: err.Error(),
		})

		return
	}

	task.Title = utils.TrimString(req.Title)
	task.Description = utils.TrimString(req.Description)
	// ignore error as it is already validated
	t, _ := enum.StatusTypeString(req.Status)
	task.Status = t
	now := time.Now()
	task.UpdatedAt = &now

	task, err = a.taskRepo.Update(r.Context(), task)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to update task",
			Details: err.Error(),
		})

		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

func (a *Task) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
			Status:  http.StatusNotFound,
			Code:    notFound,
			Title:   taskNotFound,
			Details: "path param 'id' cannot be empty",
		})

		return
	}

	err := a.taskRepo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			utils.WriteJSONError(w, http.StatusNotFound, utils.ErrorDescription{
				Status:  http.StatusNotFound,
				Code:    notFound,
				Title:   taskNotFound,
				Details: err.Error(),
			})

			return
		}

		utils.WriteJSONError(w, http.StatusInternalServerError, utils.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to delete task",
			Details: err.Error(),
		})

		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
