//go:build !windows
// +build !windows

// reverseSSH - a lightweight ssh server_2 with a reverse connection feature
// Copyright (C) 2021  Ferdinor <ferdinor@mailbox.org>

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"os/user"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
)

func createPty(s ssh.Session, shell string) {
	var (
		ptyReq, winCh, _ = s.Pty()
		cmd              = exec.CommandContext(s.Context(), shell)
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
	if currentUser, err := user.Current(); err == nil {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", currentUser.HomeDir))
	}
	f, err := pty.Start(cmd)
	if err != nil {
		log.Fatalln("Could not start shell:", err)
	}
	go func() {
		for win := range winCh {
			winSize := &pty.Winsize{Rows: uint16(win.Height), Cols: uint16(win.Width)}
			pty.Setsize(f, winSize)
		}
	}()

	go func() {
		io.Copy(f, s)
		s.Close()
	}()
	go func() {
		io.Copy(s, f)
		s.Close()
	}()

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			log.Println("Session ended with error:", err)
			s.Exit(255)
			return
		}
		log.Println("Session ended normally")
		s.Exit(cmd.ProcessState.ExitCode())
		return

	case <-s.Context().Done():
		log.Printf("Session terminated: %s", s.Context().Err())
		return
	}
}
