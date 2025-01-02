package services

import (
	"crudspanner/config"
	"crudspanner/model"
	"crudspanner/repositories"
	"errors"
	"strings"
)

type UserService interface {
	Registration(user *model.User) (*model.User, error)
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

func (s *userService) Registration(user *model.User) (*model.User, error) {

	// Validate input
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}

	// Check for existing user
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Prepare new user for creation
	var addUser model.User
	addUser.Name = strings.TrimSpace(user.Name)
	addUser.Email = strings.ToLower(strings.TrimSpace(user.Email))
	addUser.Address = strings.TrimSpace(user.Address)
	addUser.Password = config.GeneratePassword(user.Password)

	return s.repo.Create(&addUser)
}

func (s *userService) GetUserById(id uint) (*model.User, error) {

	user, err := s.repo.Get(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Redact sensitive information
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateUser(id uint, user *model.User) (*model.User, error) {
	existingUser, err := s.repo.Get(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields selectively
	if strings.TrimSpace(user.Name) != "" {
		existingUser.Name = strings.TrimSpace(user.Name)
	}
	if strings.TrimSpace(user.Email) != "" {
		existingUser.Email = strings.ToLower(strings.TrimSpace(user.Email))
	}
	if strings.TrimSpace(user.Address) != "" {
		existingUser.Address = strings.TrimSpace(user.Address)
	}

	return s.repo.Update(id, existingUser)
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists before deletion
	_, err := s.repo.Get(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Proceed with deletion
	return s.repo.Delete(id)
}

func (s *userService) GetAllUsers() ([]model.User, error) {

	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
