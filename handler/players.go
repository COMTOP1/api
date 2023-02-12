package handler

import (
	"bytes"
	"encoding/json"
)

type Player struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	FileName    string `json:"file_name"`
	DateOfBirth int64  `json:"date_of_birth"`
	Captain     bool   `json:"captain"`
	Team        uint64 `json:"team"`
}

func (s *Session) GetPlayerById(id uint64) (p Player, err error) {
	err = s.getf("public/player/%d", id).Into(&p)
	if err != nil {
		return Player{}, err
	}
	return p, err
}

func (s *Session) ListAllPlayers(token string) (p []Player, err error) {
	err = s.getToken(token, "public/players").Into(&p)
	if err != nil {
		return []Player{}, err
	}
	return p, err
}

func (s *Session) ListAllPlayersByTeam(token string, teamId uint64) (p []Player, err error) {
	err = s.getfToken(token, "public/players/%d", teamId).Into(&p)
	if err != nil {
		return []Player{}, err
	}
	return p, err
}

func (s *Session) AddPlayer(player Player, token string) (p Player, err error) {
	p1, err := json.Marshal(player)
	if err != nil {
		return Player{}, err
	}
	err = s.putToken(token, "internal/player", *bytes.NewBuffer(p1)).Into(&p)
	if err != nil {
		return Player{}, err
	}
	return p, err
}

func (s *Session) EditPlayer(player Player, token string) (p Player, err error) {
	p1, err := json.Marshal(player)
	if err != nil {
		return Player{}, err
	}
	err = s.patchToken(token, "internal/player", *bytes.NewBuffer(p1)).Into(&p)
	if err != nil {
		return Player{}, err
	}
	return p, err
}

func (s *Session) DeletePlayer(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/player/%d", id).JSON()
	return err
}
