package repositories

import (
	"crudspanner/model"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var commonUser = model.User{
	Model:    gorm.Model{ID: 1},
	Name:     "testName",
	Email:    "testEmail",
	Address:  "testAddress",
	Password: "testPassword",
}

func mockDatabase() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDb, err := gorm.Open(mysql.New((mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})), &gorm.Config{})

	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}
	return gormDb, mock
}

func TestCreateUser(t *testing.T) {

	mockDb, mock := mockDatabase()
	defer func() {
		sqlDB, _ := mockDb.DB()
		sqlDB.Close()
	}()

	// Arrange: Mock the database insert behavior
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO .*users.*`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act: Use the repository to create the user
	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.Create(&commonUser)

	// Assert: Verify the result and mock expectations
	assert.NoError(t, err)
	assert.Equal(t, commonUser.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserError(t *testing.T) {
	mockDb, mock := mockDatabase()
	defer func() {
		db, _ := mockDb.DB()
		db.Close()
	}()

	repo := NewUserRepository(mockDb)

	// Arrange: mock the database to simulate an error during Create
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	// Act: attempt to create the user
	result, err := repo.Create(&commonUser)

	// Assert: check that the error is returned
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	mockDb, mock := mockDatabase()
	defer func() {
		db, _ := mockDb.DB()
		db.Close()
	}()

	id := uint(1)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(id, 1).
		WillReturnRows(rows)
	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.Get(id)

	assert.NoError(t, err)
	assert.Equal(t, commonUser.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserError(t *testing.T) {
	mockDb, mock := mockDatabase()

	id := uint(1)

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(id, 1).
		WillReturnError(assert.AnError)
	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.Get(id)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetAllUsers(t *testing.T) {
	mockDb, mock := mockDatabase()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).AddRow(1, "testName", "testEmail", "testAddress").AddRow(2, "testName2", "testEmail2", "testAddress2")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL$").WillReturnRows(rows)

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.GetAll()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsersError(t *testing.T) {
	mockDb, mock := mockDatabase()

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL$").WillReturnError(assert.AnError)

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.GetAll()
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	mockDb, mock := mockDatabase()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	userRepository := NewUserRepository(mockDb)
	TestGetUser(t)
	result, err := userRepository.Update(1, &commonUser)
	assert.NoError(t, err)
	assert.Equal(t, commonUser.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserError(t *testing.T) {
	mockDb, mock := mockDatabase()

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.Update(1, &commonUser)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestUpdateUserErrorDB(t *testing.T) {
	mockDb, mock := mockDatabase()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` ").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.Update(1, &commonUser)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestDeleteUser(t *testing.T) {
	mockDb, mock := mockDatabase()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	userRepository := NewUserRepository(mockDb)
	err := userRepository.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserError(t *testing.T) {
	mockDb, mock := mockDatabase()

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	userRepository := NewUserRepository(mockDb)
	err := userRepository.Delete(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestDeleteUserErrorDB(t *testing.T) {
	mockDb, mock := mockDatabase()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` ").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	userRepository := NewUserRepository(mockDb)
	err := userRepository.Delete(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestFindByEmail(t *testing.T) {
	mockDb, mock := mockDatabase()

	email := "testEmail"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).AddRow(1, "testName", "testEmail", "testAddress")

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE email = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(email, 1).
		WillReturnRows(rows)

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.FindByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, commonUser.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmailError(t *testing.T) {
	mockDb, mock := mockDatabase()

	email := "testEmail"

	mock.ExpectQuery("^SELECT \\* FROM `users` WHERE email = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
		WithArgs(email, 1).
		WillReturnError(assert.AnError)

	userRepository := NewUserRepository(mockDb)
	result, err := userRepository.FindByEmail(email)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
