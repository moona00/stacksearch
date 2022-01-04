package main

import (
	"fmt"
	"os"

	"./stacksearch"

	"github.com/gocolly/colly"
)

func main() {
	var msg, sort string

	switch len(os.Args) {
	case 1:
		msg = "Please insert something to search"
	case 2:
		sort = stacksearch.SortOptions["-r"]
	case 3:
		s, ok := stacksearch.SortOptions[os.Args[2]]
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

	c := colly.NewCollector()
	c.Visit(fmt.Sprintf("https://www.stackoverflow.com/search?tab=%s&q=%s", sort, os.Args[1]))
}
