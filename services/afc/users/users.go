package users

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

type (
	// UserRepo defines all user interactions
	UserRepo interface {
		GetUser(email string) (User, error)
		GetUserFull(email string) (UserFull, error)
		ListAllUsers() ([]User, error)
		ListContactUsers() ([]User, error)
	}

	// Store contains our dependency
	Store struct {
		scope *gocb.Scope
	}
)

type (
	//User represents a user object to be used when not all data is required
	User struct {
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
		Password  string    `json:"password,omitempty"`
		Temp      string    `json:"temp,omitempty"`
		Hash      []uint8   `json:"hash,omitempty"`
		Salt      []uint8   `json:"salt,omitempty"`
		LastLogin time.Time `json:"last_login,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		CreatedBy User      `json:"created_by,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		UpdatedBy User      `json:"updated_by,omitempty"`
		DeletedAt time.Time `json:"deleted_at,omitempty"`
		DeletedBy User      `json:"deleted_by,omitempty"`
	}
)

// NewStore creates a new store
func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
