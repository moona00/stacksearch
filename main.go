package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gocolly/colly"
)

const (
	inputError   = "Input error"
	browserError = "Browser error"
)

var sortOptions = map[string]string{
	"-r": "relevance",
	"-n": "newest",
	"-a": "active",
	"-v": "votes",
}

func main() {
	var msg, sort string

	switch len(os.Args) {
	case 1:
		msg = "Please insert something to search"
	case 2:
		sort = sortOptions["-r"]
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

	col := colly.NewCollector()

	i := 0
	links := make([]string, 15)

	col.OnHTML("#questions .question-hyperlink", func(e *colly.HTMLElement) {
		fmt.Printf("%d %s\n", i, e.Text)
		links[i] = e.Attr("href")
		i++
	})

	url := fmt.Sprintf("https://www.stackoverflow.com/questions/tagged/%s?tab=%s", os.Args[1], sort)
	fmt.Println(url)

	err := col.Visit(url)
	if err != nil {
		fmt.Println("error")
	}

	var inputLink int
	_, err = fmt.Scan(&inputLink)
	if err != nil {
		panic(inputError)
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
