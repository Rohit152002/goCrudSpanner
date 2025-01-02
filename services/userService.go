package services

import (
	"crudspanner/model"
	"crudspanner/repositories"
)

type UserService interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserById(id uint) (*model.User, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]model.User, error)
	UpdateUser(id uint, user *model.User) (*model.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *model.User) (*model.User, error) {
	return s.repo.Create(user)
}

func (s *userService) GetUserById(id uint) (*model.User, error) {
	return s.repo.Get(id)
}

func (s *userService) UpdateUser(id uint, user *model.User) (*model.User, error) {
	return s.repo.Update(id, user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll()
}
