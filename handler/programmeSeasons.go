package handler

import (
	"bytes"
	"encoding/json"
)

type ProgrammeSeason struct {
	Id     uint64 `json:"id"`
	Season string `json:"season"`
}

func (s *Session) GetProgrammeSeasonById(id uint64) (p ProgrammeSeason, err error) {
	err = s.getf("public/programmeSeason/%d", id).Into(&p)
	return p, err
}

func (s *Session) ListAllProgrammeSeasons() (p []ProgrammeSeason, err error) {
	err = s.get("public/programmeSeasons").Into(&p)
	return p, err
}

func (s *Session) AddProgrammeSeason(programmeSeason ProgrammeSeason, token string) (p ProgrammeSeason, err error) {
	programmeSeason1, err := json.Marshal(programmeSeason)
	if err != nil {
		return ProgrammeSeason{}, err
	}
	err = s.putToken(token, "internal/programmeSeason", *bytes.NewBuffer(programmeSeason1)).Into(&p)
	return p, err
}

func (s *Session) EditProgrammeSeason(programmeSeason ProgrammeSeason, token string) (p ProgrammeSeason, err error) {
	programmeSeason1, err := json.Marshal(programmeSeason)
	if err != nil {
		return ProgrammeSeason{}, err
	}
	err = s.patchToken(token, "internal/programmeSeason", *bytes.NewBuffer(programmeSeason1)).Into(&p)
	return p, err
}

func (s *Session) DeleteProgrammeSeason(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/programmeSeason/%d", id).JSON()
	return err
}
