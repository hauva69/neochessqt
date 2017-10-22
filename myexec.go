// +build !windows

package main

import "os/exec"

// Myexec for os specific switches
func Exec(engine string) *exec.Cmd {
	cmd := exec.Command(engine)
	return cmd
}
