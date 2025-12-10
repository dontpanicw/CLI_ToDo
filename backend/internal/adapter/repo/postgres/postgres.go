package postgres

import (
	"CLI_ToDo/backend/internal/entity"
	"CLI_ToDo/backend/internal/port"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Проверка реализации интерфейсов
var (
	_ port.TaskRepo = (*PostgresRepository)(nil)
)

// PostgresRepository объединяет все репозитории
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) Create(name, description string) (int, error) {
	const query = `INSERT INTO tasks (name, description) VALUES ($1, $2) RETURNING id`

	var id int
	err := p.db.QueryRow(query, name, description).Scan(&id)
	return id, err
}

func (p *PostgresRepository) Read(id int) (*entity.Task, error) {
	const query = `SELECT id, name, description, completed FROM tasks WHERE id = $1`

	task := &entity.Task{}
	err := p.db.QueryRow(query, id).Scan(&task.Id, &task.Name, &task.Description, &task.Completed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return task, err
}

func (p *PostgresRepository) Update(id int, name, description string) (int, error) {
	const query = `UPDATE tasks SET name = $1, description = $2 WHERE id = $3 RETURNING id`

	var updatedID int
	err := p.db.QueryRow(query, name, description, id).Scan(&updatedID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, sql.ErrNoRows
	}
	return updatedID, err
}

func (p *PostgresRepository) Delete(id int) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (p *PostgresRepository) MarkDone(id int) error {
	const query = `UPDATE tasks SET completed = TRUE WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
