package main

import (
	"flag"
	"fmt"
	"os"

	"./crawler" // TODO: move to absolute import
	"github.com/PuerkitoBio/goquery"
)

var concurrency = 10

type BasicParser struct {
}

func (d BasicParser) ParsePage(doc *goquery.Document) crawler.ScrapeResult {
	data := crawler.ScrapeResult{}
	data.Title = doc.Find("title").First().Text()
	data.H1 = doc.Find("h1").First().Text()
	data.URL = "" // TODO: implement
	return data
}

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	p := BasicParser{}
	crawlResults := crawler.Crawl(args[0], p, concurrency)
	fmt.Printf("%+v\n", crawlResults)
}
