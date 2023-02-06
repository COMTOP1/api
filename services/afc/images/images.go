package images

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"strings"
)

type (
	ImageRepo interface {
		GetImageById(id uint64) (Image, error)
		ListAllImages() ([]Image, error)
		AddImage(image *Image) (Image, error)
		DeleteImage(id uint64) error
	}

	Store struct {
		scope *gocb.Scope
	}

	Image struct {
		Id       uint64 `json:"id"`
		FileName string `json:"file_name"`
		Caption  string `json:"caption"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (m *Store) GetImageById(id uint64) (i Image, err error) {
	m.scope.BucketName()
	result, err := m.scope.Collection("images").Get("image:"+strconv.FormatUint(id, 10), &gocb.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return Image{}, fmt.Errorf("image doesn't exist: %s", id)
		} else {
			return Image{}, fmt.Errorf("failed to get image: %w", err)
		}
	}

	err = result.Content(&i)
	if err != nil {
		return Image{}, fmt.Errorf("failed to get image: %w", err)
	}
	return i, err
}

func (m *Store) ListAllImages() (i []Image, err error) {
	query, err := m.scope.Query("SELECT `id`, `file_name`, `caption` FROM images", &gocb.QueryOptions{})
	if err != nil {
		return []Image{}, fmt.Errorf("failed to get all images: %w", err)
	}
	for query.Next() {
		var result Image
		err = query.Row(&result)
		if err != nil {
			return []Image{}, fmt.Errorf("failed to get all images: %w", err)
		}
		i = append(i, result)
	}

	if err := query.Err(); err != nil {
		return []Image{}, fmt.Errorf("failed to get all images: %w", err)
	}
	return i, err
}

func (m *Store) AddImage(i *Image) error {
	result, err := m.scope.Query("SELECT `id` FROM images WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{i.Id},
	})
	if err != nil {
		return fmt.Errorf("failed to get image in add image: %w", err)
	}
	for result.Next() {
		return fmt.Errorf("id already exists")
	}
	mut, err := m.scope.Collection("images").Insert("image:"+strconv.FormatUint(i.Id, 10), i, &gocb.InsertOptions{})
	fmt.Println(mut)
	if err != nil {
		return fmt.Errorf("failed to add image: %w", err)
	}
	return err
}

func (m *Store) DeleteImage(id uint64) error {
	result, err := m.scope.Query("SELECT `id` FROM images WHERE `id` = $1", &gocb.QueryOptions{
		Adhoc:                false,
		PositionalParameters: []interface{}{id},
	})
	if err != nil {
		return fmt.Errorf("failed to get image in delete image: %w", err)
	}
	if result.Next() {
		mut, err := m.scope.Collection("images").Remove("image:"+strconv.FormatUint(id, 10), &gocb.RemoveOptions{})
		fmt.Println(mut)
		if err != nil {
			return fmt.Errorf("failed to edit image: %w", err)
		}
		return err
	} else {
		return fmt.Errorf("id doesn't exists")
	}
}
