package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	con "iu9-networks/lab4/cmd/config"
	"log"
	"os"
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
			ssh.Password(con.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // for the demo only, use proper host key verification in production
	}

	client, err := ssh.Dial("tcp", "localhost:31337", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	RunCmd(client, "ls")
	RunCmd(client, "mkdir fsdafsd")
	RunCmd(client, "ls")
}
