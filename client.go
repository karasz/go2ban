package main

import (
	"fmt"
	"net"

	"github.com/karasz/go2ban/utils"
)

func main() {
	conn, err := net.DialUnix("unix", nil,
		&net.UnixAddr{utils.GoSock, "unix"})
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}
	var buf [1024]byte
	n, err := conn.Read(buf[:])

	fmt.Println(string(buf[:n]))
}
