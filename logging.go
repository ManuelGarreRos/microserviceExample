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

func (s *LoggingService) GetOk(ctx context.Context) (m string, err error) {
	defer func(start time.Time) {
		fmt.Printf("Message=%v err=%v took=%v \n", m, err, time.Since(start))
	}(time.Now())

	return s.next.GetOk(ctx)
}

func (s *LoggingService) PdfReport(ctx context.Context) (err error) {
	defer func(start time.Time) {
		if err != nil {
			fmt.Printf("PDF generator get err=%v; took=%v \n", err, time.Since(start))
		} else {
			fmt.Printf("Generate a report PDF; took=%v \n", time.Since(start))
		}
	}(time.Now())

	return s.next.PdfReport(ctx)
}
