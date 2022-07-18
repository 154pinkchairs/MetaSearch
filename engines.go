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
type bing struct{}

func getRequest(url string, values url.Values, header http.Header) (*html.Node, error) {
	req, err := http.NewRequest("GET", url + values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = header

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func rankScore(rank int) float64 {
	return 1.0 / (60.0 + float64(rank))
}

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

	nodeCh := make(chan *html.Node, 1)

	go genClassShallow(root, "result", nodeCh)

	rank := 0
	for resultNode := range nodeCh {
		var r result

		linkNode := getClassFirst(resultNode, "result__a")

		r.Link = getAttribute(linkNode, "href")
		r.Title = getText(linkNode)

		snippetNode := getClassFirst(resultNode, "result__snippet")

		if snippetNode != nil {
			r.Description = getText(snippetNode)
		}

		r.SearchEngines = []string{"duckduckgo"}
		r.score = rankScore(rank)

		resultCh <- r

		rank += 1
	}
}

func (_ google) search(q string, resultCh chan <- result) {
	root, err := getRequest("https://www.google.com/search?",
		url.Values{
			"q": {q},
			"gbv": {"1"},
		},
		http.Header{
			"User-Agent": {userAgent},
			"Host": {"www.google.com"},
		})
	if err != nil {
		return
	}

	nodeCh := make(chan *html.Node, 1)

	go genWithAttributeValueShallow(root, "class", "Gx5Zad fP1Qef xpd EtOod pkphOe", nodeCh)

	rank := 0
	for resultNode := range nodeCh {
		var r result

		headerNode := getWithAttributeValueFirst(resultNode, "class", "egMi0 kCrYT")

		linkNode := getElementFirst(headerNode, "a")

		hrefQuery := getAttribute(linkNode, "href")

		hrefVals, err := url.ParseQuery(hrefQuery)
		if err != nil {
			continue
		}

		r.Link = hrefVals.Get("/url?q")

		titleNode := getWithAttributeValueFirst(headerNode, "class", "BNeawe vvjwJb AP7Wnd")

		r.Title = getText(titleNode)

		descriptionNode := getWithAttributeValueFirst(resultNode, "class", "BNeawe s3v9rd AP7Wnd")

		if descriptionNode != nil {
			r.Description = getText(descriptionNode)
		}

		r.SearchEngines = []string{"google"}
		r.score = rankScore(rank)

		resultCh <- r

		rank += 1
	}
}

func (_ bing) search(q string, resultCh chan <- result) {
	root, err := getRequest("https://www.bing.com/search?",
		url.Values{
			"q": {q},
		},
		http.Header{
			"User-Agent": {userAgent},
			"Host": {"www.bing.com"},
		})
	if err != nil {
		return
	}

	resultsNode := getId(root, "b_results")

	if resultsNode == nil {
		return
	}

	nodeCh := make(chan *html.Node, 1)

	go genDirectChildrenWithAttribute(resultsNode, "class", "b_algo", nodeCh)

	rank := 0
	for resultNode := range nodeCh {
		var r result

		/*
		headerNode := getElementFirst(resultNode, "a")
		r.Link = getAttribute(headerNode, "href")
		r.Title = getText(headerNode)

		descriptionNode := getWithAttributeValueFirst(resultNode, "class", "lineclamp4")
		if descriptionNode != nil {
			r.Description = getText(descriptionNode)
		}
		*/

		titleNode := getElementFirst(resultNode, "h2")
		r.Title = getText(titleNode)

		linkNode := getWithAttributeValueFirst(resultNode, "class", "b_attribution")
		r.Link = getText(linkNode)

		descriptionNode := getElementFirst(resultNode, "p")
		if descriptionNode != nil {
			r.Description = getText(descriptionNode)
		}

		r.SearchEngines = []string{"bing"}

		r.score = rankScore(rank)

		resultCh <- r

		rank += 1
	}
}
