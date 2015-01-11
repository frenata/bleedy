package blog

import (
	"fmt"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

// TODO: These should either be part of the Post struct or live in a config file.
const (
	titlePre    string = "Title: "
	authorPre   string = "Author: "
	tagPre      string = "Tag: "
	datePre     string = "Date: "
	templatePre string = "Template: "
	dateFormat  string = "2 January 2006 @ 3:04pm"
)

// Post holds the content of a post, parsed from a file, metadata and body content.
type Post struct {
	title    string
	author   string
	tag      string
	body     string
	template string
	date     time.Time
	//comments Comments // TODO: bool for disqus comments?
}

// NewPost takes a byte slice (from a markdowned text file), a date, and creates a new Post object.
// The date should typically be the last modification time of the file, and will be overwritten
// if a date is manually set in the Post metadata.
func NewPost(raw []byte, date time.Time) *Post {
	p := new(Post)
	content := string(raw)

	c := strings.SplitN(content, "Body:", 2)
	if len(c) < 2 {
		return nil
	}

	// TODO: does this need validation / error checking?
	p.body = string(blackfriday.MarkdownCommon([]byte(c[1])))
	meta := strings.Split(c[0], "\n")

	// new version of meta yanking: just grab each line, check if it validates any Meta tags, and assign it properly
	for _, m := range meta {
		if ok, out := p.validateMeta(m, titlePre); ok {
			p.title = out
		} else if ok, out := p.validateMeta(m, authorPre); ok {
			p.author = out
		} else if ok, out := p.validateMeta(m, tagPre); ok {
			p.tag = out
		} else if ok, out := p.validateMeta(m, templatePre); ok {
			p.template = out
		} else if ok, out := p.validateMeta(m, datePre); ok {
			p.setDate(out)
		}

	}

	return p
}

// parses the string against the specified dateformat, if it validates, manually set the post date
func (p *Post) setDate(s string) {
	d, err := time.Parse(dateFormat, s)
	if err != nil {
		return
	}
	p.date = d
}

// checks the string for valid metadata, as defined in the constant prefixes, and returns the data.
func (p *Post) validateMeta(m, pre string) (ok bool, output string) {
	if strings.HasPrefix(m, pre) {
		ok, output = true, strings.TrimSpace(strings.TrimPrefix(m, pre))
	} else {
		ok, output = false, ""
	}

	return ok, output
}

// String prints the Post meta content and body.
func (p *Post) String() string {
	t := "Title: " + p.title + "\n"
	a := "Author: " + p.author + "\n"
	g := "Tag: " + p.tag + "\n"
	d := "Date: " + p.date.Format(dateFormat) + "\n"

	return t + a + g + d + "\n" + p.body
}

// Template returns the Post specific template. If one was not set in the metadata, this should return ""
func (p *Post) Template() string {
	return p.template
}

// Format takes a template file and creates a []byte representing an html document populated with the Post content,
// ready for writing to a file.
func (p *Post) Format(template string) ([]byte, error) {
	bytes := []byte(fmt.Sprint(p))

	return bytes, nil
}
