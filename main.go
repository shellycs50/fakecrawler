package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

// type Fetcher interface {
// 	// Fetch returns the body of URL and
// 	// a slice of URLs found on that page.
// 	Fetch(url string) (body string, urls []string, err error)
// }

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher func(string, *sync.Mutex, *[]byte, *map[string]struct{}) (string, []string, error), wg *sync.WaitGroup, prevFetchedUrls *map[string]struct{}, mu *sync.Mutex, data *[]byte) {

	defer wg.Done()

	fmt.Println("Fetching urls from ", url)

	if depth <= 0 {
		return
	}

	_, urls, err := fetcher(url, mu, data, prevFetchedUrls)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, wg, prevFetchedUrls, mu, data)
	}

}

func main() {
	var urls = make(map[string]struct{})
	var data = []byte("")
	// using mutex to lock the map, using map for constant time lookups. (previously used slice but with the amount of blocking you could argue quicker with a channel for small datasets but for large datasets a map is quicker.)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("https://io-academy.uk/", 9, Fetch, &wg, &urls, &mu, &data)
	wg.Wait()
	// os.WriteFile("./urls.txt", []byte(fmt.Sprintf("%v", urls)), 0644)
	os.WriteFile("./urls.txt", data, 0644)
}

func Fetch(url string, mu *sync.Mutex, data *[]byte, prevFetchedUrls *map[string]struct{}) (body string, urls []string, err error) {
	var url_list []string
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		mu.Lock()
		if _, ok := (*prevFetchedUrls)[url]; !ok {
			*data = append(*data, []byte(url+"\n")...)
			(*prevFetchedUrls)[url] = struct{}{}
		}
		mu.Unlock()
		link := e.Attr("href")
		url_list = append(url_list, link)
	})
	c.Visit(url)
	return "", url_list, nil
}
