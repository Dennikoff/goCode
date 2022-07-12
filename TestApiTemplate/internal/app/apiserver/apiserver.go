package apiserver

import (
	"TestProj/internal/app/store/sqlstore"
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	st := sqlstore.New(db)

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	srv := newServer(st, sessionStore)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
