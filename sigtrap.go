package go2ban

import (
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
				os.Exit(0)
			case syscall.SIGQUIT:
				os.Exit(1)
			case syscall.SIGHUP:
			case syscall.SIGUSR1:
				Srv.DumpCells()
			}
		}
	}()
}
