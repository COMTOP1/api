package users

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
    "strings"
)

type (
	// UserRepo defines all user interactions
    UserRepo interface {
        GetUser(email string) (User, error)
        GetUserFull(email string) (UserFull, error)
        ListAllUsers() ([]User, error)
        ListContactUsers() ([]User, error)
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

// GetUser returns basic user information to be used for other services.
func (m *Store) GetUser(email string) (u User, err error) {
    result, err := m.scope.Collection("users").Get("user:"+email, nil)
    if err != nil {
        if strings.Contains(err.Error(), "document not found") {
            return User{}, fmt.Errorf("user doesn't exist: %s", email)
        } else {
            return User{}, fmt.Errorf("failed to get user: %w", err)
        }
    }

    err = result.Content(&u)
    if err != nil {
        return User{}, fmt.Errorf("failed to get user: %w", err)
    }
    return u, nil
}

// GetUserFull will return all user information to be used for profile and management.
func (m *Store) GetUserFull(email string) (u UserFull, err error) {
	result, err := m.scope.Collection("users").Get("user:"+email, nil)
	if err != nil {
        if strings.Contains(err.Error(), "document not found") {
            return UserFull{}, fmt.Errorf("user doesn't exist: %s", email)
        } else {
            return UserFull{}, fmt.Errorf("failed to get user: %w", err)
        }
	}

    err = result.Content(&u)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get user: %w", err)
	}
    fmt.Println(u)
	return u, nil
}

// ListAllUsers returns all users
// It doesn't return the full User object
func (m *Store) ListAllUsers() (u []User, err error) {
	query, err := m.scope.Query("SELECT `email`, `name`, `phone`, `team`, `role`, `image` FROM users ", &gocb.QueryOptions{})
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var result User
		err := query.Row(&result)
		if err != nil {
			return []User{}, fmt.Errorf("failed to get users: %w", err)
		}
		fmt.Println(result)
		u = append(u, result)
	}

	if err := query.Err(); err != nil {
		return []User{}, fmt.Errorf("failed to get users: %w", err)
	}
	return u, nil
}

func (m *Store) ListContactUsers() (u []User, err error) {
	query, err := m.scope.Query("SELECT `email`, `name`, `role` FROM users WHERE `role`='ProgrammeEditor' OR `role`='LeagueSecretary' OR `role`='Treasurer' OR `role`='SafeguardingOfficer' OR `role`='ClubSecretary' OR `role`='Chairperson' ORDER BY CASE `role` WHEN 'ProgrammeEditor' THEN 1 WHEN 'LeagueSecretary' THEN 2 WHEN 'Treasurer' THEN 3 WHEN 'SafeguardingOfficer' THEN 4 WHEN 'ClubSecretary' THEN 5 WHEN 'Chairperson' THEN 6 ELSE 0 END DESC", &gocb.QueryOptions{})
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var result User
		err := query.Row(&result)
		if err != nil {
			return []User{}, fmt.Errorf("failed to get all users: %w", err)
		}
		fmt.Println(result)
		u = append(u, result)
	}

	if err := query.Err(); err != nil {
		return []User{}, fmt.Errorf("failed to get all users: %w", err)
	}
	return u, nil
}

func (m *Store) AddUser(u *UserFull) error {
    fmt.Println(u)
	result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{u.Email},
	})
	if err != nil {
        if !strings.Contains(err.Error(), "document not found") {
            return fmt.Errorf("failed to get user: %w", err)
        }
	}
	for result.Next() {
		return fmt.Errorf("email already exists")
	}
	mut, err := m.scope.Collection("users").Insert("user:"+u.Email, u, nil)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	fmt.Println(mut)
	return nil
}

func (m *Store) EditUser(u *UserFull) error {
    result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
        Adhoc:                false,
        PositionalParameters: []interface{}{u.Email},
        })
    if err != nil {
        if !strings.Contains(err.Error(), "document not found") {
            return fmt.Errorf("failed to get user: %w", err)
        }
    }
    if result.Next() {
        mut, err := m.scope.Collection("users").Upsert("user:"+u.Email, u, nil)
        fmt.Println(mut)
        if err != nil {
            return err
        }
        return nil
    } else {
        return fmt.Errorf("email doesn't exists")
    }
}

func (m *Store) DeleteUser(email string) error {
    result, err := m.scope.Query("SELECT `email` FROM users WHERE `email` = $1", &gocb.QueryOptions{
        Adhoc:                false,
        PositionalParameters: []interface{}{email},
        })
    if err != nil {
        if !strings.Contains(err.Error(), "document not found") {
            return fmt.Errorf("failed to get user: %w", err)
        }
    }
    if result.Next() {
        mut, err := m.scope.Collection("users").Remove("user:"+email, nil)
        fmt.Println(mut)
        if err != nil {
            return err
        }
        return nil
    } else {
        return fmt.Errorf("email doesn't exists")
    }
}

/*{
"email": "liam.burnand@bswdi.co.uk",
"hash": "[110, 251, 168, 191, 159, 99, 214, 88, 249, 152, 116, 228, 221, 182, 124, 166, 16, 78, 228, 58, 30, 87, 190, 70, 150, 228, 194, 209, 64, 125, 141, 212, 75, 253, 1, 211, 223, 55, 172, 244, 27, 151, 157, 168, 100, 34, 167, 24, 2, 73, 165, 46, 204, 112, 108, 55, 101, 220, 7, 38, 100, 67, 2, 31]",
"name": "Liam Burnand",
"phone": 447426534286,
"role": "Webmaster",
"salt": "[86, 210, 152, 163, 212, 214, 24, 121, 105, 222, 115, 42, 86, 255, 82, 108, 191, 216, 133, 147, 11, 12, 216, 191, 57, 112, 86, 124, 34, 104, 175, 118, 30, 44, 164, 153, 57, 50, 254, 36, 168, 32, 58, 131, 139, 133, 7, 11, 6, 242, 169, 168, 75, 241, 250, 48, 44, 76, 214, 107, 42, 132, 107, 191]",
"updated_on": "2022-08-09 13:17:46"
}*/
