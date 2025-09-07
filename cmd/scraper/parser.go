package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	Text     string
	Children []*Node
}

// printTree recursively prints a Node tree with indentation
func printNewTree(n *Node, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s- %s\n", indent, n.Text)
	for _, child := range n.Children {
		printNewTree(child, depth+1)
	}
}

// printTree shows the parsed HTML tree in a human-friendly structure
func printTree(n *html.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	switch n.Type {
	case html.ElementNode:
		if len(n.Attr) > 0 {
			fmt.Printf("%sDepth: %d, Tag: <%s> [", indent, depth, n.Data)
			for i, a := range n.Attr {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s=%q", a.Key, a.Val)
			}
			fmt.Println("]")
		} else {
			fmt.Printf("%sDepth: %d, Tag: <%s>\n", indent, depth, n.Data)
		}

	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%sDepth: %d Text: %q\n", indent, depth, text)
		}
	default:
		// remove node
	}

	// recurse
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printTree(c, depth+1)
	}
}

func buildTree(n *html.Node, parent_node *Node) {

	if n.Data == "ul" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			buildTree(c, parent_node)
		}
	} else if n.Data == "li" {
		newNode := &Node{}
		var ulNode *html.Node = nil
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				newNode.Text = strings.TrimSpace(c.Data)
			} else if c.Type == html.ElementNode && c.Data == "ul" {
				ulNode = c
			}
		}

		parent_node.Children = append(parent_node.Children, newNode)
		if ulNode != nil {
			buildTree(ulNode, newNode)
		}

	}

}

func extractBrackets(s string, tag string) string {
	// Compile regex: <(.*?)>  → non-greedy match between < and >
	re := regexp.MustCompile(`<([^>]*)>`)

	// Find all matches
	matches := re.FindAllStringSubmatch(s, -1)

	var results []string
	for _, match := range matches {
		// match[1] is the part inside <>
		if strings.Contains(match[1], tag) {
			results = append(results, "<"+match[1]+">")
			// fmt.Printf("%s\n", "<"+match[1]+">")
		}

	}

	for _, value := range results {
		s = strings.ReplaceAll(s, value, "")
	}
	return s
}

func findFirstUL(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "ul" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ul := findFirstUL(c); ul != nil {
			return ul
		}
	}
	return nil
}

// func main() {
// 	raw := `<div><div><div><ul><li><span>Complete <!-- -->all<!-- --> of the following</span><ul><div><span></span><li><span>Complete <!-- -->1<!-- --> of the following</span><ul><li data-test="ruleView-A.1.1.1"><div data-test="ruleView-A.1.1.1-result">Must have completed at least <span>1</span> of the following: <div><ul style="margin-top:5px;margin-bottom:5px"><li><span><a href="#/courses/view/66eda9af415235a3711518c9" target="_blank">MATH136</a> <!-- -->-<!-- --> <!-- -->Linear Algebra 1 for Honours Mathematics<!-- --> <span style="margin-left:5px">(0.50)</span></span></li><li><span><a href="#/courses/view/65b5329e98f0445a65101f57" target="_blank">MATH146</a> <!-- -->-<!-- --> <!-- -->Linear Algebra 1 (Advanced Level)<!-- --> <span style="margin-left:5px">(0.50)</span></span></li></ul></div></div></li><li data-test="ruleView-A.1.1.2"><div data-test="ruleView-A.1.1.2-result">Earned a minimum grade of <span>70%</span> in at least <span>1</span> of the following: <div><ul style="margin-top:5px;margin-bottom:5px"><li><span><a href="#/courses/view/65b90528f81cd8a84113167d" target="_blank">MATH106</a> <!-- -->-<!-- --> <!-- -->Applied Linear Algebra 1<!-- --> <span style="margin-left:5px">(0.50)</span></span></li><li><span><a href="#/courses/view/65b905a35f44041eaff296d1" target="_blank">MATH114</a> <!-- -->-<!-- --> <!-- -->Linear Algebra for Science<!-- --> <span style="margin-left:5px">(0.50)</span></span></li><li><span><a href="#/courses/view/65b90628b44f43033f69e76a" target="_blank">MATH115</a> <!-- -->-<!-- --> <!-- -->Linear Algebra for Engineering<!-- --> <span style="margin-left:5px">(0.50)</span></span></li></ul></div></div></li></ul></li></div><li data-test="ruleView-A"><div data-test="ruleView-A-result">Earned a minimum cumulative average of <span>60</span></div></li></ul></li></ul></div></div></div>`
// 	data := strings.ReplaceAll(raw, "<!-- -->", "")
// 	// data2 := strings.ReplaceAll(data, "<div>", "")
// 	parsed_data := extractBrackets(extractBrackets(extractBrackets(data, "href"), "div"), "span")
// 	// fmt.Printf("%s\n", parsed_data)

// 	doc, err := html.Parse(strings.NewReader(parsed_data))
// 	if err != nil {
// 		panic(err)
// 	}
// 	// printTree(doc, 0)
// 	newTree := &Node{}
// 	buildTree(findFirstUL(doc), newTree)
// 	printNewTree(newTree, 0)

// }
