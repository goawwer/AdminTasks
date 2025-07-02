package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/goawwer/admintasks/internal/models"
	"github.com/google/uuid"
)

type TaskRepository struct {
	storage *Storage
}

var (
	tablename string = "tasks"
)

func (task *TaskRepository) Create(ctx context.Context, t *models.Task) (*models.Task, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (title, description)
		VALUES ($1, $2)
		RETURNING id`,
		tablename,
	)

	err := task.storage.db.QueryRowContext(
		ctx,
		query,
		t.Title,
		t.Description,
	).Scan(&t.ID)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (task *TaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s", tablename)
	rows, err := task.storage.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task

	for rows.Next() {
		t := new(models.Task)

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Done,
			&t.CreatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (task *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, bool, error) {
	query := fmt.Sprintf(`
		SELECT id, title, description, done, created_at
		FROM %s
		WHERE id = $1
	`, tablename)

	t := &models.Task{}
	err := task.storage.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return t, true, nil
}

func (task *TaskRepository) Update(ctx context.Context, id uuid.UUID, updatedTask *models.Task) (*models.Task, error) {
	taskExisted, ok, err := task.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("there is no task with that id: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s 
		SET title = $1, description = $2, done = $3
		WHERE id = $4
		RETURNING title, description, done`,
		tablename,
	)

	var result models.Task

	err = task.storage.db.QueryRow(
		query,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Done,
		id,
	).Scan(
		&result.Title,
		&result.Description,
		&result.Done,
	)

	if err != nil {
		return nil, err
	}

	result.ID = taskExisted.ID

	return &result, nil
}

func (task *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	taskExisted, ok, err := task.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE title = $1", tablename)
		_, err := task.storage.db.Exec(query, taskExisted.Title)
		if err != nil {
			return err
		}
	}

	return nil
}
