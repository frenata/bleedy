package bleedy

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

// Formatter defines an interface that can Parse an input byte slice into some data, later work with a template file to
// Format that data into another byte slice, and can also pretty print itself with String() for simple representation.
type Formatter interface {
	Parse(b []byte, date time.Time) (Formatter, error)
	Format(templateFile string) ([]byte, error)
	String() string
	Date() string
}

// Post holds the content of a post, parsed from a file, metadata and body content.
type PostFormatter struct {
	Title    string
	Author   string
	Tag      string
	Template string
	Body     template.HTML //string //[]byte
	date     time.Time
	config   *configPostFormatter
	//comments Comments // TODO: bool for disqus comments?
}

// config struct to be read from a file - remove all the consts
type configPostFormatter struct {
	titlePre    string
	bodyPre     string
	authorPre   string
	tagPre      string
	datePre     string
	templatePre string
	dateFormat  string
}

// NewPostFormatter takes the contents of a configuration file and returns a correctly initialized PostFormatter
func NewPostFormatter(conf []string) (*PostFormatter, error) {
	//for now just convert the config string to the dateformat
	if len(conf) != 8 {
		return nil, errors.New("improper config file")
	}
	c := &configPostFormatter{}
	p := &PostFormatter{config: c}

	// initialize the constants
	c.titlePre = strings.TrimPrefix(conf[0], "titlePre: ")
	c.authorPre = strings.TrimPrefix(conf[1], "authorPre: ")
	c.tagPre = strings.TrimPrefix(conf[2], "tagPre: ")
	c.templatePre = strings.TrimPrefix(conf[3], "templatePre: ")
	c.datePre = strings.TrimPrefix(conf[4], "datePre: ")
	c.bodyPre = strings.TrimPrefix(conf[5], "bodyPre: ")
	c.dateFormat = strings.TrimPrefix(conf[6], "dateFormat: ")
	return p, nil
}

func newPost(c *configPostFormatter) (*PostFormatter, error) {
	p := &PostFormatter{config: c}
	return p, nil
}

// PostFormatter takes a byte slice (from a markdowned text file), a date, and creates a new Post object.
// The date should typically be the last modification time of the file, and will be overwritten
// if a date is manually set in the Post metadata.
func (p *PostFormatter) Parse(b []byte, date time.Time) (Formatter, error) {
	newP, err := newPost(p.config)
	if err != nil {
		return nil, err
	}

	content := string(b)
	c := strings.SplitN(content, newP.config.bodyPre, 2)
	if len(c) < 2 {
		return nil, errors.New("invalid post file - no content detected: ")
	}

	// TODO: does this need validation / error checking?
	bf := string(blackfriday.MarkdownCommon([]byte(c[1])))
	newP.Body = template.HTML(strings.TrimSpace(bf))

	// new version of meta yanking: just grab each line, check if it validates any Meta tags, and assign it properly
	meta := strings.Split(c[0], "\n")
	for _, m := range meta {
		if ok, out := newP.validateMeta(m, newP.config.titlePre); ok {
			newP.Title = out
		} else if ok, out := newP.validateMeta(m, newP.config.authorPre); ok {
			newP.Author = out
		} else if ok, out := newP.validateMeta(m, newP.config.tagPre); ok {
			newP.Tag = out
		} else if ok, out := newP.validateMeta(m, newP.config.templatePre); ok {
			newP.Template = out
		} else if ok, out := newP.validateMeta(m, newP.config.datePre); ok {
			newP.setDate(out)
		} else if newP.date.IsZero() {
			newP.date = date
		}
	}
	return newP, nil
}

// parses the string against the specified dateformat, if it validates, manually set the post date
func (p *PostFormatter) setDate(s string) {
	d, err := time.Parse(p.config.dateFormat, s)
	if err != nil {
		return
	}
	p.date = d
}

// checks the string for valid metadata, as defined in the constant prefixes, and returns the data.
func (p *PostFormatter) validateMeta(m, pre string) (ok bool, output string) {
	if strings.HasPrefix(m, pre) {
		ok, output = true, strings.TrimSpace(strings.TrimPrefix(m, pre))
	} else {
		ok, output = false, ""
	}

	return ok, output
}

// String prints the Post meta content and body.
func (p *PostFormatter) String() string {
	t := "Title: " + p.Title + "\n"
	a := "Author: " + p.Author + "\n"
	g := "Tag: " + p.Tag + "\n"
	d := "Date: " + p.Date() + "\n"

	return t + a + g + d + "\n" + string(p.Body)
}

// Format takes a template file and creates a []byte representing an html document populated with the Post content,
// ready for writing to a file.
func (p *PostFormatter) Format(templateFile string) ([]byte, error) {
	buf := &bytes.Buffer{} // byte buffer to use for template execution

	t, err := template.ParseFiles(templateFile) // load the template file
	if err != nil {
		return nil, err
	}

	err = t.Execute(buf, p) // add the post content into the template file
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil // write the buffer to a []byte prepped to write to a file.
}

// Format the date into configured readable format.
func (p *PostFormatter) Date() string {
	return p.date.Format(p.config.dateFormat)
}
