package utils

import (
	"fmt"
	"github.com/karasz/go2ban/common"
	"os"
	"os/signal"
	"syscall"
)

func TrapSignals() {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)

		for sig := range sigchan {
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				fmt.Println("[INFO] SIGTERM or SIGINT: Terminating process")
				os.Remove(common.GoSock)
				os.Exit(0)

			case syscall.SIGQUIT:
				fmt.Println("[INFO] SIGQUIT: Shutting down")
				os.Remove(common.GoSock)
				os.Exit(1)

			case syscall.SIGHUP:
				fmt.Println("[INFO] SIGHUP: Hanging up")

			case syscall.SIGUSR1:
				fmt.Println("[INFO] SIGUSR1: Reloading")

			}
		}
	}()
}
