package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"shellycs50.com/crawler/userinterface"
)

type myFirstFetchStruct struct {
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string, mu *sync.Mutex, data *[]byte, prevFetchedUrls *map[string]struct{}, allowInternalLinks bool) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, prevFetchedUrls *map[string]struct{}, mu *sync.Mutex, data *[]byte, allowInternalLinks bool) {

	defer wg.Done()

	if depth <= 0 {
		return
	}

	_, foundUrls, err := fetcher.Fetch(url, mu, data, prevFetchedUrls, allowInternalLinks)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, newUrl := range foundUrls {
		shouldCrawl := false
		mu.Lock()
		if _, ok := (*prevFetchedUrls)[newUrl]; !ok {
			(*prevFetchedUrls)[newUrl] = struct{}{}
			shouldCrawl = true
		}
		mu.Unlock()
		if shouldCrawl {
			wg.Add(1)
			go Crawl(newUrl, depth-1, fetcher, wg, prevFetchedUrls, mu, data, allowInternalLinks)
		}
	}

}

func main() {
	// by using a map and a byte slice we gain the ability to efficiently check if a url has been queried before, but
	// ONLY add it to output.txt if a response is received. However, we are using (in the realm of) double the memory.

	var urls = make(map[string]struct{})
	var data = []byte("")
	// using mutex to lock the map (formally it's locking code that accesses the map), using map for constant time lookups. (previously used slice but with the amount of blocking you could argue quicker with a channel for small datasets but for large datasets a map is quicker.)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var this_is_confusing Fetcher = myFirstFetchStruct{}
	user_url, intdepth, user_filename, allow_internal_links := userinterface.GetCrawlArgs()
	wg.Add(1)
	Crawl(user_url, intdepth, this_is_confusing, &wg, &urls, &mu, &data, allow_internal_links)
	wg.Wait()
	os.WriteFile(user_filename+".txt", data, 0644)
}

func (f myFirstFetchStruct) Fetch(url string, mu *sync.Mutex, data *[]byte, prevFetchedUrls *map[string]struct{}, allowInternalLinks bool) (body string, urls []string, err error) {
	var url_list []string
	c := colly.NewCollector()
	c.OnResponse(func(*colly.Response) {
		fmt.Printf("Visited %v\n", url)
		mu.Lock()
		*data = append(*data, []byte(url+"\n")...)
		mu.Unlock()
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if allowInternalLinks && !isExternalLink(link) {
			link = prependDomain(url, link)
		}
		url_list = append(url_list, link)
	})
	c.Visit(url)
	return "", url_list, nil
}

func isExternalLink(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}

func prependDomain(url, link string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") || len(url) < 9 {
		// url is very weird get out of here now.
		return link
	}
	if strings.Contains(url[8:], "/") {
		domainEndIndex := strings.Index(url[8:], "/") + 8
		return url[:domainEndIndex] + link
	}
	return url + "/" + link
}
