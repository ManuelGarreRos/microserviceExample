package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetMessage(context.Context) (*Message, error)
}

type MessageService struct {
	url string
}

func NewMessageService(url string) Service {
	return &MessageService{
		url: url,
	}
}

func (s *MessageService) GetMessage(ctx context.Context) (*Message, error) {
	res, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fact := &Message{}
	if err := json.NewDecoder(res.Body).Decode(&fact); err != nil {
		return nil, err
	}

	return fact, nil
}
