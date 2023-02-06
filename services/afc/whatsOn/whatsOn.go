package whatsOn

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
	"time"
)

type (
	WhatsOnRepo interface {
		GetWhatsOnById(id uint64) (WhatsOn, error)
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
		Id            uint64 `json:"id"`
		Title         string `json:"title"`
		ImageFileName string `json:"image,omitempty"`
		Content       string `json:"content"`
		Date          int64  `json:"date"`
		DateOfEvent   int64  `json:"date_of_event"`
		Delete        bool   `json:"delete,omitempty"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetWhatsOnById(id uint64) (w WhatsOn, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("whats_on").Get("whats_on:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return WhatsOn{}, fmt.Errorf("whatsOn doesn't exist: %d", id)
		} else {
			return WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
		}
	}

	err = result.Content(&w)
	if err != nil {
		return WhatsOn{}, fmt.Errorf("failed to get whatsOn: %w", err)
	}
	return w, err
}

func (m *Store) GetWhatsOnLatest() (w WhatsOn, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `date_of_event` FROM whats_on ORDER BY `date_of_event` DESC LIMIT 1", &gocb.QueryOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return WhatsOn{}, fmt.Errorf("no whatsOn exist")
		} else {
			return WhatsOn{}, fmt.Errorf("failed to get whatsOn latest: %w", err)
		}
	}
	if query.Next() {
		err = query.Row(&w)
		if err != nil {
			return WhatsOn{}, fmt.Errorf("failed to get whatsOn latest: %w", err)
		}
	}
	return w, err
}

func (m *Store) ListAllWhatsOnEventPast() (w []WhatsOn, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `date_of_event` FROM whats_on WHERE date_of_event < ? ORDER BY date_of_event DESC", &gocb.QueryOptions{
		PositionalParameters: []interface{}{time.Now().Unix()},
	})
	if err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event past: %w", err)
	}
	for query.Next() {
		var result WhatsOn
		err = query.Row(&result)
		if err != nil {
			return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event past: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event past: %w", err)
	}
	return w, err
}

func (m *Store) ListAllWhatsOn() (w []WhatsOn, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `date_of_event` FROM whats_on", &gocb.QueryOptions{})
	if err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn: %w", err)
	}
	for query.Next() {
		var result WhatsOn
		err = query.Row(&result)
		if err != nil {
			return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn: %w", err)
	}
	return w, err
}

func (m *Store) ListAllWhatsOnEventFuture() (w []WhatsOn, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `date_of_event` FROM whats_on WHERE date_of_event >= ? ORDER BY date_of_event", &gocb.QueryOptions{
		PositionalParameters: []interface{}{time.Now().Unix()},
	})
	if err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event future: %w", err)
	}
	for query.Next() {
		var result WhatsOn
		err = query.Row(&result)
		if err != nil {
			return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event future: %w", err)
		}
		w = append(w, result)
	}

	if err := query.Err(); err != nil {
		return []WhatsOn{}, fmt.Errorf("failed to get all whatsOn event future: %w", err)
	}
	return w, err
}

func (m *Store) AddWhatsOn(w *WhatsOn) error {
	result, err := m.scope.Query("SELECT `id` FROM whats_on WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get whatsOn in add whatsOn: %w", err)
		}
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	_, err = m.scope.Collection("whats_on").Insert("whats_on:"+strconv.FormatUint(w.Id, 10), w, &gocb.InsertOptions{})
	if err != nil {
		return fmt.Errorf("failed to add whatsOn: %w", err)
	}
	return nil
}

func (m *Store) EditWhatsOn(w *WhatsOn) error {
	result, err := m.scope.Query("SELECT `id` FROM whats_on WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{w.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get whatsOn in edit whatsOn: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("whats_on").Upsert("whats_on:"+strconv.FormatUint(w.Id, 10), w, &gocb.UpsertOptions{})
		if err != nil {
			return fmt.Errorf("failed to edit whatsOn: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeleteWhatsOn(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM whats_on WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get whatsOn in delete whatsOn: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("whats_on").Remove("whats_on:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete whatsOn: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
