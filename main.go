package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/gocolly/colly"
)

const (
	inputError   = "Input error"
	browserError = "Browser error"
)

var sortOptions = map[string]string{
	"-n": "newest",
	"-a": "active",
	"-b": "bounties",
	"-u": "unanswered",
	"-f": "frequent",
	"-v": "votes",
}

func scanInput() string {
	var i string
	_, err := fmt.Scan(&i)
	if err != nil {
		panic(inputError)
	}

	return i
}

func main() {
	var msg, sort string
	var page int

	col := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	i := 0
	links := make([]string, 15)

	col.OnHTML("#questions .question-hyperlink", func(e *colly.HTMLElement) {
		fmt.Printf("%d %s\n", i, e.Text)
		links[i] = e.Attr("href")
		i++
	})

	var newPage = func(inc bool) {
		if inc {
			page++
		} else {
			if page > 1 {
				page--
			} else {
				return
			}
		}

		i = 0
		links = make([]string, 15)
		exec.Command("clear").Start()

		url := fmt.Sprintf("https://www.stackoverflow.com/questions/tagged/%s?tab=%s&page=%d&pagesize=15", os.Args[1], sort, page)
		fmt.Println(url)
		err := col.Visit(url)
		if err != nil {
			fmt.Println("error")
		}
	}

	var visitURL = func() {
		inputLink, err := strconv.Atoi(scanInput())
		for err == nil {
			fmt.Println("Valid inputs: ',', '.', 'e', '{0-14}'")
			inputLink, err = strconv.Atoi(scanInput())
		}

		err = os.Chdir("/")
		if err != nil {
			panic(browserError)
		}

		cmd := []string{
			"open",
			fmt.Sprintf("https://www.stackoverflow.com%s", links[inputLink]),
		}
		err = exec.Command(cmd[0], cmd[1]).Start()
		if err != nil {
			panic(err)
		}
	}

	switch len(os.Args) {
	case 1:
		msg = "Insert something to search"
	case 2:
		sort = sortOptions["-a"]
	case 3:
		s, ok := sortOptions[os.Args[2]]
		sort = s
		if !ok {
			msg = "Invalid sort option"
		}
	default:
		msg = "Too many arguments"
	}

	if msg != "" {
		fmt.Println(msg)
		os.Exit(1)
	}

	newPage(true)

	for s := scanInput(); s != "e"; s = scanInput() {
		switch s {
		case ",":
			newPage(false)
		case ".":
			newPage(true)
		default:
			visitURL()
		}
	}
}
