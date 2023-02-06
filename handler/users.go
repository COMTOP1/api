package handler

import (
	"bytes"
	"encoding/json"
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
		User
		Password  string `json:"password,omitempty"`
		Temp      string `json:"temp,omitempty"`
		Hash      []byte `json:"hash,omitempty"`
		Salt      []byte `json:"salt,omitempty"`
		LastLogin int64  `json:"last_login,omitempty"`
		CreatedAt int64  `json:"created_at,omitempty"`
		CreatedBy User   `json:"created_by,omitempty"`
		UpdatedAt int64  `json:"updated_at,omitempty"`
		UpdatedBy User   `json:"updated_by,omitempty"`
		DeletedAt int64  `json:"deleted_at,omitempty"`
		DeletedBy User   `json:"deleted_by,omitempty"`
	}
)

func (s *Session) GetUserByEmail(email, token string) (u User, err error) {
	err = s.getfToken(token, "internal/user/email/%s", email).Into(&u)
	return u, err
}

func (s *Session) GetUserById(id uint64, token string) (u User, err error) {
	err = s.getfToken(token, "internal/user/id/%d", id).Into(&u)
	return u, err
}

func (s *Session) GetUserByEmailFull(email, token string) (u UserFull, err error) {
	err = s.getfToken(token, "internal/user/email/%s/full", email).Into(&u)
	return u, err
}

func (s *Session) GetUserByIdFull(id uint64, token string) (u UserFull, err error) {
	err = s.getfToken(token, "internal/user/id/%d/full", id).Into(&u)
	return u, err
}

func (s *Session) GetUserByToken(token string) (u User, err error) {
	err = s.getToken(token, "internal/user/").Into(&u)
	return u, err
}

func (s *Session) GetUserByTokenFull(token string) (u UserFull, err error) {
	err = s.getToken(token, "internal/user/full").Into(&u)
	return u, err
}

func (s *Session) ListAllUsers(token string) (u []User, err error) {
	err = s.getToken(token, "internal/user/all").Into(&u)
	return u, err
}

func (s *Session) ListContactUsers() (u []User, err error) {
	err = s.get("public/contacts").Into(&u)
	return u, err
}

func (s *Session) ListTeamManagersUsers(teamId uint64) (u []User, err error) {
	err = s.getf("public/team/managers/%d", teamId).Into(&u)
	return u, err
}

func (s *Session) AddUser(user UserFull, token string) (u UserFull, err error) {
	u1, err := json.Marshal(user)
	if err != nil {
		return UserFull{}, err
	}
	err = s.putToken(token, "internal/user/admin", *bytes.NewBuffer(u1)).Into(&u)
	return u, err
}

func (s *Session) EditUser(user UserFull, token string) (u UserFull, err error) {
	u1, err := json.Marshal(user)
	if err != nil {
		return UserFull{}, err
	}
	err = s.patchToken(token, "internal/user/admin", *bytes.NewBuffer(u1)).Into(&u)
	return u, err
}

func (s *Session) DeleteUserFromEmail(email, token string) (err error) {
	_, err = s.deletefToken(token, "internal/user/admin/email/%s", email).JSON()
	return err
}

func (s *Session) DeleteUserFromId(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/user/admin/id/%d", id).JSON()
	return err
}
