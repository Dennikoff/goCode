package sqlstore

import (
	"TestProj/internal/app/model"
	"TestProj/internal/app/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"Insert into users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		u.Name, u.Email, u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, email, password FROM users u where u.id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrorRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, email, password FROM users u where u.email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrorRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
