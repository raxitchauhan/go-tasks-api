package model

import (
	"time"

	"go-tasks-api/internal/enum"
	"go-tasks-api/internal/utils"

	"github.com/google/uuid"
)

type TaskCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a TaskCreateRequest) Validate() []utils.FieldError {
	err := make([]utils.FieldError, 0)
	if a.Title == "" {
		err = append(err, utils.FieldError{
			Field:   "title",
			Message: "field is required",
		})
	}

	return err
}

type TaskCreateResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Status      enum.StatusType `json:"status"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}

type TaskUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (a TaskUpdateRequest) Validate() []utils.FieldError {
	vErr := make([]utils.FieldError, 0)
	if a.Title == "" {
		vErr = append(vErr, utils.FieldError{
			Field:   "title",
			Message: "field is required",
		})
	}

	if _, err := enum.StatusTypeString(a.Status); err != nil {
		vErr = append(vErr, utils.FieldError{
			Field:   "status",
			Message: "invalid status value: " + err.Error(),
		})
	}

	return vErr
}

type Task struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Status      enum.StatusType `json:"status"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}
