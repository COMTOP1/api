package news

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	NewsRepo interface {
	}

	Store struct {
		scope *gocb.Scope
	}

	News struct {
		Id       uint64 `json:"id"`
		Title    string `json:"title"`
		FileName string `json:"file_name"`
		Content  string `json:"content"`
		Date     int64  `json:"date"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetNewsById(id uint64) (n News, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("news").Get("news:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return News{}, fmt.Errorf("news doesn't exist: %d", id)
		} else {
			return News{}, fmt.Errorf("failed to get news: %w", err)
		}
	}

	err = result.Content(&n)
	if err != nil {
		return News{}, fmt.Errorf("failed to get news: %w", err)
	}
	return n, err
}

func (m *Store) GetNewsLatest() (n News, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `image`, `content`, `date`, `date_of_event` FROM news ORDER BY `date_of_event` DESC LIMIT 1", &gocb.QueryOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return News{}, fmt.Errorf("no news exist")
		} else {
			return News{}, fmt.Errorf("failed to get news latest: %w", err)
		}
	}
	if query.Next() {
		err = query.Row(&n)
		if err != nil {
			return News{}, fmt.Errorf("failed to get news latest: %w", err)
		}
	}
	return n, err
}

func (m *Store) ListAllNews() (n []News, err error) {
	query, err := m.scope.Query("SELECT `id`, `title`, `file_name`, `content`, `date` FROM news", &gocb.QueryOptions{})
	if err != nil {
		return []News{}, fmt.Errorf("failed to get all news: %w", err)
	}
	for query.Next() {
		var result News
		err = query.Row(&result)
		if err != nil {
			return []News{}, fmt.Errorf("failed to get all news: %w", err)
		}
		n = append(n, result)
	}

	if err := query.Err(); err != nil {
		return []News{}, fmt.Errorf("failed to get all news: %w", err)
	}
	return n, err
}

func (m *Store) AddNews(n *News) error {
	result, err := m.scope.Query("SELECT `id` FROM news WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{n.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get news in add news: %w", err)
		}
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	_, err = m.scope.Collection("news").Insert("news:"+strconv.FormatUint(n.Id, 10), n, &gocb.InsertOptions{})
	if err != nil {
		return fmt.Errorf("failed to add news: %w", err)
	}
	return nil
}

func (m *Store) EditNews(n *News) error {
	result, err := m.scope.Query("SELECT `id` FROM news WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{n.Id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get news in edit news: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("news").Upsert("news:"+strconv.FormatUint(n.Id, 10), n, &gocb.UpsertOptions{})
		if err != nil {
			return fmt.Errorf("failed to edit news: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}

func (m *Store) DeleteNews(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM news WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "document not found") {
			return fmt.Errorf("failed to get news in delete news: %w", err)
		}
	}
	if result.Next() {
		_, err = m.scope.Collection("news").Remove("news:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete news: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
