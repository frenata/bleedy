package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/frenata/bleedy"
)

const (
	postsDir    string = "posts"
	postsExt    string = ".md"
	staticDir   string = "static"
	staticExt   string = ".html"
	templateDir string = "templates"
	templateExt string = ".html"
	postConfig  string = "post.config"
	blogConfig  string = "blog.config"
)

func readConfig(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func check() {
	pc, _ := readConfig(postConfig)
	p, _ := blog.NewPostFormatter(pc)
	_, _ = readConfig(blogConfig)
	b := blog.NewBlog(postsDir, postsExt, staticDir, staticExt, templateDir, templateExt)
	b.SetFormatter(p)
	for {
		time.Sleep(time.Second * 4)
		b.UpdateScan()
	}
}

func main() {
	go check()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(staticDir))))
}
