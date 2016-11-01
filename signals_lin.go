// +build !windows

package main

import (
	"os"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGTSTP,
}
