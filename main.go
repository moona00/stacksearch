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
	baseURL      = "https://www.stackoverflow.com"
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
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		panic(inputError)
	}

	return input
}

func main() {
	var sort string
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

		url := fmt.Sprintf(
			"%s/questions/tagged/%s?tab=%s&page=%d&pagesize=15",
			baseURL,
			os.Args[1],
			sort,
			page,
		)
		err := col.Visit(url)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var visitInputURL = func(input string) {
		link, _ := strconv.Atoi(input)

		cmd := []string{
			"open",
			fmt.Sprintf("%s%s", baseURL, links[link]),
		}

		err := exec.Command(cmd[0], cmd[1]).Start()
		if err != nil {
			panic(err)
		}
	}

	switch len(os.Args) {
	case 1:
		panic("Insert something to search")
	case 2:
		sort = sortOptions["-a"]
	case 3:
		s, ok := sortOptions[os.Args[2]]
		sort = s
		if !ok {
			panic("Invalid sort option")
		}
	default:
		panic("Too many arguments")
	}

	newPage(true)

	for s := scanInput(); s != "e"; s = scanInput() {
		switch s {
		case ",":
			newPage(false)
		case ".":
			newPage(true)
		default:
			visitInputURL(s)
		}
	}
}
