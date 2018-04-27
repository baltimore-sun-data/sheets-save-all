package main

import (
	"fmt"
	"os"

	"github.com/baltimore-sun-data/sheets-save-all/sheets"
)

func main() {
	c := sheets.FromArgs(os.Args[1:])
	if err := c.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
