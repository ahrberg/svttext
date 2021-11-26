package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ahrberg/svttext/internal/svt"
)

func main() {
	client := svt.NewClient()

	// Program inputs
	page := os.Args[1]

	// Get news page from SVT
	text, err := client.GetNews(page)

	if err != nil {
		log.Fatalf("Error getting page: %s", err.Error())
	}

	fmt.Print(text)
}
