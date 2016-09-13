package go2ban

import (
	"fmt"

	"github.com/karasz/go2ban/jail"
)

type server struct {
	Js []*jail.Jail
}

var Srv server

func (s *server) DumpCells() {
	for _, j := range s.Js {
		fmt.Println(j.Cells)
		fmt.Println(len(j.Cells))
	}
}
