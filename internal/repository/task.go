package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-tasks-api/internal/model"
)

type taskRepo struct {
	db *sql.DB
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -package mocks -destination=./mocks/task_mock.go -source=task.go
type TaskConnector interface {
	Create(ctx context.Context, a model.Task) error
	Get(ctx context.Context, id string) (model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task model.Task) (model.Task, error)
	Delete(ctx context.Context, id string) error
}

// NewTaskRepo creates a new Task repository
func NewTaskRepo(db *sql.DB) TaskConnector {
	return &taskRepo{
		db,
	}
}

func (a *taskRepo) Create(ctx context.Context, task model.Task) error {
	insertSQL := `INSERT INTO tasks.tasks (id, title, description, created_at) values ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING;`

	_, err := a.db.ExecContext(ctx, insertSQL, task.ID.String(), task.Title, task.Description, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (a *taskRepo) Get(ctx context.Context, id string) (model.Task, error) {
	getTaskSQL := `SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks where id = $1 AND is_active = true;`

	rows := a.db.QueryRowContext(ctx, getTaskSQL, id)
	if rows.Err() != nil {
		return model.Task{}, fmt.Errorf("failed to query task: %w", rows.Err())
	}
	var task model.Task
	if err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, ErrNoRows
		}

		return model.Task{}, fmt.Errorf("failed to scan task: %w", err)
	}

	return task, nil
}

func (a *taskRepo) List(ctx context.Context) ([]model.Task, error) {
	listSQL := `SELECT id, title, description, status, created_at, updated_at FROM tasks.tasks WHERE is_active = true;`

	rows, err := a.db.QueryContext(ctx, listSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tasks, nil
}

func (a *taskRepo) Update(ctx context.Context, task model.Task) (model.Task, error) {
	updateSQL := `
		UPDATE tasks.tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, status, created_at, updated_at;
	`

	var updated model.Task
	err := a.db.QueryRowContext(
		ctx,
		updateSQL,
		task.ID.String(),
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, ErrNoRows
		}
		return model.Task{}, fmt.Errorf("failed to update task: %w", err)
	}

	return updated, nil
}

func (a *taskRepo) Delete(ctx context.Context, id string) error {
	deleteSQL := `UPDATE tasks.tasks SET is_active = false WHERE id = $1;`

	res, err := a.db.ExecContext(ctx, deleteSQL, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNoRows
	}

	return nil
}
