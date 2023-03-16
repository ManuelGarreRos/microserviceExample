package main

import (
	"context"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"strconv"
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
	pdf.CellFormat(210, lineHt, "Circle sample report", "", 0, "C", false, 0, "")
	circles := []Circle{
		{"5678", "Diabetes", 5},
		{"9123", "Esguince de rodilla", 8},
		{"0000", "Niño sano", 20},
		{"2001", "Esguince de tobillo", 7},
		{"1111", "Epilepsia", 2},
		{"3000", "Oncologia", 15},
	}

	var maxUsers int
	for _, c := range circles {
		if c.users > maxUsers {
			maxUsers = c.users
		}
	}

	// Set X and Y position for drawing
	x := 20.0
	y := 60.0

	// Set width and height of each bar
	barWidth := 11.0
	barHeight := 100.0

	pdf.SetFont("Arial", "B", 6)

	// Draw bars for each circle
	for i, c := range circles {
		// Calculate X and Y position for this bar
		xPos := x + (float64(i) * barWidth * 2.5)
		yPos := y + barHeight

		// Calculate bar height based on circle users
		bar := float64(c.users) / float64(maxUsers) * barHeight

		// Draw the bar
		pdf.Rect(xPos, yPos-bar, barWidth, bar, "F")

		// Draw the circle name below the bar
		pdf.SetXY(xPos, yPos+20)
		pdf.Cell(barWidth, 10, c.name)

		// Draw the number of users above the bar
		pdf.SetXY(xPos, yPos+10)
		pdf.Cell(barWidth, 10, strconv.Itoa(c.users))
	}

	err := pdf.OutputFileAndClose("circle-report.pdf")
	if err != nil {
		panic(err)
	}

	return nil
}
