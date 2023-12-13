package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"iu9-networks/lab6/internal/client"
	rss "iu9-networks/lab6/models"
	"net/http"
	"os"
	"time"
)

const (
	host     = "students.yss.su"
	port     = "21"
	login    = "ftpiu8"
	password = "3Ru7yOTA"
)

func getRSSData() (rss.RSS, error) {
	resp, err := http.Get("https://news.rambler.ru/rss/Guadeloupe/")
	if err != nil {
		return rss.RSS{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rss.RSS{}, err
	}

	var result rss.RSS
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return rss.RSS{}, err
	}

	return result, nil
}

func removeFiles(files ...string) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Failed to remove file %s: %v\n", file, err)
		}
	}
}

func main() {
	rssData, err := getRSSData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	currentTime := time.Now()
	fileName := fmt.Sprintf("<%s>_%s.txt", "Lisov Aleksey", currentTime.Format("20060102150405"))

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	for _, item := range rssData.Channel.Items {
		file.WriteString(item.Title + "\n")
		file.WriteString(item.Description + "\n")
		file.WriteString(item.Link + "\n")
		file.WriteString("------------------------------\n")
	}

	fmt.Println("Data saved to file:", fileName)

	var currentConnection client.Connection

	currentConnection.Connect(host, port, login, password)

	currentConnection.UploadFile(fileName)

	removeFiles(fileName)
}
