package teststore

import (
	"TestProj/internal/app/model"
	"TestProj/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	r.users[u.Email] = u
	u.ID = len(r.users)

	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	ok := false
	var u *model.User
	for _, user := range r.users {
		if user.ID == id {
			ok = true
			u = user
			break
		}
	}
	if ok == false {
		return nil, store.ErrorRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if ok == false {
		return nil, store.ErrorRecordNotFound
	}

	return u, nil
}
