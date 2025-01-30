package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

var raw = `
<!DOCTYPE html>
<html>
  <body>
    <h1>My First Heading</h1>
    <p>My first paragraph.</p>
    <p>HTML images are defined with the img tag:</p>
    <img src="xxx.jpg" width="104" height="142">
  </body>
</html>
`

func main() {
	doc, err := html.Parse(bytes.NewReader([]byte(raw)))

	if err != nil {
		fmt.Fprintf(os.Stderr, "parse failed: %s\n", err)
		os.Exit(-1)
	}

	words, pics := countWordsAndImages(doc)

	fmt.Printf("%d words and %d images\n", words, pics)
}

func countWordsAndImages(doc *html.Node) (int, int) {
	var words, pics int
	visit(doc, &words, &pics)
	return words, pics
}

func visit(node *html.Node, words, pics *int) {
	if node.Type == html.TextNode {
		*words += len(strings.Fields(node.Data))
	} else if node.Type == html.ElementNode && node.Data == "img" {
		*pics++
	}

	for c := range node.ChildNodes() {
		visit(c, words, pics)
	}
}
