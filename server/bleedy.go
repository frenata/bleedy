package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/frenata/bleedy"
)

const (
	postConfig string = "post.config"
	blogConfig string = "blog.config"
)

func readConfig(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func check(b *bleedy.Blog) {
	for {
		b.UpdateScan()
		time.Sleep(time.Second * 10)
	}
}

func main() {
	var err error
	pc, bc := []string{}, []string{}
	b, p := &bleedy.Blog{}, &bleedy.PostFormatter{}
	l := log.New(os.Stdout, "bleedy server: ", log.Ltime)

	if pc, err = readConfig(postConfig); err != nil {
		l.Println(err)
	} else if p, err = bleedy.NewPostFormatter(pc); err != nil {
		l.Println(err)
	} else if bc, err = readConfig(blogConfig); err != nil {
		l.Println(err)
	} else if b, err = bleedy.NewBlog(bc, l); err != nil {
		l.Println(err)
	}

	if err != nil {
		l.Fatal("Fatal: failed to create blog properly")
	}
	b.SetFormatter(p)

	go check(b)
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(b.Output()))))
}
