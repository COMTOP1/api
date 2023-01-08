package whatsOn

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
    "strconv"
)

type (
	WhatsOnRepo interface {
		GetWhatsOn(id uint64) (WhatsOn, error)
        GetWhatsOnLatest() (WhatsOn, error)
        ListAllWhatsOn() ([]WhatsOn, error)
        AddWhatsOn(w *WhatsOn) error
        EditWhatsOn(w *WhatsOn) error
        DeleteWhatsOn(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	WhatsOn struct {
        Id uint64 `json:"id"`
        Title string `json:"title"`
        ImageFileName string `json:"image,omitempty"`
        Content string `json:"content"`
        Date uint64 `json:"date"`
        DateOfEvent uint64 `json:"dateOfEvent"`
        Delete bool `json:"delete,omitempty"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetWhatsOn(id uint64) (w WhatsOn, err error) {
    result, err := m.scope.Collection("whatsOn").Get("whatsOn:"+strconv.FormatUint(id, 10), nil)
    if err != nil {
        return WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
    }

    err = result.Content(&w)
    if err != nil {
        return WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
    }
    return w, err
}

func (m *Store) GetWhatsOnLatest() (w WhatsOn, err error) {
    query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `dateOfEvent` FROM whatsOn ORDER BY `dateOfEvent` DESC LIMIT 1", &gocb.QueryOptions{})
    if err != nil {
        return WhatsOn{}, err
    }
    if query.Next() {
        err = query.Row(&w)
        if err != nil {
            return WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
        }
    }
    return w, err
}

func (m *Store) ListAllWhatsOn() (w []WhatsOn, err error) {
    query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `dateOfEvent` FROM whatsOn ", &gocb.QueryOptions{})
    if err != nil {
        return nil, err
    }
    for query.Next() {
        var result WhatsOn
        err = query.Row(&result)
        if err != nil {
            return []WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
        }
        w = append(w, result)
    }

    if err := query.Err(); err != nil {
        return []WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
    }
    return w, nil
}

func (m *Store) AddWhatsOn(w *WhatsOn) error {
    result, err := m.scope.Query("SELECT `id` FROM whatsOn WHERE `id` = $1", &gocb.QueryOptions{
        Adhoc:                false,
        PositionalParameters: []interface{}{w.Id},
        })
    if err != nil {
        return err
    }
    for result.Next() {
        return fmt.Errorf("id already exists")
    }
    mut, err := m.scope.Collection("whatsOn").Insert("whatsOn:"+strconv.FormatUint(w.Id, 10), w, nil)
    fmt.Println(mut)
    if err != nil {
        return fmt.Errorf("failed to add whatsOn: %w", err)
    }
    return nil
}

func (m *Store) EditWhatsOn(w *WhatsOn) error {
    result, err := m.scope.Query("SELECT `id` FROM whatsOn WHERE `id` = $1", &gocb.QueryOptions{
        Adhoc:                false,
        PositionalParameters: []interface{}{w.Id},
        })
    if err != nil {
        return err
    }
    if result.Next() {
        mut, err := m.scope.Collection("whatsOn").Upsert("whatsOn:"+strconv.FormatUint(w.Id, 10), w, nil)
        fmt.Println(mut)
        if err != nil {
            return err
        }
        return nil
    } else {
        return fmt.Errorf("id doesn't exists")
    }
}

func (m *Store) DeleteWhatsOn(id uint64) error {
    result, err := m.scope.Query("SELECT `id` FROM whatsOn WHERE `id` = $1", &gocb.QueryOptions{
        Adhoc:                false,
        PositionalParameters: []interface{}{id},
        })
    if err != nil {
        return err
    }
    if result.Next() {
        mut, err := m.scope.Collection("whatsOn").Remove("whatsOn:"+strconv.FormatUint(id, 10), nil)
        fmt.Println(mut)
        if err != nil {
            return err
        }
        return nil
    } else {
        return fmt.Errorf("id doesn't exists")
    }
}