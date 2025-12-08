package port

import "CLI_ToDo/backend/internal/entity"

type TaskUseCase interface {
	Create(name, description string) (int, error)
	Read(Id int) (*entity.Task, error)
	Update(id int, name, description string) (int, error)
	Delete(id int) error
}
