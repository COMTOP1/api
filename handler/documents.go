package handler

import (
	"bytes"
	"encoding/json"
)

type Document struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

func (s *Session) GetDocumentById(id uint64) (d Document, err error) {
	err = s.getf("public/document/%d", id).Into(&d)
	if err != nil {
		return Document{}, err
	}
	return d, err
}

func (s *Session) ListAllDocuments() (d []Document, err error) {
	err = s.get("public/documents").Into(&d)
	if err != nil {
		return []Document{}, err
	}
	return d, err
}

func (s *Session) AddDocument(document Document, token string) (d Document, err error) {
	d1, err := json.Marshal(document)
	if err != nil {
		return Document{}, err
	}
	err = s.putToken(token, "internal/document", *bytes.NewBuffer(d1)).Into(&d)
	if err != nil {
		return Document{}, err
	}
	return d, err
}

func (s *Session) EditDocument(document Document, token string) (d Document, err error) {
	d1, err := json.Marshal(document)
	if err != nil {
		return Document{}, err
	}
	err = s.patchToken(token, "internal/document", *bytes.NewBuffer(d1)).Into(&d)
	if err != nil {
		return Document{}, err
	}
	return d, err
}

func (s *Session) DeleteDocument(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/document/%d", id).JSON()
	return err
}
