package documents

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	DocumentRepo interface {
		GetDocumentById(id uint64) (Document, error)
		ListAllDocuments() ([]Document, error)
		AddDocument(p *Document) (Document, error)
		EditDocument(p *Document) (Document, error)
		DeleteDocument(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Document struct {
		Id       uint64 `json:"id"`
		Name     string `json:"name"`
		FileName string `json:"file_name"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetDocumentById(id uint64) (d Document, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("documents").Get("document:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Document{}, fmt.Errorf("document doesn't exist: %d", id)
		} else {
			return Document{}, fmt.Errorf("failed to get document: %w", err)
		}
	}

	err = result.Content(&d)
	if err != nil {
		return Document{}, fmt.Errorf("failed to get document: %w", err)
	}
	return d, err
}

func (m *Store) ListAllDocuments() (d []Document, err error) {
	query, err := m.scope.Query("SELECT `id`, `name`, `file_name` FROM documents", &gocb.QueryOptions{})
	if err != nil {
		return []Document{}, fmt.Errorf("failed to get all documents: %w", err)
	}
	for query.Next() {
		var result Document
		err = query.Row(&result)
		if err != nil {
			return []Document{}, fmt.Errorf("failed to get all documents: %w", err)
		}
		d = append(d, result)
	}

	if err := query.Err(); err != nil {
		return []Document{}, fmt.Errorf("failed to get all documents: %w", err)
	}
	return d, err
}

func (m *Store) AddDocument(document *Document) error {
	result, err := m.scope.Query("SELECT `id` FROM documents WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{document.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get document in add document: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("documents").Insert("document:"+strconv.FormatUint(document.Id, 10), document, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add document: %w", err)
	}
	return err
}

func (m *Store) DeleteDocument(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM documents WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get document in delete document: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("documents").Remove("document:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit document: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
