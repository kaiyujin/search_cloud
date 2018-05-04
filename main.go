package main

import (
	"flag"
	"fmt"
	"os"
	"search_cloud/lib"
)

func main() {
	//go routineを使う
	flag.Parse()
	k := flag.Arg(0)
	c := lib.Read()
	m := lib.SearchModel{c}
	p(m.FindBacklog, k)
	//p(model.FindGithub)
}

func p(fn func(string) ([]lib.Result, error), keyword string) {
	results, err := fn(keyword)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, m := range results {
		fmt.Println(m.URL, m.Title)
	}
}
