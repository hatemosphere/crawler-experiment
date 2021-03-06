package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	URL   string
	Title string
	H1    string
}

type Parser interface {
	ParsePage(*goquery.Document) ScrapeResult
}

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func extractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	}
	return foundUrls
}

func resolveRelative(baseURL string, hrefs []string) []string {
	internalUrls := []string{}

	for _, href := range hrefs {
		if strings.HasPrefix(href, baseURL) {
			internalUrls = append(internalUrls, href)
		}

		if strings.HasPrefix(href, "/") {
			resolvedURL := fmt.Sprintf("%s%s", baseURL, href)
			internalUrls = append(internalUrls, resolvedURL)
		}
	}

	return internalUrls
}

func trimTrash(links []string) []string {
	validLinks := []string{}
	for _, l := range links {
		// According to https://tools.ietf.org/html/rfc3986#section-4.1
		// query should come first, so this ugly trimming should work
		if strings.Contains(l, "?") || strings.Contains(l, "#") {
			var index int
			for n, str := range l {
				if strconv.QuoteRune(str) == "'?'" || strconv.QuoteRune(str) == "'#'" {
					index = n
					break
				}
			}
			l = l[:index]
		}
		validLinks = append(validLinks, l)
	}
	return validLinks

}

func crawlPage(baseURL, targetURL string, parser Parser, token chan struct{}) ([]string, ScrapeResult) {

	token <- struct{}{}
	fmt.Println("Requesting: ", targetURL)
	resp, err := getRequest(targetURL)
	if err != nil {
		fmt.Println(err)
	}
	<-token

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	pageResults := parser.ParsePage(doc)
	links := extractLinks(doc)
	foundUrls := resolveRelative(baseURL, links)
	normalizedUrls := trimTrash(foundUrls)
	pageResults.URL = targetURL
	return normalizedUrls, pageResults
}

func parseStartURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

func Crawl(startURL string, parser Parser, concurrency int) []ScrapeResult {
	results := []ScrapeResult{}
	worklist := make(chan []string)
	var n int
	n++
	var tokens = make(chan struct{}, concurrency)
	go func() { worklist <- []string{startURL} }()
	seen := make(map[string]bool)
	baseDomain := parseStartURL(startURL)

	m := &sync.Mutex{}
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(baseDomain, link string, parser Parser, token chan struct{}) {
					foundLinks, pageResults := crawlPage(baseDomain, link, parser, token)
					m.Lock()
					results = append(results, pageResults)
					m.Unlock()
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(baseDomain, link, parser, tokens)
			}
		}
	}
	return results
}
