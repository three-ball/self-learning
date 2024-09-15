package main

import (
	"fmt"
	"log"
)

type Context struct {
	url     string
	content string
	data    any
}

type Handler func(*Context) error

func CheckingUrlHandler(ctx *Context) error {
	fmt.Printf("Checking URL: %s\n", ctx.url)
	return nil
}

func FetchContentHandler(c *Context) error {
	fmt.Printf("Fetching content from url: %s\n", c.url)
	c.content = "Hello, World!"
	return nil
}

func ExtractDataHandler(c *Context) error {
	fmt.Printf("Extracting data from content: %s\n", c.content)
	c.data = "Hello, World!"
	return nil
}

func SaveDataHandler(c *Context) error {
	fmt.Printf("Saving data: %v\n", c.data)
	return nil
}

type HandlerNode struct {
	handler Handler
	next    *HandlerNode
}

func (node *HandlerNode) Execute(url string) error {
	ctx := Context{
		url: url,
	}

	if node == nil || node.handler == nil {
		return nil
	}

	if err := node.handler(&ctx); err != nil {
		return err
	}

	return node.next.Execute(url)

	// nextNode := node.next
	// for nextNode != nil {
	// 	if err := nextNode.handler(&ctx); err != nil {
	// 		return err
	// 	}

	// 	nextNode = nextNode.next
	// }

	// return nil
}

func NewCrawler(handler ...Handler) HandlerNode {
	node := HandlerNode{}

	if len(handler) > 0 {
		node.handler = handler[0]
	}

	currentNode := &node

	for i := 1; i < len(handler); i++ {
		currentNode.next = &HandlerNode{handler: handler[i]}
		currentNode = currentNode.next
	}

	return node
}

type WebCrawler struct {
	handler HandlerNode
}

func (wc WebCrawler) Crawl(url string) {
	if err := wc.handler.Execute(url); err != nil {
		log.Println(err)
	}
}

func main() {
	WebCrawler{
		handler: NewCrawler(
			CheckingUrlHandler,
			FetchContentHandler,
			ExtractDataHandler,
			SaveDataHandler,
		),
	}.Crawl("https://www.google.com")
}
