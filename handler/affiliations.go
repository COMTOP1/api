package handler

import (
	"bytes"
	"encoding/json"
)

type Affiliation struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Website  string `json:"website"`
	FileName string `json:"file_name"`
}

func (s *Session) GetAffiliationById(id uint64) (a Affiliation, err error) {
	err = s.getf("public/affiliation/%d", id).Into(&a)
	if err != nil {
		return Affiliation{}, err
	}
	return a, err
}

func (s *Session) ListAllAffiliations() (a []Affiliation, err error) {
	err = s.get("public/affiliations").Into(&a)
	if err != nil {
		return []Affiliation{}, err
	}
	return a, err
}

func (s *Session) AddAffiliation(affiliation Affiliation, token string) (a Affiliation, err error) {
	a1, err := json.Marshal(affiliation)
	if err != nil {
		return Affiliation{}, err
	}
	err = s.putToken(token, "internal/affiliation", *bytes.NewBuffer(a1)).Into(&a)
	if err != nil {
		return Affiliation{}, err
	}
	return a, err
}

func (s *Session) DeleteAffiliation(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/affiliation/%d", id).JSON()
	return err
}
