package pty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func NewPty() (*os.File, string, error) {
	master, err := os.OpenFile("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, "", err
	}

	err = unlockpt(master)
	if err != nil {
		return nil, "", err
	}

	slave, err := ptsname(master)
	if err != nil {
		return nil, "", err
	}

	return master, slave, nil
}

func ExecWithPty(command string, raw bool, args ...string) (io.ReadWriteCloser, error) {
	master, slaveName, err := NewPty()
	if err != nil {
		return nil, err
	}

	if raw {
		err = ttySetRaw(master)
		if err != nil {
			return nil, err
		}
	}

	slave, err := os.OpenFile(slaveName, syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	defer slave.Close()

	cmd := exec.Command(command, args...)
	cmd.Stdin = slave
	cmd.Stdout = slave
	cmd.Stderr = slave
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return master, nil
}

func ioctl(fd, flag, data uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, flag, data); err != 0 {
		return err
	}
	return nil
}

func unlockpt(f *os.File) error {
	var data int32
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&data)))
}

func ptsname(f *os.File) (string, error) {
	var pty int32
	err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&pty)))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%v", pty), nil
}

func ttySetRaw(f *os.File) error {
	termios, err := unix.IoctlGetTermios(int(f.Fd()), unix.TCGETS)
	if err != nil {
		return err
	}

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(int(f.Fd()), unix.TCSETS, termios)
	if err != nil {
		return err
	}

	return nil
}
