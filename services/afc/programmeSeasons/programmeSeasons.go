package programmeSeasons

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	Repo interface {
		GetProgrammeSeasonById(id uint64) (ProgrammeSeason, error)
		ListAllProgrammesSeasons() ([]ProgrammeSeason, error)
		AddProgrammeSeason(p *ProgrammeSeason) (ProgrammeSeason, error)
		EditProgrammeSeason(p *ProgrammeSeason) (ProgrammeSeason, error)
		DeleteProgrammeSeason(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	ProgrammeSeason struct {
		Id     uint64 `json:"id"`
		Season string `json:"season"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetProgrammeSeasonById(id uint64) (p ProgrammeSeason, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("programme_seasons").Get("programme_season:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return ProgrammeSeason{}, fmt.Errorf("programme season doesn't exist: %d", id)
		} else {
			return ProgrammeSeason{}, fmt.Errorf("failed to get programme season: %w", err)
		}
	}

	err = result.Content(&p)
	if err != nil {
		return ProgrammeSeason{}, fmt.Errorf("failed to get programme season: %w", err)
	}
	return p, err
}

func (m *Store) ListAllProgrammeSeasons() (p []ProgrammeSeason, err error) {
	query, err := m.scope.Query("SELECT `id`, `season` FROM programme_seasons", &gocb.QueryOptions{})
	if err != nil {
		return []ProgrammeSeason{}, fmt.Errorf("failed to get all programme seasons: %w", err)
	}
	for query.Next() {
		var result ProgrammeSeason
		err = query.Row(&result)
		if err != nil {
			return []ProgrammeSeason{}, fmt.Errorf("failed to get all programme seasons: %w", err)
		}
		p = append(p, result)
	}

	if err := query.Err(); err != nil {
		return []ProgrammeSeason{}, fmt.Errorf("failed to get all programme seasons: %w", err)
	}
	return p, err
}

func (m *Store) AddProgrammeSeason(p *ProgrammeSeason) error {
	result, err := m.scope.Query("SELECT `id` FROM programme_seasons WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get programme season in add programme season: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("programme_seasons").Insert("programme_season:"+strconv.FormatUint(p.Id, 10), p, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add programme season: %w", err)
	}
	return err
}

func (m *Store) EditProgrammeSeason(p *ProgrammeSeason) error {
	result, err := m.scope.Query("SELECT `id` FROM programme_seasons WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get programme season in edit programme season: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("programme_seasons").Upsert("programme_season:"+strconv.FormatUint(p.Id, 10), p, &gocb.UpsertOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit programme season: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeleteProgrammeSeason(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM programme_seasons WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get programme season in delete programme season: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("programme_seasons").Remove("programme_season:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit programme season: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
