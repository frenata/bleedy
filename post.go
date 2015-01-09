package blog

import (
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

const (
	titlePre   string = "Title: "
	authorPre  string = "Author: "
	tagPre     string = "Tag: "
	datePre    string = "Date: "
	dateFormat string = "2 January 2006 @ 3:04pm"
	//dateFormat string = "January 2, 2006"
	//dateFormat string = time.RFC822Z
)

type Post struct {
	title  string
	author string
	tag    string
	body   string
	date   time.Time
	//comments Comments
}

func NewPost(raw []byte, date time.Time) *Post {
	content := string(raw)

	c := strings.SplitN(content, "Body:", 2)
	if len(c) < 2 {
		return nil
	}
	meta := c[0]
	b := []byte(c[1])
	body := string(blackfriday.MarkdownCommon(b))

	m := strings.SplitN(meta, "\n", 4)
	if len(m) < 3 {
		return nil
	}
	title := strings.TrimPrefix(m[0], titlePre)
	author := strings.TrimPrefix(m[1], authorPre)
	tag := strings.TrimPrefix(m[2], tagPre)
	if len(m) > 3 {
		s := strings.TrimPrefix(m[3], datePre)
		s = strings.TrimSpace(s)
		d, err := time.Parse(dateFormat, s)
		if err == nil {
			date = d
		} else {
			//fmt.Println(d)
			//fmt.Println(err)
		}
	}

	//comments := nil //to implement
	return &Post{title, author, tag, body, date}
}

func (p *Post) String() string {
	t := "Title: " + p.title + "\n"
	a := "Author: " + p.author + "\n"
	g := "Tag: " + p.tag + "\n"
	d := "Date: " + p.date.Format(dateFormat) + "\n"

	return t + a + g + d + "\n" + p.body
}
