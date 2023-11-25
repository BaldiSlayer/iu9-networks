package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

var (
	Login    = ""
	Password = "123"
	Ip       = "127.0.0.1"
	Port     = 31337
)

func RunCmd(client *ssh.Client, cmd string) {
	fmt.Println("Running cmd: ", cmd)
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run command: ", err)
	}
}

func main() {
	config := &ssh.ClientConfig{
		User: "",
		Auth: []ssh.AuthMethod{
			ssh.Password(Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // for the demo only, use proper host key verification in production
	}

	client, err := ssh.Dial("tcp", "localhost:31337", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	RunCmd(client, "ping yandex.ru")
}
