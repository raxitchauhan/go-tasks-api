package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"go-tasks-api/internal/enum"
	"go-tasks-api/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jose/go-jose/v4/testutils/require"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type taskSuite struct {
	suite.Suite
	repo TaskConnector
	db   sqlmock.Sqlmock
}

func TestTask(t *testing.T) {
	suite.Run(t, new(taskSuite))
}

func (s *taskSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.repo = NewTaskRepo(db)
	s.db = mock
}

func (s *taskSuite) TearDownTest() {
	s.NoError(s.db.ExpectationsWereMet())
}

func (s *taskSuite) TestCreateSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.New()
	request := model.Task{
		ID:        mockUUID,
		Title:     "doc",
		CreatedAt: now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks.tasks (id, title, description, created_at) 
										values ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING;`)).
		WithArgs(
			request.ID.String(),
			request.Title,
			request.Description,
			request.CreatedAt,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Create(ctx, request)
	s.NoError(err)
}

func (s *taskSuite) TestCreateError() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.New()
	mockError := errors.New("db error")

	request := model.Task{
		ID:        mockUUID,
		Title:     "doc",
		CreatedAt: now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks.tasks (id, title, description, created_at) 
										values ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING;`)).
		WithArgs(
			request.ID.String(),
			request.Title,
			request.Description,
			request.CreatedAt,
		).WillReturnError(mockError)

	err := s.repo.Create(ctx, request)
	s.Error(err)
	s.True(errors.Is(err, mockError))
}

func (s *taskSuite) TestGetTaskSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.New()

	expected := model.Task{
		ID:          mockUUID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Todo,
		CreatedAt:   now,
	}
	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks where id = $1 AND is_active = true;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"title",
					"description",
					"status",
					"created_at",
					"updated_at",
				}).
				AddRow(
					expected.ID.String(),
					expected.Title,
					expected.Description,
					expected.Status,
					expected.CreatedAt,
					expected.UpdatedAt,
				))

	got, err := s.repo.Get(ctx, mockUUID.String())
	s.NoError(err)
	s.Equal(got, expected)
}

func (s *taskSuite) TestGetTaskError() {
	ctx := context.Background()
	mockUUID := uuid.New()
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks where id = $1 AND is_active = true;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(mockError)

	got, err := s.repo.Get(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, mockError))
	s.Equal(got, model.Task{})
}

func (s *taskSuite) TestListTasksSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID1 := uuid.New()
	mockUUID2 := uuid.New()

	expected := []model.Task{
		{
			ID:          mockUUID1,
			Title:       "test title 1",
			Description: "test description 1",
			Status:      enum.Status_Todo,
			CreatedAt:   now,
		},
		{
			ID:          mockUUID2,
			Title:       "test title 2",
			Description: "test description 2",
			Status:      enum.Status_Done,
			CreatedAt:   now,
			UpdatedAt:   &now,
		},
	}
	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks WHERE is_active = true;`)).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"title",
					"description",
					"status",
					"created_at",
					"updated_at",
				}).
				AddRow(
					expected[0].ID.String(),
					expected[0].Title,
					expected[0].Description,
					expected[0].Status,
					expected[0].CreatedAt,
					expected[0].UpdatedAt,
				).
				AddRow(
					expected[1].ID.String(),
					expected[1].Title,
					expected[1].Description,
					expected[1].Status,
					expected[1].CreatedAt,
					expected[1].UpdatedAt,
				),
		)

	got, err := s.repo.List(ctx)
	s.NoError(err)
	s.Equal(got, expected)
}

func (s *taskSuite) TestListTasksError() {
	ctx := context.Background()
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks WHERE is_active = true;`)).
		WillReturnError(mockError)

	got, err := s.repo.List(ctx)
	s.Error(err)
	s.True(errors.Is(err, mockError))
	s.Nil(got)
}

func (s *taskSuite) TestListTasksEmpty() {
	ctx := context.Background()

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks WHERE is_active = true;`)).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"title",
					"description",
					"status",
					"created_at",
					"updated_at",
				}))

	got, err := s.repo.List(ctx)
	s.NoError(err)
	s.Empty(got)
}

func (s *taskSuite) TestGetTaskNoRows() {
	ctx := context.Background()
	mockUUID := uuid.New()

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks where id = $1 AND is_active = true;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(sql.ErrNoRows)

	got, err := s.repo.Get(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, sql.ErrNoRows))
	s.Equal(got, model.Task{})
}

func (s *taskSuite) TestUpdateTaskSuccess() {
	ctx := context.Background()
	mockUUID := uuid.New()
	now := time.Now()
	mockTask := model.Task{
		ID:          mockUUID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Done,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	s.db.ExpectQuery(regexp.QuoteMeta(`UPDATE tasks.tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, status, created_at, updated_at;`)).
		WithArgs(
			mockUUID.String(),
			mockTask.Title,
			mockTask.Description,
			mockTask.Status,
			mockTask.UpdatedAt,
		).
		WillReturnRows(sqlmock.NewRows(
			[]string{
				"id",
				"title",
				"description",
				"status",
				"created_at",
				"updated_at",
			}).
			AddRow(
				mockTask.ID.String(),
				mockTask.Title,
				mockTask.Description,
				mockTask.Status,
				mockTask.CreatedAt,
				mockTask.UpdatedAt,
			))

	task, err := s.repo.Update(ctx, mockTask)
	s.NoError(err)
	s.Equal(task, mockTask)
}

func (s *taskSuite) TestUpdateTaskError() {
	ctx := context.Background()
	mockUUID := uuid.New()
	now := time.Now()
	mockTask := model.Task{
		ID:          mockUUID,
		Title:       "test title",
		Description: "test description",
		Status:      enum.Status_Done,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	s.db.ExpectQuery(regexp.QuoteMeta(`UPDATE tasks.tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, status, created_at, updated_at;`)).
		WithArgs(
			mockUUID.String(),
			mockTask.Title,
			mockTask.Description,
			mockTask.Status,
			mockTask.UpdatedAt,
		).
		WillReturnError(errors.New("db error"))

	task, err := s.repo.Update(ctx, mockTask)
	s.Error(err)
	s.Equal(task, model.Task{})
}

func (s *taskSuite) TestDeleteTaskSuccess() {
	ctx := context.Background()
	mockUUID := uuid.New()

	s.db.ExpectExec(regexp.QuoteMeta(`UPDATE tasks.tasks SET is_active = false WHERE id = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Delete(ctx, mockUUID.String())
	s.NoError(err)
}

func (s *taskSuite) TestDeleteTaskError() {
	ctx := context.Background()
	mockUUID := uuid.New()

	s.db.ExpectExec(regexp.QuoteMeta(`UPDATE tasks.tasks SET is_active = false WHERE id = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(errors.New("db error"))

	err := s.repo.Delete(ctx, mockUUID.String())
	s.Error(err)
}
