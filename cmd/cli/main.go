package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ahrberg/svttext/internal/svt"
)

func main() {
	// Program usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: svttext [OPTION]... [PAGE]\n")
		fmt.Fprintf(os.Stderr, "Read news from SVT Text\n\n")
		fmt.Fprintf(os.Stderr, "Example: svttext --colors 100\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Program inputs
	colors := flag.Bool("colors", false, "colorize the output")

	flag.Parse()

	page := "100" // Default to page 100 if nothing else specified

	if len(os.Args) >= 2 && !strings.HasPrefix(os.Args[len(os.Args)-1], "-") {
		page = os.Args[len(os.Args)-1]
	}

	if !pageValid(page) {
		fmt.Fprintf(os.Stderr, "Error: Page number not valid, must be nnn. Use `svttext --help` for more help.\n")
		os.Exit(1)
	}

	// Get news from SVT
	client := svt.NewClient()
	text, err := client.GetNews(page)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}

	// Color output if decired
	if *colors {
		text = svt.ColorPage(text)
	}

	// Output
	fmt.Print(text)
}

func pageValid(page string) bool {
	r := regexp.MustCompile(`^\d{3}$`)
	return r.MatchString(page)
}
