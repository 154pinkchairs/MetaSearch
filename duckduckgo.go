package main

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type duckduckgo struct{}

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
