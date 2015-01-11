package blog

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Blog struct {
	input    files                //  read directory/ext
	output   files                //  write directory/ext
	template files                // template direcoty/ext
	hash     map[string]time.Time // hash map to check for updates
}

type files struct {
	dir string
	ext string
	def string
}

func NewBlog(iDir, iExt, oDir, oExt, tDir, tExt string) *Blog {
	b := &Blog{hash: make(map[string]time.Time)}
	b.SetInput(iDir, iExt)
	b.SetOutput(oDir, oExt)
	b.SetTemplate(tDir, tExt)
	b.SetDefaultTemplate("")
	return b
}

// SetInput sets the input directory and file extension, this must be set when calling New, so not likely anyone will
// use this unless changing directories on the fly is somehow required, included for completeness though
func (b *Blog) SetInput(dir, ext string) {
	b.input.dir = dir
	b.input.ext = ext
}

func (b *Blog) SetOutput(dir, ext string) {
	b.output.dir = dir
	b.output.ext = ext
}

func (b *Blog) SetTemplate(dir, ext string) {
	b.template.dir = dir
	b.template.ext = ext
}

func (b *Blog) SetDefaultTemplate(def string) {
	if def != "" {
		b.template.def = def
	} else {
		b.template.def = "default"
	}
}

func (b *Blog) readFile(file string, date time.Time) (*Post, error) {
	content, err := ioutil.ReadFile(path.Join(b.input.dir, file) + b.input.ext)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// TODO: call a parse function, read to struct?
	//ioutil.WriteFile(path.Join(p.staticDir, file)+p.staticExt, output, 0666)
	//fmt.Println(string(content))
	p := NewPost(content, date)
	//fmt.Println(p)
	return p, nil
}

func (b *Blog) writeFile(file string, p *Post) error {
	//quick and dirty, replacing with templating
	name := ""
	if name = p.Template(); name == "" {
		name = b.template.def
	}
	template := path.Join(b.template.dir, name) + b.template.ext

	output, err := p.Format(template)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(b.output.dir, file)+b.output.ext, output, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Update checks the read directory for changes to files
func (b *Blog) Update() {
	files, err := ioutil.ReadDir(b.input.dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		n := f.Name()
		switch {
		case f.IsDir():
			continue // TODO: descend into directory, auto-tag with name
		case strings.HasSuffix(n, b.input.ext):
			n := strings.TrimSuffix(n, b.input.ext)
			if _, ok := b.hash[n]; ok {
				fmt.Printf("%v was last updated %v ago\n", n, time.Since(b.hash[n]))
				if b.hash[n] == f.ModTime() {
					continue
				}
			}
			b.hash[n] = f.ModTime()
			fmt.Printf("Update %v\n", n)
			p, err := b.readFile(n, f.ModTime())
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println(p)a
			err = b.writeFile(n, p)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
