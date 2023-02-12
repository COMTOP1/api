package sponsors

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	SponsorRepo interface {
		GetSponsorById(id uint64) (Sponsor, error)
		ListAllSponsors() ([]Sponsor, error)
		ListAllSponsorsMinimal() ([]Sponsor, error)
		ListAllSponsorsByTeamId(teamId string) ([]Sponsor, error)
		AddSponsor(s *Sponsor) error
		EditSponsor(s *Sponsor) error
		DeleteSponsor(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Sponsor struct {
		Id       uint64 `json:"id"`
		Name     string `json:"name"`
		Website  string `json:"website"`
		FileName string `json:"file_name,omitempty"`
		Purpose  string `json:"purpose,omitempty"`
		Team     string `json:"team,omitempty"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetSponsorById(id uint64) (w Sponsor, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("sponsors").Get("sponsor:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Sponsor{}, fmt.Errorf("sponsor doesn't exist: %d", id)
		} else {
			return Sponsor{}, fmt.Errorf("failed to get sponsor: %w", err)
		}
	}

	err = result.Content(&w)
	if err != nil {
		return Sponsor{}, fmt.Errorf("failed to get sponsor: %w", err)
	}
	return w, err
}

func (m *Store) ListAllSponsors() (w []Sponsor, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `website`, `file_name`, `purpose`, `team` FROM sponsors ", &gocb.QueryOptions{})
	if err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	for query.Next() {
		var result Sponsor
		err = query.Row(&result)
		if err != nil {
			return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	return w, err
}

func (m *Store) ListAllSponsorsMinimal() (w []Sponsor, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `website` FROM sponsors ", &gocb.QueryOptions{})
	if err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	for query.Next() {
		var result Sponsor
		err = query.Row(&result)
		if err != nil {
			return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	return w, err
}

func (m *Store) ListAllSponsorsByTeamId(teamId string) (w []Sponsor, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `website`, `file_name`, `purpose`, `team` FROM sponsors WHERE `team` = $1", &gocb.QueryOptions{
		PositionalParameters: []interface{}{teamId},
		Adhoc:                true,
	})
	if err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	for query.Next() {
		var result Sponsor
		err = query.Row(&result)
		if err != nil {
			return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []Sponsor{}, fmt.Errorf("failed to get all sponsor: %w", err)
	}
	return w, err
}

func (m *Store) AddSponsor(w *Sponsor) error {
	result, err := m.scope.Query("SELECT `id` FROM sponsors WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get sponsor in add sponsor: %w", err)
		}
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	_, err = m.scope.Collection("sponsors").Insert("sponsor:"+strconv.FormatUint(w.Id, 10), w, &gocb.InsertOptions{})
	if err != nil {
		return fmt.Errorf("failed to add sponsor: %w", err)
	}
	return nil
}

func (m *Store) EditSponsor(w *Sponsor) error {
	result, err := m.scope.Query("SELECT `id` FROM sponsors WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get sponsor in edit sponsor: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("sponsors").Upsert("sponsor:"+strconv.FormatUint(w.Id, 10), w, &gocb.UpsertOptions{})
		if err != nil {
			return fmt.Errorf("failed to edit sponsor: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeleteSponsor(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM sponsors WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get sponsor in delete sponsor: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("sponsors").Remove("sponsor:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete sponsor: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
