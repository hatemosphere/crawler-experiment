package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/hatemosphere/crawler-experiment/crawler"
)

const concurrency = 10

type BasicParser struct {
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func (d BasicParser) ParsePage(doc *goquery.Document) crawler.ScrapeResult {
	data := crawler.ScrapeResult{}
	data.Title = doc.Find("title").First().Text()
	data.H1 = doc.Find("h1").First().Text()
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
	PrettyPrint(crawlResults)
}
