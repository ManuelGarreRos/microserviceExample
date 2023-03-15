package main

import (
	"context"
	"fmt"
	"time"
)

type LoggingService struct {
	next Service
}

func NewLoggingService(next Service) Service {
	return &LoggingService{
		next: next,
	}
}

func (s *LoggingService) GetMessage(ctx context.Context) (fact *Message, err error) {
	defer func(start time.Time) {
		fmt.Printf("fact=%v err=%v took=%v \n", fact, err, time.Since(start))
	}(time.Now())

	return s.next.GetMessage(ctx)
}
