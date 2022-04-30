package main

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type searchEngine interface {
	search(string, chan <- result)
}

type result struct {
	Title string
	Link string
	Description string
	SearchEngines []string
	score float64
}

type duckduckgo struct{}

type google struct{}

func (_ duckduckgo) search(q string, resultCh chan <- result) {
	vals := url.Values(map[string][]string{
		"q": {q},
		"b": {""},
	})

	req, err := http.NewRequest("POST", "https://html.duckduckgo.com/html/", strings.NewReader(vals.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://html.duckduckgo.com")
	req.Header.Set("Referer", "https://html.duckduckgo.com")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	resultNodes := getClassShallow(root, "result")

	for i, resultNode := range resultNodes {
		var r result

		linkNode := getClassFirst(resultNode, "result__a")

		r.Link = getAttribute(linkNode, "href")
		r.Title = getText(linkNode)

		snippetNode := getClassFirst(resultNode, "result__snippet")

		if snippetNode != nil {
			r.Description = getText(snippetNode)
		}

		r.SearchEngines = []string{"duckduckgo"}
		r.score = 1 / (60.0 + float64(i))

		resultCh <- r
	}
}

func (_ google) search(q string, resultCh chan <- result) {
	vals := url.Values{
		"q": {q},
		"gbv": {"1"},
	}

	searchUrl := "https://www.google.com/search?" + vals.Encode()

	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Host", "www.google.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	resultNodes := getWithAttributeValueShallow(root, "class", "ZINbbc luh4tb xpd O9g5cc uUPGi")

	for i, resultnode := range resultNodes {
		var r result

		headerNode := getWithAttributeValueFirst(resultnode, "class", "egMi0 kCrYT")

		linkNode := getElementFirst(headerNode, "a")

		hrefQuery := getAttribute(linkNode, "href")

		hrefVals, err := url.ParseQuery(hrefQuery)
		if err != nil {
			continue
		}

		r.Link = hrefVals.Get("/url?q")

		titleNode := getWithAttributeValueFirst(headerNode, "class", "BNeawe vvjwJb AP7Wnd")

		r.Title = getText(titleNode)

		descriptionNode := getWithAttributeValueFirst(resultnode, "class", "BNeawe s3v9rd AP7Wnd")

		if descriptionNode != nil {
			r.Description = getText(descriptionNode)
		}

		r.SearchEngines = []string{"google"}
		r.score = 1 / (60.0 + float64(i))

		resultCh <- r
	}
}
