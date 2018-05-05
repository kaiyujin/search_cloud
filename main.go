package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"search_cloud/lib"
)

func main() {
	flag.Parse()
	k := flag.Arg(0)
	c := lib.Read()
	m := lib.SearchModel{c}
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan int)
	go p(m.FindBacklog, k, ch)
	go p(m.FindGithub, k, ch)
	go p(m.FindConfluence, k, ch)
	for i := 0; i < 3; i++ {
		<-ch
	}
}

func p(fn func(string) ([]lib.Result, error), keyword string, ch chan int) {
	results, err := fn(keyword)
	if err != nil {
		log.Fatal(results, err)
		os.Exit(1)
	}
	for _, m := range results {
		fmt.Println(m.URL, m.Title)
	}
	ch <- 0
}
