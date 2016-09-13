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
		if len(j.Cells) == 0 {
			fmt.Println("No Cells active")
		} else {
			for k, v := range j.Cells {
				fmt.Println("IP: ", k, "\tFirst observed at: ", v)
			}
		}
	}
}
