package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
)


func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		container := n.Parent.Parent.Parent
		for _, atr := range container.Attr{
			if atr.Key == "class" && atr.Val == "cell" {

				for _, atr := range n.Attr{
					if atr.Key == "href" {
						links = append(links, atr.Val)
					}
				}

			}
		}

	}

	for c:= n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

func parse(url string)  {
	links, err := findLinks(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
	}

	writeUrls(links, "urls2")
}

func writeUrls(urls []string, fileName string)  {

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Unable to create file:", err)
	}

	defer file.Close()

	for _, url := range urls{
		file.WriteString(url + "\n")
	}

}



func index(c *gin.Context)  {

	file, err := os.Open("urls")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	urls := []string{}
	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	json, _ := json.Marshal(urls)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"urls": string(json),
	})


}

type resItem struct {
	counter int
	visited bool
}

func getAnswer(c *gin.Context)  {
	data := c.PostForm("data")
	response := make(map[string]map[string]int)
	json.Unmarshal([]byte(data), &response)

}

func main() {
	parse("https://hosting-pulse.ru/top-sites")

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*.html")
	//router.Static("/assets", "./assets")
	router.GET("/", index)
	router.POST("/result", getAnswer)
	router.Run(":8080")
}


