package main

import (
	"strings"

	"golang.org/x/net/html"
)

func hasClass(n *html.Node, class string) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classes := strings.Split(attr.Val, " ")

				for _, nclass := range classes {
					if nclass == class {
						return true
					}
				}
			}
		}
	}

	return false
}

func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}

func getClassShallow(n *html.Node, class string) []*html.Node {
	if hasClass(n, class) {
		return []*html.Node{n}
	}
	
	nodes := make([]*html.Node, 0)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		newnodes := getClassShallow(c, class)

		nodes = append(nodes, newnodes...)
	}

	return nodes
}

func getClassFirst(n *html.Node, class string) *html.Node {
	if hasClass(n, class) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		o := getClassFirst(c, class)

		if o != nil {
			return o
		}
	}

	return nil
}

// If val is "" then matches all nodes
func getWithAttributeValueShallow(n *html.Node, attr string, val string) []*html.Node {
	if n.Type == html.ElementNode && getAttribute(n, attr) == val {
		return []*html.Node{n}
	}

	nodes := make([]*html.Node, 0)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		newnodes := getWithAttributeValueShallow(c, attr, val)

		nodes = append(nodes, newnodes...)
	}

	return nodes
}

func getWithAttributeValueFirst(n *html.Node, attr string, val string) *html.Node {
	if n.Type == html.ElementNode && getAttribute(n, attr) == val {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		o := getWithAttributeValueFirst(c, attr, val)

		if o != nil {
			return o
		}
	}

	return nil
}

func getElementFirst(n *html.Node, elem string) *html.Node {
	if n.Type == html.ElementNode && n.Data == elem {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		o := getElementFirst(c, elem)

		if o != nil {
			return o
		}
	}

	return nil
}

func getTextBuilder(n *html.Node, b *strings.Builder) {
	if n.Type == html.TextNode {
		b.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getTextBuilder(c, b)
	}
}

func getText(n *html.Node) string {
	var b strings.Builder

	getTextBuilder(n, &b)

	return b.String()
}
