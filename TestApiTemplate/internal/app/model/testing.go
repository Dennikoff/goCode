package model

func TestUser() *User {
	return &User{
		Email:    "user@example.com",
		Password: "pass111",
		Name:     "example",
	}
}
