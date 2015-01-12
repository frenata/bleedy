/*
bleedy handles creation of a Blog that will scan an input directory for new/modified files (markdown for instance), and
parse the metadata and content of those files (content with github.com/russross/blackfriday) and create files of the same
name in html format in a designated output directory.

NewBlog creates a new blog, SetInput/Output/Template allow finetuned or changing control of the directory/formats Blog scans for.

The primary method is Blog.Update(), which scans for the new/modified files, checking their last modification date against an
internal map. Changes trigger calls to read the file, create a new Post struct (see post.go), format it, and write it to the
output.
*/
package bleedy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	log      *log.Logger
	Formatter
}

// small struct to make Blog look prettier, defines a location and type of file, as well as a default filename.
type files struct {
	dir string
	ext string
	def string
}

// NewBlog creates a new Blog object, populated with all the directories/extensions for input/ouput/templates.
// It also allocates the hashmap for checking file modification times.
func NewBlog(conf []string, log *log.Logger) (*Blog, error) {
	b := &Blog{hash: make(map[string]time.Time)}
	if err := b.config(conf); err != nil {
		return nil, err
	}
	b.log = log
	return b, nil
}

// SetFormatter sets the Formatter for Blog to use
func (b *Blog) SetFormatter(f Formatter) {
	b.Formatter = f
}

// Config sets the directory and file extensions for input,output,and template directories.
func (b *Blog) config(conf []string) error {
	if len(conf) != 8 {
		return errors.New("improper config file")
	}
	// initialize the constants
	b.input.dir = strings.TrimPrefix(conf[0], "inputDir: ")
	b.input.ext = strings.TrimPrefix(conf[1], "inputExt: ")
	b.output.dir = strings.TrimPrefix(conf[2], "outputDir: ")
	b.output.ext = strings.TrimPrefix(conf[3], "outputExt: ")
	b.template.dir = strings.TrimPrefix(conf[4], "templateDir: ")
	b.template.ext = strings.TrimPrefix(conf[5], "templateExt: ")
	b.setDefaultTemplate(strings.TrimPrefix(conf[6], "defaultTem: "))

	return nil
}

// SetDefaultTemplate sets the filename for the default post template. This should be in the template directory.
func (b *Blog) setDefaultTemplate(def string) {
	if def != "" {
		b.template.def = def
	} else {
		b.template.def = "default"
	}
}

// Output returns the ouput directory of the blog, for serving HTTP
func (b *Blog) Output() string {
	return b.output.dir
}

// reads from the specified input file (markdown), creates a new Post, and returns it.
func (b *Blog) readFile(file string, date time.Time) (Formatter, error) {
	name := path.Join(b.input.dir, file) + b.input.ext
	content, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	f, err := b.Formatter.Parse(content, date)
	if err != nil {
		return nil, errors.New(fmt.Sprint(err) + name)
	}

	return f, nil
}

// formats and writes the content of a Post to the specified file
func (b *Blog) writeFile(file string, f Formatter) error {
	template := path.Join(b.template.dir, b.template.def) + b.template.ext

	output, err := f.Format(template)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(b.output.dir, file)+b.output.ext, output, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Break the logic in two: func Scan() which looks for and detects new/modified posts, and func Update(filename)
// which updates a specific post. Both can be re-written to be unexported, then a simpler method ScanAndUpdate replicate
// the old functionality of update. This will allow for easier implementationg of other update mechanisms, like ones
// that include loading all the Posts into a [] and sorting them, and/or creating multipost pages.

// UpdateScan checks the read directory for changes to files. If it detects changes (based on last-modified time),
// it reads the input file and creates an output file of the same name.
// Designed to be called continously in a loop.
func (b *Blog) UpdateScan() {
	b.update(b.scan(false))
}

// scans the input directory for new/modified files, adds them to the update list. If force == true, adds all files to the list.
func (b *Blog) scan(force bool) (update []os.FileInfo) {
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
				b.log.Printf("%v was last updated %v ago\n", n, time.Since(b.hash[n]))
				if b.hash[n] == f.ModTime() && !force {
					continue // file has not been modified since the last check, ignore it
				}
			}
			update = append(update, f)
			b.hash[n] = f.ModTime() // store the last modified time
		}

	}
	return update
}

// takes a list of files, reads each, formats it, and writes it to the output
func (b *Blog) update(files []os.FileInfo) {
	for _, f := range files {
		n := strings.TrimSuffix(f.Name(), b.input.ext) // remove the suffix
		b.log.Printf("Update %v\n", n)
		post, err := b.readFile(n, f.ModTime()) // read the file, creating a post
		if err != nil {
			b.log.Println(err, "Will continue trying.")
			return
		}

		err = b.writeFile(n, post) // write the post to a file
		if err != nil {
			b.log.Println(err, "Will continue trying.")
			return
		}
	}
}
