package main

import (
	"context"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

type Service interface {
	GetMessage(context.Context) (*Message, error)
	GetOk(ctx context.Context) (string, error)
	PdfReport(ctx context.Context) error
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

func (s *MessageService) GetOk(ctx context.Context) (string, error) {
	m := "OK"
	return m, nil
}

func (s *MessageService) PdfReport(ctx context.Context) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	_, lineHt := pdf.GetFontSize()
	pdf.SetXY(0, 20)
	pdf.CellFormat(210, lineHt, "Sample report", "", 0, "C", false, 0, "")
	pdf.SetLineWidth(0.5)
	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(10, 50, 100, 100, "F")
	pdf.Line(10, 75, 110, 75)
	err := pdf.OutputFileAndClose("sample-report.pdf")
	if err != nil {
		panic(err)
	}

	return nil
}
