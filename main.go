package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	maxDepth    int
	domainOnly  bool
	visited     = make(map[string]bool)
	visitedLock sync.Mutex
	baseHost    string
	wg          sync.WaitGroup
)

func crawl(link string, depth int) {
	defer wg.Done()

	if depth > maxDepth {
		return
	}

	visitedLock.Lock()
	if visited[link] {
		visitedLock.Unlock()
		return
	}
	visited[link] = true
	visitedLock.Unlock()

	resp, err := http.Get(link)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Error fetching site: %v\n", err)
		return
	}
	defer resp.Body.Close()

	htmlDoc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing html: %v\n", err)
		return
	}

	var links []string
	walkDocAndCollectLinks(htmlDoc, link, &links)

	for _, l := range links {
		wg.Add(1)
		go crawl(l, depth+1)
	}
}

func walkDocAndCollectLinks(n *html.Node, base string, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				u, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}

				abs := resolveURL(u, base)
				if abs != "" && shouldVisit(abs) {
					*links = append(*links, abs)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkDocAndCollectLinks(c, base, links)
	}
}

func resolveURL(u *url.URL, base string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	resolved := baseURL.ResolveReference(u)

	return trimTrailingSlash(resolved.String())
}

func shouldVisit(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	if domainOnly && parsed.Host != baseHost {
		return false
	}

	return strings.HasPrefix(parsed.Scheme, "http")
}

func trimTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	return url
}

func writeToFile() {
	f, err := os.Create("sitemap.txt")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	visitedLock.Lock()
	defer visitedLock.Unlock()

	var links []string
	for link := range visited {
		links = append(links, link)
	}

	sort.Strings(links)

	for _, link := range links {
		f.WriteString(link + "\n")
	}

	log.Println("Sitemap written to sitemap.txt")
}

func main() {
	flag.IntVar(&maxDepth, "depth", 3, "max crawl depth")
	flag.BoolVar(&domainOnly, "domain-only", true, "only crawl same domain")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("Usage: web-scraper <url>")
	}

	siteUrl := flag.Arg(0)
	u, err := url.Parse(siteUrl)
	if err != nil {
		log.Fatalf("Invalid URL provided: %v", err)
	}

	baseHost = u.Host

	wg.Add(1)
	go crawl(siteUrl, 0)
	wg.Wait()

	writeToFile()
}
