package main

import (
	"log"
	"net/http"
	"time"

	"github.com/frenata/blog"
)

const (
	postsDir        string = "server/posts"
	postsExt        string = ".md"
	staticDir       string = "server/static"
	staticExt       string = ".html"
	templateDir     string = "server/templates"
	templateExt     string = ".t"
	defaultTemplate string = "blog"
)

func check() {
	b := blog.NewBlog(postsDir, postsExt, staticDir, staticExt, templateDir, templateExt)
	for {
		time.Sleep(time.Second * 4)
		b.Update()
	}
}

func main() {
	go check()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(staticDir))))
}
