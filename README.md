# sheets-save-all
Save all sheets in a Google Sheets doc as CSV. Requires an API token to access Google Sheets.


## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```shell
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/baltimore-sun-data/sheets-save-all
```


## Usage
```shell
$ sheets-save-all -h
sheets-save-all is a tool to save all sheets in Google Sheets document.

-path and -filename are Go templates and can use any property of the document
or sheet object respectively. See gopkg.in/Iwark/spreadsheet.v2 for properties.

Usage of sheets-save-all:

  -client-secret string
        Google client secret
  -filename string
        File name for files (default "{{.Properties.Index}} {{.Properties.Title}}.csv")
  -path string
        Path to save files in (default "{{.Properties.Title}}")
  -quiet
        don't log activity
  -sheet string
        Google Sheet ID (default $GOOGLE_CLIENT_SECRET)
```
