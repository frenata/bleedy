package blog

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Posts struct {
	dir  string
	hash map[string]time.Time
}

func NewPosts(dir string) *Posts {
	p := &Posts{dir, nil}
	return p
}

func (p *Posts) writePost(file string) error {
	b, e := ioutil.ReadFile(file)
	fmt.Println(string(b))
	return e
}

func (p *Posts) scanFiles() (files []string) {

	return files
}

func (p *Posts) CheckPosts() {
	files, err := ioutil.ReadDir(p.dir)
	os.Chdir(p.dir)
	d, err := os.Open(".")
	files, err := d.Readdirnames(-1)

	fmt.Println(files)
	s := "test.txt"
	t, err := os.Stat(s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Last modified %v\n", t.ModTime())
	}
	p.writePost(s)
}
