package main

import (
	"flag"

	log "github.com/mgutz/logxi/v1"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

func main() {
	var (
		root = flag.String("root", "test_dir", "Root directory to serve")
		user = flag.String("user", "user", "Username for login")
		pass = flag.String("pass", "123456", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "localhost", "Host")
	)
	flag.Parse()
	if *root == "" {
		log.Error("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
	}

	log.Info("Starting ftp server on", "host: ", opts.Hostname, "port: ", opts.Port)
	log.Info("Connection with: ", "name: ", *user, "password: ", *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
	}
}
