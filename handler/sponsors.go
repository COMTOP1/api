package handler

import (
	"bytes"
	"encoding/json"
)

type Sponsor struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Website  string `json:"website"`
	FileName string `json:"file_name"`
	Purpose  string `json:"purpose"`
	Team     string `json:"team"`
}

func (s *Session) GetSponsorById(id uint64) (sp Sponsor, err error) {
	err = s.getf("public/sponsor/%d", id).Into(&sp)
	if err != nil {
		return Sponsor{}, err
	}
	return sp, err
}

func (s *Session) ListAllSponsors() (sp []Sponsor, err error) {
	err = s.get("public/sponsors").Into(&sp)
	if err != nil {
		return []Sponsor{}, err
	}
	return sp, err
}

func (s *Session) ListAllSponsorsMinimal() (sp []Sponsor, err error) {
	err = s.get("public/sponsors/minimal").Into(&sp)
	if err != nil {
		return []Sponsor{}, err
	}
	return sp, err
}

func (s *Session) ListAllSponsorsByTeamId(teamId uint64) (sponsors []Sponsor, err error) {
	err = s.getf("public/sponsors/%d", teamId).Into(&sponsors)
	if err != nil {
		return []Sponsor{}, err
	}
	return sponsors, err
}

func (s *Session) AddSponsor(sp1 Sponsor, token string) (sp Sponsor, err error) {
	sp2, err := json.Marshal(sp1)
	if err != nil {
		return Sponsor{}, err
	}
	err = s.putToken(token, "internal/sponsor", *bytes.NewBuffer(sp2)).Into(&sp)
	if err != nil {
		return Sponsor{}, err
	}
	return sp, err
}

func (s *Session) EditSponsor(sp1 Sponsor, token string) (sp Sponsor, err error) {
	sp2, err := json.Marshal(sp1)
	if err != nil {
		return Sponsor{}, err
	}
	err = s.patchToken(token, "internal/sponsor", *bytes.NewBuffer(sp2)).Into(&sp)
	if err != nil {
		return Sponsor{}, err
	}
	return sp, err
}

func (s *Session) DeleteSponsor(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/sponsor/%d", id).JSON()
	return err
}
