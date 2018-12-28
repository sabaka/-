package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sabaka/fileHelper"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	fmt.Println("Execution Started")
	pathToFile := getPathToFIle()

	fileHelper.DoOnEachLine(pathToFile, checkLink)

	fmt.Println("Execution finished")
	// Wait for any key to prevent console closure
	bufio.NewScanner(os.Stdin).Scan()
}

/*
	Main action
 */
func checkLink(link string) {
	findAllLinks(addPrefix(link))
}

/*
	Finds all links (href) on page and it's children
 */
func findAllLinks(link string) {

	var visited = make(map[string]bool)
	visited[link] = true

	hostname := getHostname(link)
	var queue = findLinksInSource(make(map[string]bool), parsePage(fetchPage(link)), hostname)

	fmt.Printf("Starting processing. Current queue: %v\n", queue)

	for l := range queue {
		fmt.Printf("Checking page %s. \n", l)
		q1 := findLinksInSource(queue, parsePage(fetchPage(l)), hostname)
		for c := range q1 {
			if !visited[c] {
				fmt.Printf("New unvisited page found %s. Adding to queue\n", c)
				queue[c] = true
			}
		}
		visited[l] = true
	}

}

/*
	Parses page to tree
 */
func parsePage(page []byte) *html.Node {
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		fmt.Printf("Can't parse page")
	}
	return doc
}

/*
	Find all links on page
 */
func findLinksInSource(links map[string]bool, n *html.Node, hostname string) map[string]bool {
	if n.FirstChild != nil {
		links = findLinksInSource(links, n.FirstChild, hostname)
	}
	if n.NextSibling != nil {
		links = findLinksInSource(links, n.NextSibling, hostname)
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.Contains(a.Val, hostname) && !links[a.Val] {

				links[a.Val] = true
			}
		}
	}
	return links
}

/*
	Returns hostname. Ex: http://example.org/lala/ -> example.org
 */
func getHostname(link string) string {
	url, _ := url.ParseRequestURI(link)
	return url.Hostname()
}

/*
	Downloads page code
 */
func fetchPage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Can't fetch page %s\n", url)
	}
	body := resp.Body
	defer body.Close()
	doc, _ := ioutil.ReadAll(body)
	//resp.Request.URL.Parse("")
	return doc
}

/*
	Adds http:// prefix to link if there is no
 */
func addPrefix(link string) string {
	if  hasPrefix(link) {
		return "http://" + link
	}
	return link
}

func hasPrefix(link string) bool {
	return !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://")
}

/*
	Gets file path
 */
func getPathToFIle() string {
	if len(os.Args) > 1 {
		return os.Args[1:][0]
	} else {
		return "links.lst"
	}
}
