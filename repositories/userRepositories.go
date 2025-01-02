package repositories

import (
	"crudspanner/interfaces"
	"crudspanner/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Get(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(id uint, user *model.User) (*model.User, error) {
	var existingUser *model.User
	var err error

	existingUser, err = r.Get(id)
	if err != nil {
		return nil, err
	}
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Address = user.Address
	existingUser.Password = user.Password
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	_, err := r.Get(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
