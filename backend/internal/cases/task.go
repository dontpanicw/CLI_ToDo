package cases

import (
	"CLI_ToDo/backend/internal/entity"
	"CLI_ToDo/backend/internal/port"
	"errors"
	"strings"
)

type taskUseCase struct {
	TaskRepo port.TaskRepo
}

// NewTaskUseCase создает новый экземпляр TaskUseCase
func NewTaskUseCase(repo port.TaskRepo) port.TaskUseCase {
	return &taskUseCase{
		TaskRepo: repo,
	}
}

func (t taskUseCase) Create(name, description string) (int, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	return t.TaskRepo.Create(name, description)
}

func (t taskUseCase) Read(id int) (*entity.Task, error) {
	if id <= 0 {
		return nil, errors.New("id must be positive")
	}
	return t.TaskRepo.Read(id)
}

func (t taskUseCase) Update(id int, name, description string) (int, error) {
	if id <= 0 {
		return 0, errors.New("id must be positive")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	return t.TaskRepo.Update(id, name, description)
}

func (t taskUseCase) Delete(id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	return t.TaskRepo.Delete(id)
}
