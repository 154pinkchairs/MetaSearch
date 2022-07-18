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

func hasId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" {
				ids := strings.Split(attr.Val, " ")

				for _, nid := range ids {
					if nid == id {
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

func genClassShallowWorker(n *html.Node, class string, ch chan <- *html.Node)  {
	if hasClass(n, class) {
		ch <- n
	}
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		genClassShallowWorker(c, class, ch)
	}
}

func genClassShallow(n *html.Node, class string, ch chan <- *html.Node) {
	genClassShallowWorker(n, class, ch)
	close(ch)
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

func getId(n *html.Node, id string) *html.Node {
	if hasId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		o := getId(c, id)

		if o != nil {
			return o
		}
	}

	return nil
}

// If val is "" then matches all nodes
func genWithAttributeValueShallowWorker(n *html.Node, attr string, val string, ch chan <- *html.Node) {
	if n.Type == html.ElementNode && getAttribute(n, attr) == val {
		ch <- n
	} else if n.Type == html.ElementNode && getAttribute(n, attr) != "" {
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		genWithAttributeValueShallowWorker(c, attr, val, ch)
	}
}

func genWithAttributeValueShallow(n *html.Node, attr string, val string, ch chan <- *html.Node) {
	genWithAttributeValueShallowWorker(n, attr, val, ch)
	close(ch)
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

func genDirectChildrenWithAttribute(n *html.Node, attr string, val string, ch chan <- *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && getAttribute(c, attr) == val {
			ch <- c
		}
	}
	close(ch)
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
