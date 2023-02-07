package users

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	// UserRepo defines all user interactions
	UserRepo interface {
		GetUserByEmail(email string) (User, error)
		GetUserById(id uint64) (User, error)
		GetUserByEmailFull(email string) (UserFull, error)
		GetUserByIdFull(id uint64) (UserFull, error)
		ListAllUsers() ([]User, error)
		ListContactUsers() ([]User, error)
		ListTeamManagersUsers(teamId uint64) ([]User, error)
		AddUser(user *UserFull) error
		EditUser(user *UserFull) error
		DeleteUser(email string) error
	}

	// Store contains our dependency
	Store struct {
		scope *gocb.Scope
	}

	//User represents a user object to be used when not all data is required
	User struct {
		Id    uint64 `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone string `json:"phone,omitempty"`
		Team  uint64 `json:"team,omitempty"`
		Role  string `json:"role"`
		Image string `json:"image,omitempty"`
	}
	// UserFull represents a user and all columns
	UserFull struct {
		User
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

// NewStore creates a new store
func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

// GetUserByEmail returns basic user information to be used for other services.
func (m *Store) GetUserByEmail(email string) (u User, err error) {
	result, err := m.scope.Query("SELECT `id` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                true,
		PositionalParameters: []interface{}{email},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return User{}, fmt.Errorf("failed to get user in get user by email: %w", err)
		}
	}
	if result.Next() {
		var d UserFull
		err := result.Row(&d)
		if err != nil {
			return User{}, err
		}
		return m.GetUserById(d.Id)
	} else {
		return User{}, fmt.Errorf("failed to get user by email: %s", email)
	}
}

// GetUserByEmail returns basic user information to be used for other services.
func (m *Store) GetUserById(id uint64) (u User, err error) {
	result, err := m.scope.Collection("users").Get("user:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return User{}, fmt.Errorf("user doesn't exist: %d", id)
		} else {
			return User{}, fmt.Errorf("failed to get user: %w", err)
		}
	}

	err = result.Content(&u)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return u, err
}

// GetUserFullByEmail will return all user information to be used for profile and management.
func (m *Store) GetUserFullByEmail(email string) (u UserFull, err error) {
	result, err := m.scope.Query("SELECT `id` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                true,
		PositionalParameters: []interface{}{email},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return UserFull{}, fmt.Errorf("failed to get user in get user full by email: %w", err)
		}
	}
	if result.Next() {
		var d UserFull
		err := result.Row(&d)
		if err != nil {
			return UserFull{}, err
		}
		return m.GetUserFullById(d.Id)
	} else {
		return UserFull{}, fmt.Errorf("failed to get user full by email: %s", email)
	}
}

// GetUserFullByEmail will return all user information to be used for profile and management.
func (m *Store) GetUserFullById(id uint64) (u UserFull, err error) {
	result, err := m.scope.Collection("users").Get("user:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return UserFull{}, fmt.Errorf("user doesn't exist: %d", id)
		} else {
			return UserFull{}, fmt.Errorf("failed to get user full: %w", err)
		}
	}

	err = result.Content(&u)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get user full: %w", err)
	}
	return u, err
}

// ListAllUsers returns all users
// It doesn't return the full User object
func (m *Store) ListAllUsers() (u []User, err error) {
	query, err := m.scope.Query("SELECT `id`, `email`, `name`, `phone`, `team`, `role`, `image` FROM users ORDER BY `id` ", &gocb.QueryOptions{})
	if err != nil {
		return []User{}, fmt.Errorf("failed to get all users: %w", err)
	}
	for query.Next() {
		var result User
		err := query.Row(&result)
		if err != nil {
			return []User{}, fmt.Errorf("failed to get all users: %w", err)
		}
		u = append(u, result)
	}

	if err := query.Err(); err != nil {
		return []User{}, fmt.Errorf("failed to get all users: %w", err)
	}
	return u, err
}

func (m *Store) ListContactUsers() (u []User, err error) {
	query, err := m.scope.Query("SELECT `email`, `name`, `role` FROM users WHERE `role`='ProgrammeEditor' OR `role`='LeagueSecretary' OR `role`='Treasurer' OR `role`='SafeguardingOfficer' OR `role`='ClubSecretary' OR `role`='Chairperson' ORDER BY CASE `role` WHEN 'ProgrammeEditor' THEN 1 WHEN 'LeagueSecretary' THEN 2 WHEN 'Treasurer' THEN 3 WHEN 'SafeguardingOfficer' THEN 4 WHEN 'ClubSecretary' THEN 5 WHEN 'Chairperson' THEN 6 ELSE 0 END DESC", &gocb.QueryOptions{})
	if err != nil {
		return []User{}, fmt.Errorf("failed to get contact users: %w", err)
	}
	for query.Next() {
		var result User
		err := query.Row(&result)
		if err != nil {
			return []User{}, fmt.Errorf("failed to get contact users: %w", err)
		}
		u = append(u, result)
	}

	if err := query.Err(); err != nil {
		return []User{}, fmt.Errorf("failed to get contact users: %w", err)
	}
	return u, err
}

func (m *Store) ListTeamManagersUsers(teamId uint64) (u []User, err error) {
	query, err := m.scope.Query("SELECT `email`, `name`, `role` FROM users WHERE `role`='Manager' AND `team` = $1", &gocb.QueryOptions{
		PositionalParameters: []interface{}{teamId},
		Adhoc:                true,
	})
	if err != nil {
		return []User{}, fmt.Errorf("failed to get team managers users: %w", err)
	}
	for query.Next() {
		var result User
		err := query.Row(&result)
		if err != nil {
			return []User{}, fmt.Errorf("failed to get team managers users: %w", err)
		}
		u = append(u, result)
	}

	if err := query.Err(); err != nil {
		return []User{}, fmt.Errorf("failed to get team managers users: %w", err)
	}
	return u, err
}

func (m *Store) AddUser(u *UserFull) error {
	result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{u.Email},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get user in add user: %w", err)
		}
	}
	for result.Next() {
		return fmt.Errorf("email already exists")
	}
	_, err = m.scope.Collection("users").Insert("user:"+strconv.FormatUint(u.Id, 10), u, &gocb.InsertOptions{})
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return err
}

func (m *Store) EditUser(u *UserFull) error {
	result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{u.Email},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get user in edit user: %w", err)
		}
	}
	if result.Next() {
		_, err := m.scope.Collection("users").Upsert("user:"+u.Email, u, &gocb.UpsertOptions{})
		if err != nil {
			return fmt.Errorf("failed to edit user: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("email doesn't exists")
	}
}

func (m *Store) DeleteUserFromEmail(email string) error {
	result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{email},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get user in delete user: %w", err)
		}
	}
	if result.Next() {
		var u User
		err = result.Row(&u)
		if err != nil {
			return fmt.Errorf("failed to delete user from email: %w", err)
		}
		_, err := m.scope.Collection("users").Remove("user:"+strconv.FormatUint(u.Id, 10), &gocb.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("email doesn't exists")
	}
}

func (m *Store) DeleteUserFromId(id uint64) error {
	result, err := m.scope.Query("SELECT `email` FROM users WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get user in delete user: %w", err)
		}
	}
	if result.Next() {
		_, err := m.scope.Collection("users").Remove("user:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
