package interfaces

import "crudspanner/model"

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	Get(id uint) (*model.User, error)
	GetAll() ([]model.User, error)
	Delete(id uint) error
	Update(id uint, user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
}
