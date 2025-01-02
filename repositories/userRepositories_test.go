package repositories

import (
	"crudspanner/model"
	"testing"
)

func TestCreate(t *testing.T) {

	user := model.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Address:  "test address",
		Password: "testpassword",
	}
}
