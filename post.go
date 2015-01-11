package blog

import (
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

const (
	titlePre    string = "Title: "
	authorPre   string = "Author: "
	tagPre      string = "Tag: "
	datePre     string = "Date: "
	templatePre string = "Template: "
	dateFormat  string = "2 January 2006 @ 3:04pm"
	//dateFormat string = "January 2, 2006"
	//dateFormat string = time.RFC822Z
)

type Post struct {
	title      string
	author     string
	tag        string
	body       string
	template   string
	date       time.Time
	dateFormat string
	//comments Comments
}

func NewPost(raw []byte, date time.Time) *Post {
	p := new(Post)
	content := string(raw)

	c := strings.SplitN(content, "Body:", 2)
	if len(c) < 2 {
		return nil
	}
	//m := c[0]
	//b := []byte(c[1])
	// TODO: does this need validation / error checking?
	p.body = string(blackfriday.MarkdownCommon([]byte(c[1])))

	meta := strings.Split(c[0], "\n")
	/* if len(m) < 3 {
		return nil
	}
	// TODO: should probably validate that these are tagged properly
	// beter yet, iterate through everything in the meta section and assign iff those variables
	// that validate
	title := strings.TrimPrefix(m[0], titlePre)
	author := strings.TrimPrefix(m[1], authorPre)
	tag := strings.TrimPrefix(m[2], tagPre)
	if len(m) > 3 {
		// TODO validate this one too
		s := strings.TrimPrefix(m[3], datePre)
		s = strings.TrimSpace(s)
		d, err := time.Parse(dateFormat, s)
		if err == nil {
			date = d
		} else {
			//fmt.Println(d)
			//fmt.Println(err)
		}
	} */
	p.dateFormat = dateFormat

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

	//comments := nil //to implement
	return p
}

func (p *Post) setDate(s string) {
	d, err := time.Parse(p.dateFormat, s)
	if err != nil {
		return
	}
	p.date = d
}

func (p *Post) validateMeta(m, pre string) (ok bool, output string) {
	if strings.HasPrefix(m, pre) {
		ok, output = true, strings.TrimSpace(strings.TrimPrefix(m, pre))
	} else {
		ok, output = false, ""
	}

	return ok, output
}

func (p *Post) String() string {
	t := "Title: " + p.title + "\n"
	a := "Author: " + p.author + "\n"
	g := "Tag: " + p.tag + "\n"
	d := "Date: " + p.date.Format(dateFormat) + "\n"

	return t + a + g + d + "\n" + p.body
}

func (p *Post) Template() string {
	return p.template
}

func (p *Post) Format(template string) ([]byte, error) {
	//bytes := []byte(fmt.Sprint(p))

	return bytes, nil
}
