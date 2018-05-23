# sheets-save-all
Save all sheets in a Google Sheets doc as CSV. Requires an [OAuth 2.0 token](https://support.google.com/googleapi/answer/6158849) to access Google Sheets.


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
        Google client secret (default $GOOGLE_CLIENT_SECRET)
  -crlf
        use Windows-style line endings
  -filename string
        file name for files (default "{{.Properties.Index}} {{.Properties.Title}}.csv")
  -path string
        path to save files in (default "{{.Properties.Title}}")
  -quiet
        don't log activity
  -sheet string
        Google Sheet ID
```
