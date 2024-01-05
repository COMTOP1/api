package affiliations

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	AffiliationRepo interface {
		GetAffiliation(id uint64) (Affiliation, error)
		ListAffiliation() ([]Affiliation, error)
		AddAffiliation(affiliation Affiliation) error
		EditAffiliation(affiliation Affiliation) error
		DeleteAffiliation(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Affiliation struct {
		Id       uint64 `json:"id"`
		Name     string `json:"name"`
		Website  string `json:"website"`
		FileName string `json:"file_name"`
	}
)

// NewStore creates a new store
func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetAffiliationById(id uint64) (p Affiliation, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("affiliations").Get("affiliation:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Affiliation{}, fmt.Errorf("affiliation doesn't exist: %d", id)
		} else {
			return Affiliation{}, fmt.Errorf("failed to get affiliation: %w", err)
		}
	}

	err = result.Content(&p)
	if err != nil {
		return Affiliation{}, fmt.Errorf("failed to get affiliation: %w", err)
	}
	return p, err
}

func (m *Store) ListAllAffiliations() (p []Affiliation, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `website`, `file_name` FROM affiliations", &gocb.QueryOptions{})
	if err != nil {
		return []Affiliation{}, fmt.Errorf("failed to get all affiliations: %w", err)
	}
	for query.Next() {
		var result Affiliation
		err = query.Row(&result)
		if err != nil {
			return []Affiliation{}, fmt.Errorf("failed to get all affiliations: %w", err)
		}
		p = append(p, result)
	}

	if err := query.Err(); err != nil {
		return []Affiliation{}, fmt.Errorf("failed to get all affiliations: %w", err)
	}
	return p, err
}

func (m *Store) AddAffiliation(p *Affiliation) error {
	result, err := m.scope.Query("SELECT `id` FROM affiliations WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{p.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get affiliation in add affiliation: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("affiliations").Insert("affiliation:"+strconv.FormatUint(p.Id, 10), p, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add affiliation: %w", err)
	}
	return err
}

func (m *Store) DeleteAffiliation(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM affiliations WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get affiliation in delete affiliation: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("affiliations").Remove("affiliation:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit affiliation: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
