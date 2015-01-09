package main

import (
	"log"
	"net/http"
	"time"

	"github.com/frenata/blog"
)

const dir string = "static"
const postsDir string = "server/posts"

func check() {
	p := blog.NewPosts(postsDir)
	for {
		time.Sleep(time.Second * 4)
		p.CheckPosts()
	}
}

func main() {
	go check()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(dir))))
}
