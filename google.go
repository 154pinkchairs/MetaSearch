package main

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type google struct{}

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
