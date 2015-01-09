package blog

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Blog struct {
	rDir string               // read directory
	rExt string               // read extension
	wDir string               // write directory
	wExt string               // write extension
	hash map[string]time.Time // hash map to check for updates
}

func NewBlog(rDir, rExt, wDir, wExt string) *Blog {
	b := &Blog{rDir, rExt, wDir, wExt, make(map[string]time.Time)}
	return b
}

func (b *Blog) readFile(file string, date time.Time) error {
	content, err := ioutil.ReadFile(path.Join(b.rDir, file) + b.rExt)
	if err != nil {
		return err
	}
	// TODO: call a parse function, read to struct?
	//ioutil.WriteFile(path.Join(p.staticDir, file)+p.staticExt, output, 0666)
	//fmt.Println(string(content))
	p := NewPost(content, date)
	fmt.Println(p)
	return nil
}

func (b *Blog) writeFile(p Post) error {

	return nil
}

// Update checks the read directory for changes to files
func (b *Blog) Update() {
	files, err := ioutil.ReadDir(b.rDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		n := f.Name()
		switch {
		case f.IsDir():
			continue // TODO: descend into directory, auto-tag with name
		case strings.HasSuffix(n, b.rExt):
			n := strings.TrimSuffix(n, b.rExt)
			if _, ok := b.hash[n]; ok {
				fmt.Printf("%v was last updated %v ago\n", n, time.Since(b.hash[n]))
				if b.hash[n] == f.ModTime() {
					continue
				}
			}
			b.hash[n] = f.ModTime()
			fmt.Printf("Update %v\n", n)
			b.readFile(n, f.ModTime())
		}

	}
}
