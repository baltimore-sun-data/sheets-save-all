package sheets

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	spreadsheet "gopkg.in/Iwark/spreadsheet.v2"
)

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("sheets-save-all", flag.ExitOnError)
	fl.StringVar(&conf.SheetID, "sheet", "", "Google Sheet ID (default $GOOGLE_CLIENT_SECRET)")
	fl.StringVar(&conf.ClientSecret, "client-secret", "", "Google client secret")
	fl.StringVar(&conf.PathTemplate, "path", "{{.Properties.Title}}", "path to save files in")
	fl.StringVar(&conf.FileTemplate, "filename", "{{.Properties.Index}} {{.Properties.Title}}.csv",
		"file name for files")
	fl.BoolVar(&conf.UseCRLF, "crlf", false, "use Windows-style line endings")

	quiet := fl.Bool("quiet", false, "don't log activity")
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`sheets-save-all is a tool to save all sheets in Google Sheets document.

-path and -filename are Go templates and can use any property of the document or sheet object respectively. See gopkg.in/Iwark/spreadsheet.v2 for properties.

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
	PathTemplate string
	FileTemplate string
	UseCRLF      bool
	Logger       *log.Logger
}

func (c *Config) Exec() error {
	if c.SheetID == "" {
		return ErrNoSheet
	}

	pt, err := template.New("path").Parse(c.PathTemplate)
	if err != nil {
		return fmt.Errorf("path template problem: %v", err)
	}
	ft, err := template.New("file").Parse(c.FileTemplate)
	if err != nil {
		return fmt.Errorf("file path template problem: %v", err)
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

	var buf strings.Builder
	if err = pt.Execute(&buf, doc); err != nil {
		return fmt.Errorf("could not use path template: %v", err)
	}
	dir := buf.String()

	os.MkdirAll(dir, os.ModePerm)
	for _, s := range doc.Sheets {
		buf.Reset()
		if err = ft.Execute(&buf, s); err != nil {
			return fmt.Errorf("could not use file path template: %v", err)
		}
		file := buf.String()

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
	w.UseCRLF = c.UseCRLF
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
