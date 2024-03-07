# Crawler
 From a starting URL, all links on the current page are visited and all links from those pages are visited until predefined depth reaches 0. The same address won't be visited twice and all visited URLS that responded with HTML are saved in urls.txt 

This began as a exercise to implement mutual exclusion or channels using a fake crawling function in ['A Tour of Go'](https://go.dev/tour/concurrency/10). 

### Todo: 
Make the url and depth optional argvs
