// +build windows

package main

import (
	"os/exec"
	"syscall"
)

// Myexec for os specific switches
func Exec(engine string) *exec.Cmd {
	cmd := exec.Command(engine)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
