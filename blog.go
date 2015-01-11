/*
blog handles creation of a Blog that will scan an input directory for new/modified files (markdown for instance), and
parse the metadata and content of those files (content with github.com/russross/blackfriday) and create files of the same
name in html format in a designated output directory.

NewBlog creates a new blog, SetInput/Output/Template allow finetuned or changing control of the directory/formats Blog scans for.

The primary method is Blog.Update(), which scans for the new/modified files, checking their last modification date against an
internal map. Changes trigger calls to read the file, create a new Post struct (see post.go), format it, and write it to the
output.
*/

package blog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

// Blog holds the definitions of file directories/extensions for inputs, outputs, and templates, as well as the hashmap
// for checking time of last modification of files in the input directory.
type Blog struct {
	input    files                //  read directory/ext
	output   files                //  write directory/ext
	template files                // template directory/ext
	hash     map[string]time.Time // hash map to check for updates
}

// small struct to make Blog look prettier, defines a location and type of file, as well as a default filename.
type files struct {
	dir string
	ext string
	def string
}

// NewBlog creates a new Blog object, populated with all the directories/extensions for input/ouput/templates.
// It also allocates the hashmap for checking file modification times.
func NewBlog(iDir, iExt, oDir, oExt, tDir, tExt string) *Blog {
	b := &Blog{hash: make(map[string]time.Time)}
	b.SetInput(iDir, iExt)
	b.SetOutput(oDir, oExt)
	b.SetTemplate(tDir, tExt)
	return b
}

// SetInput sets the directory and file extension for markeddown text files.
func (b *Blog) SetInput(dir, ext string) {
	b.input.dir = dir
	b.input.ext = ext
}

// SetOutput sets the directory and file extension for generated html files.
func (b *Blog) SetOutput(dir, ext string) {
	b.output.dir = dir
	b.output.ext = ext
}

// SetTemplate sets the directory and file extension for template files.
func (b *Blog) SetTemplate(dir, ext string) {
	b.template.dir = dir
	b.template.ext = ext
	b.SetDefaultTemplate("")
}

// SetDefaultTemplate sets the filename for the default post template. This should be in the template directory.
func (b *Blog) SetDefaultTemplate(def string) {
	if def != "" {
		b.template.def = def
	} else {
		b.template.def = "default"
	}
}

// reads from the specified input file (markdown), creates a new Post, and returns it.
func (b *Blog) readFile(file string, date time.Time) (*Post, error) {
	name := path.Join(b.input.dir, file) + b.input.ext
	content, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	p, err := NewPost(content, date)
	if err != nil {
		return nil, errors.New(fmt.Sprint(err) + name)
	}

	return p, nil
}

// formats and writes the content of a Post to the specified file
func (b *Blog) writeFile(file string, p *Post) error {
	name := ""
	if p.Template == "" {
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

// Update checks the read directory for changes to files. If it detects changes (based on last-modified time),
// it reads the input file and creates an output file of the same name.
// Designed to be called continously in a loop.
func (b *Blog) Update() {
	// read all the files in the input directory
	files, err := ioutil.ReadDir(b.input.dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// check each file
	for _, f := range files {
		n := f.Name()
		switch {
		case f.IsDir(): // TODO: descend into directory, auto-tag with name
			continue
		case strings.HasSuffix(n, b.input.ext): // check the suffix
			n := strings.TrimSuffix(n, b.input.ext) // remove the suffix
			if _, ok := b.hash[n]; ok {             // is it already in the hashmap?
				//fmt.Printf("%v was last updated %v ago\n", n, time.Since(b.hash[n])) // TODO: log not print
				if b.hash[n] == f.ModTime() {
					continue // file has not been modified since the last check, ignore it
				}
			}
			b.hash[n] = f.ModTime()                 // store the last modified time
			fmt.Printf("Update %v\n", n)            // TODO: log not print
			post, err := b.readFile(n, f.ModTime()) // read the file, creating a post
			if err != nil {
				fmt.Println(err) // TODO: log not print
				return
			}

			err = b.writeFile(n, post) // write the post to a file
			if err != nil {
				fmt.Println(err) // TODO: log not print
				return
			}
		}

	}
}
