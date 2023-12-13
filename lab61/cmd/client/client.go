package main

import (
	"fmt"
	"iu9-networks/lab6/internal/client"
	"log"
	"os"
)

const (
	host     = "localhost"
	port     = "2121"
	login    = "user"
	password = "123456"

	fileName = "lisov"
)

func createFile(name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func contains(files []string, fileName string) bool {
	for _, value := range files {
		if value == fileName {
			return true
		}
	}
	return false
}

func testFileUploading(currentConnection *client.Connection) {
	createFile(fileName)

	currentConnection.UploadFile(fileName)

	sl, err := currentConnection.GetDirectoryContent(".")
	if err != nil {
		fmt.Println("Failed upload test")
		panic(nil)
	}

	if !contains(sl, fileName) {
		fmt.Println("Failed upload test")
		panic(nil)
	}

	fmt.Println("Success upload test")
}

func fileExist(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	} else {
		return false
	}
}

func testFileDownloading(currentConnection *client.Connection) {
	currentConnection.DownloadFile(fileName)

	if !fileExist(fileName + "1") {
		fmt.Println("Failed download test")
		panic(nil)
	}

	fmt.Println("Success download test")
}

func removeFiles(files ...string) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Failed to remove file %s: %v\n", file, err)
		}
	}
}

func test1(currentConnection *client.Connection) {
	testFileUploading(currentConnection)
	testFileDownloading(currentConnection)

	removeFiles(fileName, fileName+"1")
}

func main() {
	var currentConnection client.Connection

	currentConnection.Connect(host, port, login, password)

	// погнали тестить
	test1(&currentConnection)

	/*sl, err := currentConnection.GetDirectoryContent("")
	if err != nil {
		panic(nil)
	}
	fmt.Println(sl)*/

	// подчистка

}
