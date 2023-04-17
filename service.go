package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"net/http"
	"os/exec"
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
		if c.Users > maxUsers {
			maxUsers = c.Users
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
		bar := float64(c.Users) / float64(maxUsers) * barHeight

		// Draw the bar
		pdf.Rect(xPos, yPos-bar, barWidth, bar, "F")

		// Draw the circle name below the bar
		pdf.SetXY(xPos, yPos+20)
		pdf.Cell(barWidth, 10, c.Name)

		// Draw the number of users above the bar
		pdf.SetXY(xPos, yPos+10)
		pdf.Cell(barWidth, 10, strconv.Itoa(c.Users))
	}

	//err := pdf.OutputFileAndClose("circle-report.pdf")
	//if err != nil {
	//	panic(err)
	//}

	// Create a new Excel file
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Add a sheet named "Data" to the file
	datasheet, err := file.NewSheet("Data")
	if err != nil {
		fmt.Println(datasheet)
		fmt.Println(err)
	}

	// Add headers to the data sheet
	headers := []string{"ID", "Name", "Users"}
	cell := []string{"A1", "B1", "C1"}
	for i, header := range headers {
		file.SetCellValue("Data", cell[i], header)
	}

	for i, circle := range circles {
		row := i + 2
		file.SetCellValue("Data", "A"+fmt.Sprint(row), circle.ID)
		file.SetCellValue("Data", "B"+fmt.Sprint(row), circle.Name)
		file.SetCellValue("Data", "C"+fmt.Sprint(row), circle.Users)
	}
	file.SetActiveSheet(datasheet)

	// Add a sheet named "Report" to the file
	reportSheet, err := file.NewSheet("Report")
	if err != nil {
		fmt.Println(reportSheet)
		fmt.Println(err)
	}

	// Save the file
	//err = file.SaveAs("circles.xlsx")
	//if err != nil {
	//	log.Fatal(err)
	//}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for idx, row := range [][]interface{}{
		{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
		{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	} {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			}},
		Title: excelize.ChartTitle{
			Name: "Fruit 3D Clustered Column Chart",
		},
	}); err != nil {
		return err
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		return err
	}

	// transform the file Book1.xlsx to Book1.pdf
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "Book1.xlsx")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
