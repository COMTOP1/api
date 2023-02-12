package handler

import (
	"bytes"
	"encoding/json"
)

type News struct {
	Id       uint64 `json:"id"`
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	Content  string `json:"content"`
	Date     int64  `json:"date"`
}

func (s *Session) GetNewsById(id uint64) (n News, err error) {
	err = s.getf("public/news/%d", id).Into(&n)
	if err != nil {
		return News{}, err
	}
	return n, err
}

func (s *Session) GetNewsLatest() (n News, err error) {
	err = s.get("public/news/latest").Into(&n)
	if err != nil {
		return News{}, err
	}
	return n, err
}

func (s *Session) ListAllNews() (n []News, err error) {
	err = s.get("public/news").Into(&n)
	if err != nil {
		return []News{}, err
	}
	return n, err
}

func (s *Session) AddNews(news News, token string) (n News, err error) {
	n1, err := json.Marshal(news)
	if err != nil {
		return News{}, err
	}
	err = s.putToken(token, "internal/news", *bytes.NewBuffer(n1)).Into(&n)
	if err != nil {
		return News{}, err
	}
	return n, err
}

func (s *Session) EditNews(news News, token string) (p News, err error) {
	n1, err := json.Marshal(news)
	if err != nil {
		return News{}, err
	}
	err = s.patchToken(token, "internal/news", *bytes.NewBuffer(n1)).Into(&p)
	if err != nil {
		return News{}, err
	}
	return p, err
}

func (s *Session) DeleteNews(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/news/%d", id).JSON()
	return err
}
