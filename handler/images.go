package handler

import (
	"bytes"
	"encoding/json"
)

type Image struct {
	Id       uint64 `json:"id"`
	FileName string `json:"file_name"`
	Caption  string `json:"caption"`
}

func (s *Session) GetImageById(id uint64) (i Image, err error) {
	err = s.getf("public/image/%d", id).Into(&i)
	if err != nil {
		return Image{}, err
	}
	return i, err
}

func (s *Session) ListAllImages() (i []Image, err error) {
	err = s.get("public/images").Into(&i)
	if err != nil {
		return []Image{}, err
	}
	return i, err
}

func (s *Session) AddImage(image Image, token string) (i Image, err error) {
	i1, err := json.Marshal(image)
	if err != nil {
		return Image{}, err
	}
	err = s.putToken(token, "internal/image", *bytes.NewBuffer(i1)).Into(&i)
	if err != nil {
		return Image{}, err
	}
	return i, err
}

func (s *Session) DeleteImage(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/image/%d", id).JSON()
	return err
}
