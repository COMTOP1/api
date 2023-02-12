package handler

import (
	"bytes"
	"encoding/json"
)

type Team struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	League        string `json:"league,omitempty"`
	Division      string `json:"division,omitempty"`
	LeagueTable   string `json:"league_table,omitempty"`
	Fixtures      string `json:"fixtures,omitempty"`
	Coach         string `json:"coach,omitempty"`
	Physio        string `json:"physio,omitempty"`
	ImageFileName string `json:"image,omitempty"`
	Active        bool   `json:"active"`
	Youth         bool   `json:"youth"`
	Ages          int    `json:"ages"`
}

func (s *Session) GetTeamById(id uint64) (t Team, err error) {
	err = s.getf("public/team/%d", id).Into(&t)
	if err != nil {
		return Team{}, err
	}
	return t, err
}

func (s *Session) ListAllTeams(token string) (t []Team, err error) {
	err = s.getToken(token, "internal/teams").Into(&t)
	if err != nil {
		return []Team{}, err
	}
	return t, err
}

func (s *Session) ListActiveTeams() (t []Team, err error) {
	err = s.get("public/teams").Into(&t)
	if err != nil {
		return []Team{}, err
	}
	return t, err
}

func (s *Session) AddTeam(team Team, token string) (t Team, err error) {
	t1, err := json.Marshal(team)
	if err != nil {
		return Team{}, err
	}
	err = s.putToken(token, "internal/team", *bytes.NewBuffer(t1)).Into(&t)
	if err != nil {
		return Team{}, err
	}
	return t, err
}

func (s *Session) EditTeam(team Team, token string) (t Team, err error) {
	t1, err := json.Marshal(team)
	if err != nil {
		return Team{}, err
	}
	err = s.patchToken(token, "internal/team", *bytes.NewBuffer(t1)).Into(&t)
	if err != nil {
		return Team{}, err
	}
	return t, err
}

func (s *Session) DeleteTeam(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/team/%d", id).JSON()
	return err
}
