package main

import (
	"fmt"
	"net"
	"os"

	"github.com/karasz/go2ban/utils"
)

func init() {
	utils.TrapSignals()
	os.Remove(utils.GoSock)
}

func main() {

	l, err := net.ListenUnix("unix", &net.UnixAddr{utils.GoSock, "unix"})
	if err != nil {
		fmt.Println("listen error:", err)
	}

	defer os.Remove(utils.GoSock)

	for {
		fd, err := l.AcceptUnix()
		if err != nil {
			fmt.Println("accept error:", err)
		}

		go startServer(fd)
	}

}

func startServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		n, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]
		fmt.Println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			fmt.Println("Write: ", err)
		}
	}
}
