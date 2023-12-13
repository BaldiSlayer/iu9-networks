package client

import (
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
)

type Connection struct {
	conn *ftp.ServerConn
}

func (c *Connection) Connect(host, port, login, password string) {
	client, err := ftp.Dial(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Login(login, password)
	if err != nil {
		log.Println("Failed to connect to FTP server:", err)
	} else {
		log.Println("Successfully connected to FTP server")
		c.conn = client
	}

}

// UploadFile - загрузка файла
func (c *Connection) UploadFile(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	err = c.conn.Stor(path, file)
	if err != nil {
		log.Fatal(err)
	}
}

// DownloadFile - выкачка файла
func (c *Connection) DownloadFile(path string) {
	r, err := c.conn.Retr(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(path + "1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, r)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateDir - создание директории
func (c *Connection) CreateDir(path string) {
	err := c.conn.MakeDir("new-directory")
	if err != nil {
		log.Fatal(err)
	}
}

// DeleteFile - Удаление файла
func (c *Connection) DeleteFile(path string) {
	err := c.conn.Delete("remote-file.txt")
	if err != nil {
		log.Fatal(err)
	}
}

// GetDirectoryContent - функция для получения содержимого директории
func (c *Connection) GetDirectoryContent(path string) ([]string, error) {
	entries, err := c.conn.List(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name)
	}

	return fileNames, nil
}
