package jail

import (
	"fmt"
	"regexp"
)

type Matcher struct {
	Lines   chan string
	Regexp  []string
	Matches chan int
}

func NewMatcher(regexp []string, lines chan string) *Matcher {
	return &Matcher{Regexp: regexp,
		Lines:   lines,
		Matches: make(chan int)}
}

func (m *Matcher) Run() {
	for {
		select {
		case str := <-m.Lines:
			for _, value := range m.Regexp {
				var re = regexp.MustCompile(value)
				if re.MatchString(str) {
					fmt.Println(str)
					m.Matches <- 1
				}
			}
		}
	}
}
