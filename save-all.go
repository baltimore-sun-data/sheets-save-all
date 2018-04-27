package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/baltimore-sun-data/sheets-save-all/sheets"
)

func main() {
	flag.Parse()
	doc := flag.Arg(0)
	if err := sheets.Save(doc); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
