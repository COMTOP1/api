package users

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
)

var _ UserRepo = &Store{}

// GetUserFull will return all user information to be used for profile and management.
func (m *Store) GetUserFull(email string) (u UserFull, err error) {
	doc, err := m.scope.Collection("users").Get("u:"+email, nil)

	err = doc.Content(&u)
	if err != nil {
		return UserFull{}, err
	}
	return u, nil
}

// GetUser returns basic user information to be used for other services.
func (m *Store) GetUser(email string) (u User, err error) {
	result, err := m.scope.Collection("users").Get("user:"+email, nil)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}

	err = result.Content(&u)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}

// ListAllUsers returns all users
// It doesn't return the full User object
func (m *Store) ListAllUsers() (u []User, err error) {
	query, err := m.scope.Query("SELECT * FROM users", &gocb.QueryOptions{})
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
