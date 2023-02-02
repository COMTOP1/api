package handler

import (
	"bytes"
	"encoding/json"
)

type WhatsOn struct {
	Id            uint64 `json:"id"`
	Title         string `json:"title"`
	ImageFileName string `json:"image,omitempty"`
	Content       string `json:"content"`
	Date          uint64 `json:"date"`
	DateOfEvent   uint64 `json:"dateOfEvent"`
	Delete        bool   `json:"delete,omitempty"`
}

func (s *Session) GetWhatsOnById(id uint64) (whatsOn WhatsOn, err error) {
	err = s.getf("public/whatsOn/%d", id).Into(&whatsOn)
	return whatsOn, err
}

func (s *Session) GetWhatsOnLatest() (whatsOn WhatsOn, err error) {
	err = s.get("public/whatsOn/latest").Into(&whatsOn)
	return whatsOn, err
}

func (s *Session) ListAllWhatsOn() (whatsOn []WhatsOn, err error) {
	err = s.get("public/whatsOn").Into(&whatsOn)
	return whatsOn, err
}

func (s *Session) AddWhatsOn(whatsOn WhatsOn, token string) (whatsOns WhatsOn, err error) {
	whatsOn1, err := json.Marshal(whatsOn)
	if err != nil {
		return WhatsOn{}, err
	}
	err = s.putToken(token, "internal/whatsOn", *bytes.NewBuffer(whatsOn1)).Into(&whatsOns)
	return whatsOns, err
}

func (s *Session) EditWhatsOn(whatsOn WhatsOn, token string) (whatsOns WhatsOn, err error) {
	whatsOn1, err := json.Marshal(whatsOn)
	if err != nil {
		return WhatsOn{}, err
	}
	err = s.patchToken(token, "internal/whatsOn", *bytes.NewBuffer(whatsOn1)).Into(&whatsOns)
	return whatsOns, err
}

func (s *Session) DeleteWhatsOn(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/whatsOn/%d", id).JSON()
	return err
}
