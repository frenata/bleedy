package main

import (
	"fmt"
	"strings"
	"time"
)

const dateFormat string = "2 January 2006 @ 3:04pm"

func main() {
	t := time.Now()
	s := fmt.Sprint(t.Format(dateFormat))
	fmt.Println(s)
	s2 := "Prefix: 9 January 2015 @ 10:30pm"
	s2 = strings.TrimPrefix(s2, "Prefix: ")
	t2, err := time.Parse(dateFormat, s2)
	fmt.Println(t2, err)
}
