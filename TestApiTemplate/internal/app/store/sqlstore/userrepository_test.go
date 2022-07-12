package sqlstore_test

import (
	"TestProj/internal/app/model"
	"TestProj/internal/app/store"
	"TestProj/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	err := s.User().Create(model.TestUser())
	assert.NoError(t, err)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	user := model.TestUser()
	userRep := s.User()
	err := userRep.Create(user)
	assert.NoError(t, err)

	u, err := userRep.Find(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	_, err = userRep.Find(user.ID + 1)
	assert.EqualError(t, err, store.ErrorRecordNotFound.Error())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	user := s.User()
	err := user.Create(&model.User{
		Email:    "user@example.org",
		Password: "1234567",
	})
	assert.NoError(t, err)

	u, err := user.FindByEmail("user@example.org")
	assert.NoError(t, err)
	assert.NotNil(t, u)

	_, err = user.FindByEmail("user1@example.org")
	assert.EqualError(t, err, store.ErrorRecordNotFound.Error())
}
