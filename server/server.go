package main

import (
	"log"
	"net/http"
	"time"

	"github.com/frenata/blog"
)

const postsDir string = "server/posts"
const postsExt string = ".md"
const staticDir string = "server/static"
const staticExt string = ".html"

func check() {
	b := blog.NewBlog(postsDir, postsExt, staticDir, staticExt)
	for {
		time.Sleep(time.Second * 4)
		b.Update()
	}
}

func main() {
	go check()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(staticDir))))
}
