package services

import (
	"crudspanner/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Get(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Update(id uint, user *model.User) (*model.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// func TestUserService_Registration(t *testing.T) {
// 	mockRepo := new(MockUserRepository)
// 	userService := NewUserService(mockRepo)

// 	// Test case: valid registration
// 	newUser := &model.User{
// 		Model:    gorm.Model{ID: 1},
// 		Name:     "John",
// 		Email:    "john@example.com",
// 		Password: "password123",
// 	}

// 	// Mock that FindByEmail will return nil (no existing user)
// 	mockRepo.On("FindByEmail", newUser.Email).Return(nil, errors.New("record not found"))

// 	// Mock Create to return the new user
// 	mockRepo.On("Create", mock.Anything).Return(newUser, nil)

// 	result, err := userService.Registration(newUser)

// 	// Validate if everything went as expected
// 	assert.NoError(t, err)
// 	assert.Equal(t, newUser.Name, result.Name)
// 	assert.Equal(t, newUser.Email, result.Email)

// 	// Ensure expectations were met
// 	mockRepo.AssertExpectations(t)

// 	// Test case: registration with missing required field
// 	invalidUser := &model.User{
// 		Name:  "",
// 		Email: "invalid@example.com",
// 	}
// 	_, err = userService.Registration(invalidUser)
// 	assert.Error(t, err)
// 	assert.Equal(t, "name, email, and password are required", err.Error())

// 	mockRepo.AssertExpectations(t)
// }

// func TestUserServiceError(t *testing.T) {
// 	mockRepo := new(MockUserRepository)
// 	userService := NewUserService(mockRepo)
// 	newUser := &model.User{
// 		Model:    gorm.Model{ID: 1},
// 		Name:     "John",
// 		Email:    "john@example.com",
// 		Password: "password123",
// 	}

//		mockRepo.On("FindByEmail", newUser.Email).Return(newUser, nil)
//		_, err := userService.Registration(newUser)
//		assert.Error(t, err)
//		assert.Equal(t, "user with this email already exists", err.Error())
//	}
func TestUserService_Registration(t *testing.T) {
	tests := []struct {
		name               string
		mockFindByEmail    func(mockRepo *MockUserRepository)
		mockCreate         func(mockRepo *MockUserRepository)
		user               *model.User
		expectedError      string
		expectedUserReturn *model.User
	}{
		{
			name: "valid registration",
			mockFindByEmail: func(mockRepo *MockUserRepository) {
				// Mock that FindByEmail will return nil (no existing user)
				mockRepo.On("FindByEmail", "john@example.com").Return(nil, errors.New("record not found"))
			},
			mockCreate: func(mockRepo *MockUserRepository) {
				// Mock Create to return the new user
				mockRepo.On("Create", mock.Anything).Return(&model.User{
					Model:    gorm.Model{ID: 1},
					Name:     "John",
					Email:    "john@example.com",
					Password: "password123",
				}, nil)
			},
			user: &model.User{
				Name:     "John",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedError: "",
			expectedUserReturn: &model.User{
				Model:    gorm.Model{ID: 1},
				Name:     "John",
				Email:    "john@example.com",
				Password: "password123",
			},
		},
		{
			name: "registration with missing required field",
			mockFindByEmail: func(mockRepo *MockUserRepository) {
				// No specific mock needed here, as this is an invalid input case
			},
			mockCreate: func(mockRepo *MockUserRepository) {
				// No need to mock create since the test will fail before hitting it
			},
			user: &model.User{
				Name:  "",
				Email: "invalid@example.com",
			},
			expectedError: "name, email, and password are required",
		},
		{
			name: "user already exists",
			mockFindByEmail: func(mockRepo *MockUserRepository) {
				// Mock FindByEmail to return an existing user
				mockRepo.On("FindByEmail", "john@example.com").Return(&model.User{
					Name:     "John",
					Email:    "john@example.com",
					Password: "password123",
				}, nil)
			},
			mockCreate: func(mockRepo *MockUserRepository) {
				// No need to mock Create since this will fail due to existing user
			},
			user: &model.User{
				Name:     "John",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedError: "user with this email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			userService := NewUserService(mockRepo)

			// Apply mock setup
			tt.mockFindByEmail(mockRepo)
			tt.mockCreate(mockRepo)

			// Call the Registration method
			result, err := userService.Registration(tt.user)

			// Check expected error
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserReturn.Name, result.Name)
				assert.Equal(t, tt.expectedUserReturn.Email, result.Email)
			}

			// Ensure expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
