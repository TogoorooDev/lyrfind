
package main
import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"bytes"
	
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func collectText(n *html.Node, buf *bytes.Buffer) {
    if n.Type == html.TextNode {
        buf.WriteString(n.Data)
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        collectText(c, buf)
    }
}

func main(){
	//var artist string
	var songname string
	if  len(os.Args) == 1 {
		fmt.Println("Usage: lyrfind name [artist]")
		return
	}
	if len(os.Args) == 2 {
		songname = os.Args[1]
		//artist = ""
	}

	params := url.Values{}
	params.Add("q", songname)
	ops, err := http.Get("http://search.azlyrics.com/search.php?" + params.Encode())
	if err != nil {
		fmt.Println("Fatal, quitting")
		fmt.Println(err)
		return
	}
	defer ops.Body.Close()
	doc, err := goquery.NewDocumentFromReader(ops.Body)
	if err != nil {
		fmt.Println("Fatal, quitting")
		fmt.Println(err)
		return
	}

	lyrobjs := doc.Find(".visitedlyr")
	attr := lyrobjs.Find("a").AttrOr("href", "fatal")
	title := &bytes.Buffer{}
	collectText(lyrobjs.Get(0), title)
	fmt.Println(title)
	
	lyrspage, err := http.Get(attr)
	if err != nil {
		fmt.Println("Fatal, quitting")
		fmt.Println(err)
		return
	}
	defer lyrspage.Body.Close()
	lyrsdoc, err := goquery.NewDocumentFromReader(lyrspage.Body)
	if err != nil {
		fmt.Println("Fatal, quitting")
		fmt.Println(err)
		return
	}
	lyrsdoc.Find("div").Each(func(i int, s *goquery.Selection){
		for _, node := range s.Nodes {
			if len(node.Attr) == 0 {
				text := &bytes.Buffer{}
				collectText(node, text)
				fmt.Println(text)
			}
		}
	})
}

