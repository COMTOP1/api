package programmes

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	ProgrammeRepo interface {
		GetProgrammeById(id uint64) (Programme, error)
		ListAllProgrammes() ([]Programme, error)
		AddProgramme(p *Programme) (Programme, error)
		EditProgramme(p *Programme) (Programme, error)
		DeleteProgramme(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Programme struct {
		Id                uint64 `json:"id"`
		Name              string `json:"name"`
		FileName          string `json:"file_name"`
		DateOfProgramme   int64  `json:"date_of_programme"`
		ProgrammeSeasonId uint64 `json:"programme_season_id"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetProgrammeById(id uint64) (p Programme, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("programmes").Get("programme:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Programme{}, fmt.Errorf("programme doesn't exist: %s", id)
		} else {
			return Programme{}, fmt.Errorf("failed to get programme: %w", err)
		}
	}

	err = result.Content(&p)
	if err != nil {
		return Programme{}, fmt.Errorf("failed to get programme: %w", err)
	}
	return p, err
}

func (m *Store) ListAllProgrammes() (p []Programme, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `file_name`, `date_of_programme`, `programme_season_id` FROM programmes", &gocb.QueryOptions{})
	if err != nil {
		return []Programme{}, fmt.Errorf("failed to get all programmes: %w", err)
	}
	for query.Next() {
		var result Programme
		err = query.Row(&result)
		if err != nil {
			return []Programme{}, fmt.Errorf("failed to get all programmes: %w", err)
		}
		p = append(p, result)
	}

	if err := query.Err(); err != nil {
		return []Programme{}, fmt.Errorf("failed to get all programmes: %w", err)
	}
	return p, err
}

func (m *Store) AddProgramme(p *Programme) error {
	result, err := m.scope.Query("SELECT `id` FROM programmes WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get programme in add programme: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("programmes").Insert("programme:"+strconv.FormatUint(p.Id, 10), p, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add programme: %w", err)
	}
	return err
}

func (m *Store) DeleteProgramme(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM programmes WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get programme in delete programme: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("programmes").Remove("programme:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit programme: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
