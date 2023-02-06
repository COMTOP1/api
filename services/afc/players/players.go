package players

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	PlayerRepo interface {
		GetPlayerById(id uint64) (Player, error)
		ListAllPlayers() ([]Player, error)
		ListAllPlayersByTeamId(teamId uint64) ([]Player, error)
		AddPlayer(p *Player) (Player, error)
		EditPlayer(p *Player) (Player, error)
		DeletePlayer(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Player struct {
		Id          uint64 `json:"id"`
		Name        string `json:"name"`
		FileName    string `json:"file_name"`
		DateOfBirth int64  `json:"date_of_birth"`
		Captain     bool   `json:"captain"`
		Team        uint64 `json:"team"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetPlayerById(id uint64) (p Player, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("players").Get("player:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Player{}, fmt.Errorf("player doesn't exist: %s", id)
		} else {
			return Player{}, fmt.Errorf("failed to get player: %w", err)
		}
	}

	err = result.Content(&p)
	if err != nil {
		return Player{}, fmt.Errorf("failed to get player: %w", err)
	}
	return p, err
}

func (m *Store) ListAllPlayers() (p []Player, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `file_name`, `date_of_birth`, `captain`, `team` FROM players", &gocb.QueryOptions{})
	if err != nil {
		return []Player{}, fmt.Errorf("failed to get all players: %w", err)
	}
	for query.Next() {
		var result Player
		err = query.Row(&result)
		if err != nil {
			return []Player{}, fmt.Errorf("failed to get all players: %w", err)
		}
		p = append(p, result)
	}

	if err := query.Err(); err != nil {
		return []Player{}, fmt.Errorf("failed to get all players: %w", err)
	}
	return p, err
}

func (m *Store) ListAllPlayersByTeamId(teamId uint64) (p []Player, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `file_name`, `date_of_birth`, `captain`, `team` FROM players WHERE `team` = $1", &gocb.QueryOptions{
		Adhoc:                true,
		PositionalParameters: []interface{}{teamId},
	})
	if err != nil {
		return []Player{}, fmt.Errorf("failed to get all players: %w", err)
	}
	for query.Next() {
		var result Player
		err = query.Row(&result)
		if err != nil {
			return []Player{}, fmt.Errorf("failed to get all players: %w", err)
		}
		p = append(p, result)
	}

	if err := query.Err(); err != nil {
		return []Player{}, fmt.Errorf("failed to get all players: %w", err)
	}
	return p, err
}

func (m *Store) AddPlayer(p *Player) error {
	result, err := m.scope.Query("SELECT `id` FROM players WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get player in add player: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("players").Insert("player:"+strconv.FormatUint(p.Id, 10), p, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add player: %w", err)
	}
	return err
}

func (m *Store) EditPlayer(p *Player) error {
	result, err := m.scope.Query("SELECT `id` FROM players WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get player in edit player: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("players").Upsert("player:"+strconv.FormatUint(p.Id, 10), p, &gocb.UpsertOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit player: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeletePlayer(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM players WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get player in delete player: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("players").Remove("player:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit player: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
