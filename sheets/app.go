package sheets

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	spreadsheet "gopkg.in/Iwark/spreadsheet.v2"
)

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("sheets-save-all", flag.ExitOnError)
	fl.StringVar(&conf.SheetID, "sheet", "", "Google Sheet ID (default $GOOGLE_CLIENT_SECRET)")
	fl.StringVar(&conf.ClientSecret, "client-secret", "", "Google client secret")
	quiet := fl.Bool("quiet", false, "don't log activity")
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`sheets-save-all is a tool to save all sheets in Google Sheets document.

Usage of sheets-save-all:

`,
		)
		fl.PrintDefaults()
	}
	_ = fl.Parse(args)

	if conf.ClientSecret == "" {
		conf.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}

	if *quiet {
		conf.Logger = log.New(ioutil.Discard, "", 0)
	} else {
		conf.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	return conf
}

var (
	ErrNoSheet error = fmt.Errorf("No sheet ID provided")
)

type Config struct {
	SheetID      string
	ClientSecret string
	Logger       *log.Logger
}

func (c *Config) Exec() error {
	if c.SheetID == "" {
		return ErrNoSheet
	}

	c.Logger.Printf("Connecting to Google Sheets for %q", c.SheetID)

	conf, err := google.JWTConfigFromJSON([]byte(c.ClientSecret), spreadsheet.Scope)
	if err != nil {
		return fmt.Errorf("could not parse credentials: %v", err)
	}

	client := conf.Client(context.Background())
	service := spreadsheet.NewServiceWithClient(client)
	doc, err := service.FetchSpreadsheet(c.SheetID)
	if err != nil {
		return fmt.Errorf("failure getting Google Sheet: %v", err)
	}

	dir := doc.Properties.Title
	os.MkdirAll(dir, os.ModePerm)
	for _, s := range doc.Sheets {
		file := fmt.Sprintf("%s.csv", s.Properties.Title)
		if err = c.makeCSV(dir, file, s.Rows); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) makeCSV(dir, file string, rows [][]spreadsheet.Cell) error {
	pathname := filepath.Join(dir, file)
	c.Logger.Printf("Writing file: %s", pathname)
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
