package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

func main() {
	flag.Parse()
	doc := flag.Arg(0)
	if err := fromSheet(doc); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

var sheetsClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

func fromSheet(sheetID string) error {
	log.Printf("Connecting to Google Sheets for %q", sheetID)

	conf, err := google.JWTConfigFromJSON([]byte(sheetsClientSecret), spreadsheet.Scope)
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
			return nil
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
