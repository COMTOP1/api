package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type (
	//User represents a user object to be used when not all data is required
	User struct {
		Id    uint64 `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone string `json:"phone,omitempty"`
		Team  string `json:"team,omitempty"`
		Role  string `json:"role"`
		Image string `json:"image,omitempty"`
	}
	// UserFull represents a user and all columns
	UserFull struct {
		User      `json:"user"`
		Password  string  `json:"password,omitempty"`
		Temp      string  `json:"temp,omitempty"`
		Hash      []uint8 `json:"hash,omitempty"`
		Salt      []uint8 `json:"salt,omitempty"`
		LastLogin int64   `json:"last_login,omitempty"`
		CreatedAt int64   `json:"created_at,omitempty"`
		CreatedBy User    `json:"created_by,omitempty"`
		UpdatedAt int64   `json:"updated_at,omitempty"`
		UpdatedBy User    `json:"updated_by,omitempty"`
		DeletedAt int64   `json:"deleted_at,omitempty"`
		DeletedBy User    `json:"deleted_by,omitempty"`
	}
)

func (s *Session) UserByEmail(email, token string) (user User, err error) {
	err = s.getfToken(token, "internal/user/%s", email).Into(&user)
	return user, err
}

func (s *Session) UserByEmailFull(email, token string) (user UserFull, err error) {
	err = s.getfToken(token, "internal/user/%s/full", email).Into(&user)
	return user, err
}

func (s *Session) UserByToken(token string) (user User, err error) {
	err = s.getToken(token, "internal/user/").Into(&user)
	return user, err
}

func (s *Session) UserByTokenFull(token string) (user UserFull, err error) {
	err = s.getToken(token, "internal/user/full").Into(&user)
	return user, err
}

func (s *Session) ListAllUsers(token string) (users []User, err error) {
	err = s.getToken(token, "internal/user/all").Into(&users)
	return users, err
}

func (s *Session) ListAllContactUsers() (users []User, err error) {
	err = s.get("public/contacts").Into(&users)
	return users, err
}

func (s *Session) AddUser(user UserFull, token string) (users UserFull, err error) {
	user1, err := json.Marshal(user)
	if err != nil {
		return UserFull{}, err
	}
	err = s.putToken(token, "internal/user/admin", *bytes.NewBuffer(user1)).Into(&users)
	return users, err
}

func (s *Session) EditUser(user UserFull, token string) (users UserFull, err error) {
	user1, err := json.Marshal(user)
	if err != nil {
		return UserFull{}, err
	}
	err = s.patchToken(token, "internal/user/admin", *bytes.NewBuffer(user1)).Into(&users)
	return users, err
}

type Message struct {
	Message string `json:"message"`
}

func (s *Session) DeleteUser(email, token string) (err error) {
	var message Message
	json1, err := s.deletefToken(token, "internal/user/admin/%s", email).JSON()
	fmt.Println(json1)
	fmt.Println(message)
	return err
}
