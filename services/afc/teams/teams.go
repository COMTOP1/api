package teams

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	TeamRepo interface {
		GetTeamById(id uint64) (Team, error)
		ListAllTeams() ([]Team, error)
		ListActiveTeas() ([]Team, error)
		AddTeam(t *Team) error
		EditTeam(t *Team) error
		DeleteTeam(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Team struct {
		Id            uint64 `json:"id"`
		Name          string `json:"name"`
		League        string `json:"league,omitempty"`
		Division      string `json:"division,omitempty"`
		LeagueTable   string `json:"league_table,omitempty"`
		Fixtures      string `json:"fixtures,omitempty"`
		Coach         string `json:"coach,omitempty"`
		Physio        string `json:"physio,omitempty"`
		ImageFileName string `json:"image,omitempty"`
		Active        bool   `json:"active"`
		Youth         bool   `json:"youth"`
		Ages          int    `json:"ages"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetTeamById(id uint64) (t Team, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("teams").Get("team:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Team{}, fmt.Errorf("team doesn't exist: %s", id)
		} else {
			return Team{}, fmt.Errorf("failed to get team: %w", err)
		}
	}

	err = result.Content(&t)
	if err != nil {
		return Team{}, fmt.Errorf("failed to get team: %w", err)
	}
	return t, err
}

func (m *Store) ListAllTeams() (t []Team, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `league`, `division`, `league_table`, `fixtures`, `coach`, `physio`, `image`, `active`, `youth`, `ages` FROM teams", &gocb.QueryOptions{})
	if err != nil {
		return []Team{}, fmt.Errorf("failed to get all teams: %w", err)
	}
	for query.Next() {
		var result Team
		err = query.Row(&result)
		if err != nil {
			return []Team{}, fmt.Errorf("failed to get all teams: %w", err)
		}
		t = append(t, result)
	}

	if err := query.Err(); err != nil {
		return []Team{}, fmt.Errorf("failed to get all teams: %w", err)
	}
	return t, err
}

func (m *Store) ListActiveTeams() (t []Team, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `league`, `division`, `league_table`, `fixtures`, `coach`, `physio`, `image`, `active`, `youth`, `ages` FROM teams WHERE `active` = true", &gocb.QueryOptions{})
	if err != nil {
		return []Team{}, fmt.Errorf("failed to get active teams: %w", err)
	}
	for query.Next() {
		var result Team
		err = query.Row(&result)
		if err != nil {
			return []Team{}, fmt.Errorf("failed to get active teams: %w", err)
		}
		t = append(t, result)
	}

	if err := query.Err(); err != nil {
		return []Team{}, fmt.Errorf("failed to get active teams: %w", err)
	}
	return t, err
}

func (m *Store) AddTeam(w *Team) error {
	result, err := m.scope.Query("SELECT `id` FROM teams WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get team in add team: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("teams").Insert("team:"+strconv.FormatUint(w.Id, 10), w, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add team: %w", err)
	}
	return err
}

func (m *Store) EditTeam(w *Team) error {
	result, err := m.scope.Query("SELECT `id` FROM teams WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get team in edit team: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("teams").Upsert("team:"+strconv.FormatUint(w.Id, 10), w, &gocb.UpsertOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit team: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeleteTeam(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM teams WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get team in delete team: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("teams").Remove("team:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit team: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
