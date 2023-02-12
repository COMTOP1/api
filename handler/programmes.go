package handler

import (
	"bytes"
	"encoding/json"
)

type Programme struct {
	Id                uint64 `json:"id"`
	Name              string `json:"name"`
	FileName          string `json:"file_name"`
	DateOfProgramme   int64  `json:"date_of_programme"`
	ProgrammeSeasonId uint64 `json:"programme_season_id"`
}

func (s *Session) GetProgrammeById(id uint64) (p Programme, err error) {
	err = s.getf("public/programme/%d", id).Into(&p)
	if err != nil {
		return Programme{}, err
	}
	return p, err
}

func (s *Session) ListAllProgrammes() (p []Programme, err error) {
	err = s.get("public/programmes").Into(&p)
	if err != nil {
		return []Programme{}, err
	}
	return p, err
}

func (s *Session) AddProgramme(p1 Programme, token string) (p Programme, err error) {
	p2, err := json.Marshal(p1)
	if err != nil {
		return Programme{}, err
	}
	err = s.putToken(token, "internal/programme", *bytes.NewBuffer(p2)).Into(&p)
	if err != nil {
		return Programme{}, err
	}
	return p, err
}

func (s *Session) DeleteProgramme(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/programme/%d", id).JSON()
	return err
}
