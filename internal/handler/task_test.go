package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-tasks-api/internal/enum"
	"go-tasks-api/internal/model"
	"go-tasks-api/internal/repository"
	"go-tasks-api/internal/repository/mocks"
	"go-tasks-api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type taskTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	connector *Task
	mockTasks *mocks.MockTaskConnector
	router    *chi.Mux
	recoder   *httptest.ResponseRecorder
}

func TestTaskHnadler(t *testing.T) {
	suite.Run(t, new(taskTestSuite))
}

// Setup test suite
func (s *taskTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockTasks = mocks.NewMockTaskConnector(s.ctrl)

	s.connector = NewTaskHandler(s.mockTasks)
	s.recoder = httptest.NewRecorder()
	s.router = chi.NewRouter()

	s.router.Post("/tasks", s.connector.Create)
	s.router.Get("/tasks/{id}", s.connector.Get)
	s.router.Get("/tasks", s.connector.List)
	s.router.Put("/tasks/{id}", s.connector.Update)
	s.router.Delete("/tasks/{id}", s.connector.Delete)
}

// Assert expectations
func (s *taskTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Success: A task was created
//
// Return: 201
func (s *taskTestSuite) TestCreateTaskSuccess() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/tasks",
		strings.NewReader(
			`{
				"title" : "test title",
				"description": "test description"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	s.mockTasks.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, a model.Task) error {
			// validate fields
			if a.Title != "test title" {
				return errors.New("incorrect params")
			}

			return nil
		})

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusCreated, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("id", string(resBody))
}

// BadRequest: `title` field was not passed in the request body
//
// Returns: 400
func (s *taskTestSuite) TestTaskBadRequestDocumentNumber() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/tasks",
		strings.NewReader(`{ "description": "test description" }`))
	s.Require().NoError(err)
	defer req.Body.Close()

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("title", string(resBody))
}

// InternalServerError: Task creation failed at database
//
// Return: 500
func (s *taskTestSuite) TestTaskCreateFailed() {
	mockDBError := errors.New("some db error")
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/tasks",
		strings.NewReader(
			`{
				"title" : "test title",
				"description": "test description"
			}`))

	s.Require().NoError(err)
	defer req.Body.Close()

	s.mockTasks.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, a model.Task) error {
			// validate fields
			if a.Title != "test title" {
				return errors.New("incorrect params")
			}

			return mockDBError
		})

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("internal_error", string(resBody))
}

// Success: Get task by ID
//
// Return: 200
func (s *taskTestSuite) TestGetTaskSuccess() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	expected := model.Task{
		ID:          taskID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Todo,
		CreatedAt:   time.Now(),
	}

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(expected, nil)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusOK, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	expectedJson, err := json.Marshal(expected)
	s.NoError(err)

	s.JSONEq(string(expectedJson), string(resBody))
}

// InternalServerError: Failed to get task by ID
//
// Return: 500
func (s *taskTestSuite) TestGetTaskFailure() {
	taskID := utils.GetMockUUID()
	mockDBError := errors.New("some-db-error")

	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(model.Task{}, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("internal_error", string(resBody))
}

// NotFound: Given task ID was not found
//
// Return: 404
func (s *taskTestSuite) TestGetTaskFailureNotFound() {
	taskID := utils.GetMockUUID()

	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(model.Task{}, repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("not_found", string(resBody))
}

// NotFound: Get task, incorrect URL param
//
// Return: 404
func (s *taskTestSuite) TestGetTaskFailureIncorrectURLNotFound() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks/", nil)
	s.Require().NoError(err)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)
}

// Success: List tasks
//
// Return: 200
func (s *taskTestSuite) TestListTasksSuccess() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks", nil)
	s.Require().NoError(err)

	expected := []model.Task{
		{
			ID:          utils.GetMockUUID(),
			Title:       "test title 1",
			Description: "test description 1",
			Status:      enum.Status_Todo,
			CreatedAt:   time.Now(),
		},
		{
			ID:          utils.GetMockUUID(),
			Title:       "test title 2",
			Description: "test description 2",
			Status:      enum.Status_Done,
			CreatedAt:   time.Now(),
		},
	}

	s.mockTasks.EXPECT().List(gomock.Any()).Return(expected, nil)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusOK, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	expectedJson, err := json.Marshal(expected)
	s.NoError(err)

	s.JSONEq(string(expectedJson), string(resBody))
}

// InternalServerError: List tasks, error at database
//
// Return: 500
func (s *taskTestSuite) TestListTasksFailure() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/tasks", nil)
	s.Require().NoError(err)

	mockDBError := errors.New("some-db-error")

	s.mockTasks.EXPECT().List(gomock.Any()).Return(nil, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("internal_error", string(resBody))
}

// UpdateTaskSuccess: Update task successfully
//
// Return: 200
func (s *taskTestSuite) TestUpdateTaskSuccess() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPut, "/tasks/"+taskID.String(),
		strings.NewReader(
			`{
				"title" : "updated title",
				"description": "updated description",
				"status": "done"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	now := time.Now()
	current := model.Task{
		ID:          taskID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Todo,
		CreatedAt:   now,
	}

	updatedAt := time.Now().Add(time.Hour)
	updated := model.Task{
		ID:          current.ID,
		Title:       "updated title",
		Description: "updated description",
		Status:      enum.Status_Done,
		CreatedAt:   current.CreatedAt,
		UpdatedAt:   &updatedAt,
	}

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(current, nil)
	s.mockTasks.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updated, nil)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusOK, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	expectedJson, err := json.Marshal(updated)
	s.NoError(err)

	s.JSONEq(string(expectedJson), string(resBody))
}

// UpdateTaskFailureNotFound: Update task failure, task not found
//
// Return: 404
func (s *taskTestSuite) TestUpdateTaskFailureNotFound() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPut, "/tasks/"+taskID.String(),
		strings.NewReader(
			`{
				"title" : "updated title",
				"description": "updated description",
				"status": "done"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(model.Task{}, repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("not_found", string(resBody))
}

// UpdateTaskFailureInternalError: Update task failure, internal server error
//
// Return: 500
func (s *taskTestSuite) TestUpdateTaskFailureInternalError() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPut, "/tasks/"+taskID.String(),
		strings.NewReader(
			`{
				"title" : "updated title",
				"description": "updated description",
				"status": "done"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	now := time.Now()
	mockDBError := errors.New("some-db-error")

	current := model.Task{
		ID:          taskID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Todo,
		CreatedAt:   now,
	}

	s.mockTasks.EXPECT().Get(gomock.Any(), taskID.String()).Return(current, nil)
	s.mockTasks.EXPECT().Update(gomock.Any(), gomock.Any()).Return(model.Task{}, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("internal_error", string(resBody))
}

// Success: Delete task successfully
//
// Return: 204
func (s *taskTestSuite) TestDeleteTaskSuccess() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodDelete, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	s.mockTasks.EXPECT().Delete(gomock.Any(), taskID.String()).Return(nil)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNoContent, s.recoder.Code)
}

// DeleteTaskFailureNotFound: Delete task failure, task not found
//
// Return: 404
func (s *taskTestSuite) TestDeleteTaskFailureNotFound() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodDelete, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	s.mockTasks.EXPECT().Delete(gomock.Any(), taskID.String()).Return(repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("not_found", string(resBody))
}

// DeleteTaskFailureInternalError: Delete task failure, internal server error
//
// Return: 500
func (s *taskTestSuite) TestDeleteTaskFailureInternalError() {
	taskID := utils.GetMockUUID()
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodDelete, "/tasks/"+taskID.String(), nil)
	s.Require().NoError(err)

	mockDBError := errors.New("some-db-error")

	s.mockTasks.EXPECT().Delete(gomock.Any(), taskID.String()).Return(mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("internal_error", string(resBody))
}
