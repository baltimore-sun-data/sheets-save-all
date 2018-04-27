package sheets

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	spreadsheet "gopkg.in/Iwark/spreadsheet.v2"
)

var ClientSecret = []byte(os.Getenv("GOOGLE_CLIENT_SECRET"))

func Save(sheetID string) error {
	log.Printf("Connecting to Google Sheets for %q", sheetID)

	conf, err := google.JWTConfigFromJSON(ClientSecret, spreadsheet.Scope)
	if err != nil {
		return fmt.Errorf("could not parse credentials: %v", err)
	}

	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)
	doc, err := service.FetchSpreadsheet(sheetID)
	if err != nil {
		return fmt.Errorf("failure getting Google Sheet: %v", err)
	}

	dir := doc.Properties.Title
	os.MkdirAll(dir, os.ModePerm)
	for _, s := range doc.Sheets {
		file := fmt.Sprintf("%s.csv", s.Properties.Title)
		if err = makeCSV(dir, file, s.Rows); err != nil {
			return err
		}
	}
	return nil
}

func makeCSV(dir, file string, rows [][]spreadsheet.Cell) error {
	pathname := filepath.Join(dir, file)
	log.Printf("Writing file: %s", pathname)
	f, err := os.Create(pathname)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, row := range rows {
		record := make([]string, 0, len(row))
		for _, cell := range row {
			record = append(record, cell.Value)
		}
		if blank(record) {
			continue
		}
		err = w.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}

func blank(record []string) bool {
	for _, s := range record {
		if s != "" {
			return false
		}
	}
	return true
}
