package apiserver

import (
	"TestProj/internal/app/model"
	"TestProj/internal/app/store/teststore"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	s := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")))
	testcases := []struct {
		name       string
		payload    interface{}
		statusCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "d.harke@yandex.ru",
				"password": "12345678",
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "invalidEmail",
			payload: map[string]string{
				"email":    "invalid@email",
				"password": "12345678",
			},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalidPassword",
			payload: map[string]string{
				"email":    "valid@email.com",
				"password": "123",
			},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalidpayload",
			payload:    3,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	s := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")))
	u := model.TestUser()
	s.store.User().Create(u)
	testcases := []struct {
		name        string
		cookieValue map[interface{}]interface{}
		statusCode  int
	}{
		{
			name: "valid",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "invalid",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID + 1,
			},
			statusCode: http.StatusUnauthorized,
		},
	}

	sc := securecookie.New([]byte("secret"), nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(SessionName, tc.cookieValue)
			fmt.Println(cookieStr)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", SessionName, cookieStr))
			fmt.Println(req.Header.Get("Cookie"))
			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}

func TestServer_HandleSessionCreate(t *testing.T) {
	s := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")))
	s.store.User().Create(&model.User{
		Email:    "d.harke@yandex.ru",
		Password: "12345678",
	})
	testcases := []struct {
		name       string
		payload    interface{}
		statusCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "d.harke@yandex.ru",
				"password": "12345678",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "invalidEmail",
			payload: map[string]string{
				"email":    "d.hark@yandex.ru",
				"password": "12345678",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "invalidPassword",
			payload: map[string]string{
				"email":    "d.harke@yandex.ru",
				"password": "123456",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "invalidpayload",
			payload:    3,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/session", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}
