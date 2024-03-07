package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type myFirstFetchStruct struct {
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string, mu *sync.Mutex, data *[]byte, prevFetchedUrls *map[string]struct{}) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, prevFetchedUrls *map[string]struct{}, mu *sync.Mutex, data *[]byte) {

	defer wg.Done()

	if depth <= 0 {
		return
	}

	_, foundUrls, err := fetcher.Fetch(url, mu, data, prevFetchedUrls)
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
			go Crawl(newUrl, depth-1, fetcher, wg, prevFetchedUrls, mu, data)
		}
	}

}

func main() {
	var urls = make(map[string]struct{})
	var data = []byte("")
	// using mutex to lock the map, using map for constant time lookups. (previously used slice but with the amount of blocking you could argue quicker with a channel for small datasets but for large datasets a map is quicker.)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var this_is_confusing Fetcher = myFirstFetchStruct{}
	wg.Add(1)
	Crawl("https://muz.li/blog/60-most-creative-portfolio-websites-of-2023", 5, this_is_confusing, &wg, &urls, &mu, &data)
	wg.Wait()
	os.WriteFile("./urls.txt", data, 0644)
}

func (f myFirstFetchStruct) Fetch(url string, mu *sync.Mutex, data *[]byte, prevFetchedUrls *map[string]struct{}) (body string, urls []string, err error) {
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
		url_list = append(url_list, link)
	})
	c.Visit(url)
	return "", url_list, nil
}
