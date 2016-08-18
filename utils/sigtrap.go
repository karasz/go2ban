package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const GoSock = "/var/run/go2ban/socket"

func TrapSignals() {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)

		for sig := range sigchan {
			switch sig {
			case syscall.SIGTERM:
				log.Println("[INFO] SIGTERM: Terminating process")
				os.Remove(GoSock)
				os.Exit(0)

			case syscall.SIGQUIT:
				log.Println("[INFO] SIGQUIT: Shutting down")
				os.Remove(GoSock)
				os.Exit(1)

			case syscall.SIGHUP:
				log.Println("[INFO] SIGHUP: Hanging up")

			case syscall.SIGUSR1:
				log.Println("[INFO] SIGUSR1: Reloading")

			}
		}
	}()
}
