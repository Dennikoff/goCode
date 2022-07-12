package model_test

import (
	"TestProj/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser()

	assert.NoError(t, u.BeforeCreate())

	assert.NotEmpty(t, u.EncryptedPassword)

}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			"valid",
			func() *model.User {
				return model.TestUser()
			},
			true,
		},
		{
			"empty pass",
			func() *model.User {
				u := model.TestUser()
				u.Password = ""
				return u
			},
			false,
		},
		{
			"empty email",
			func() *model.User {
				u := model.TestUser()
				u.Email = ""
				return u
			},
			false,
		},
		{
			"incorrect email",
			func() *model.User {
				u := model.TestUser()
				u.Email = "d.harkeyandex.ru"
				return u
			},
			false,
		},
		{
			"encrypted password",
			func() *model.User {
				u := model.TestUser()
				u.Password = ""
				u.EncryptedPassword = "some password"
				return u
			},
			true,
		},
		{
			"pass is less than 6 symbols",
			func() *model.User {
				u := model.TestUser()
				u.Password = "12345"
				return u
			},
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}

}
