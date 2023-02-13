package handler

import (
	"bytes"
	"encoding/json"
)

type WhatsOn struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	FileName    string `json:"file_name,omitempty"`
	Content     string `json:"content"`
	Date        int64  `json:"date"`
	DateOfEvent int64  `json:"date_of_event"`
	Delete      bool   `json:"delete,omitempty"`
}

func (s *Session) GetWhatsOnById(id uint64) (w WhatsOn, err error) {
	err = s.getf("public/whatsOn/%d", id).Into(&w)
	if err != nil {
		return WhatsOn{}, err
	}
	return w, err
}

func (s *Session) GetWhatsOnLatest() (w WhatsOn, err error) {
	err = s.get("public/whatsOn/latest").Into(&w)
	if err != nil {
		return WhatsOn{}, err
	}
	return w, err
}

func (s *Session) ListAllWhatsOnEventPast() (w []WhatsOn, err error) {
	err = s.get("public/whatsOn/past").Into(&w)
	if err != nil {
		return []WhatsOn{}, err
	}
	return w, err
}

func (s *Session) ListAllWhatsOn() (w []WhatsOn, err error) {
	err = s.get("public/whatsOn").Into(&w)
	if err != nil {
		return []WhatsOn{}, err
	}
	return w, err
}

func (s *Session) ListAllWhatsOnEventFuture() (w []WhatsOn, err error) {
	err = s.get("public/whatsOn/future").Into(&w)
	if err != nil {
		return []WhatsOn{}, err
	}
	return w, err
}

func (s *Session) AddWhatsOn(whatsOn WhatsOn, token string) (w WhatsOn, err error) {
	w1, err := json.Marshal(whatsOn)
	if err != nil {
		return WhatsOn{}, err
	}
	err = s.putToken(token, "internal/whatsOn", *bytes.NewBuffer(w1)).Into(&w)
	if err != nil {
		return WhatsOn{}, err
	}
	return w, err
}

func (s *Session) EditWhatsOn(whatsOn WhatsOn, token string) (w WhatsOn, err error) {
	w1, err := json.Marshal(whatsOn)
	if err != nil {
		return WhatsOn{}, err
	}
	err = s.patchToken(token, "internal/whatsOn", *bytes.NewBuffer(w1)).Into(&w)
	if err != nil {
		return WhatsOn{}, err
	}
	return w, err
}

func (s *Session) DeleteWhatsOn(id uint64, token string) (err error) {
	_, err = s.deletefToken(token, "internal/whatsOn/%d", id).JSON()
	return err
}
