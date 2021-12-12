package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ahrberg/svttext/internal/svt"
	tea "github.com/charmbracelet/bubbletea"
)

// Set by ldflags
var version string

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
	interactive := flag.Bool("interactive", false, "start interactive mode\nuse arrow keys to navigate pages\nor enter page number to go to page")
	versionFlag := flag.Bool("version", false, "prints svttext version")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

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

	// Output
	if *interactive {
		p := tea.NewProgram(model{
			page:   page,
			colors: *colors,
			client: client,
		}, tea.WithAltScreen())
		if err := p.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		res, err := client.GetNews(page)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		// Color output if decired
		var text string

		for _, page := range res.Text {
			text += page
		}

		if *colors {
			text = svt.ColorPage(text)
		}

		fmt.Print(text)
	}
}

func pageValid(page string) bool {
	r := regexp.MustCompile(`^\d{3}$`)
	return r.MatchString(page)
}

func getPage(client *svt.Client, page string, colors bool) tea.Cmd {
	return func() tea.Msg {
		res, err := client.GetNews(page)

		if err != nil {
			return news{
				errMsg: err.Error(),
			}
		}

		return news{
			text:       res.Text,
			pageNumber: res.PageNumber,
			prevPage:   res.PrevPage,
			nextPage:   res.NextPage,
		}
	}
}

type news struct {
	text       []string
	errMsg     string
	pageNumber string
	prevPage   string
	nextPage   string
}

type model struct {
	page         string // current page
	subPageIndex int    // sub page index
	pageInp      string // user input for page search
	news         news   // page content with details
	colors       bool   // color output
	client       *svt.Client
}

func (m model) Init() tea.Cmd {
	return getPage(m.client, m.page, m.colors)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case news:
		m.news = msg
		return m, nil

	case tea.KeyMsg:

		// Switch over key pressed
		switch msg.String() {

		// Clear page search
		case "esc":
			if m.pageInp != "" {
				m.pageInp = ""
				return m, nil
			}

		// Start page search
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if len(m.pageInp) == 2 {
				m.page = m.pageInp + msg.String()
				m.pageInp = ""
				m.subPageIndex = 0
				return m, getPage(m.client, m.page, m.colors)
			} else {
				m.pageInp += msg.String()
			}

		// Exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// Search page decremental
		case "left", "h":
			if m.news.prevPage != "" && (m.news.prevPage < m.page) {
				m.page = m.news.prevPage
				m.subPageIndex = 0
				return m, getPage(m.client, m.page, m.colors)
			} else {
				return m, nil
			}

		// Search page incremental
		case "right", "l":
			if m.news.nextPage != "" && (m.news.nextPage > m.page) {
				m.page = m.news.nextPage
				m.subPageIndex = 0
				return m, getPage(m.client, m.page, m.colors)
			} else {
				return m, nil
			}
		// Change sub page
		case "up", "k":
			if m.subPageIndex > 0 {
				m.subPageIndex--
				return m, nil
			}
		case "down", "j":
			if m.subPageIndex < len(m.news.text)-1 {
				m.subPageIndex++
				return m, nil
			}
		}

	}

	return m, nil
}

func (m model) View() string {

	res := ""

	if m.pageInp != "" {
		res += fmt.Sprintf("🔎 %s\n\n", m.pageInp)
	}

	if len(m.news.text) > 0 {
		res += m.news.text[m.subPageIndex]
	}

	if m.news.errMsg != "" {
		res += m.news.errMsg
	}

	if m.colors {
		res = svt.ColorPage(res)
	}

	return res
}
