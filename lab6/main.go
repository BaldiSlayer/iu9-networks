package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Category    string `xml:"category"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL  string `xml:"url,attr"`
		Type string `xml:"type,attr"`
	} `xml:"enclosure"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

func main() {
	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS news_lisov (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255),
			link VARCHAR(255),
			description TEXT,
			category VARCHAR(100),
			pub_date VARCHAR(100),
			enclosure_url VARCHAR(255),
			enclosure_type VARCHAR(50)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		ALTER TABLE news_lisov CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	`)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get("https://briansk.ru/rss20_briansk.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel

	var rss RSS
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range rss.Items {
		// Проверка наличия новости в базе данных
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM news_lisov WHERE title = ?", item.Title).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		// Если новости нет в базе данных, добавляем ее
		if count == 0 {
			_, err := db.Exec("INSERT INTO news_lisov (title, link, description, category, pub_date, enclosure_url, enclosure_type) VALUES (?, ?, ?, ?, ?, ?, ?)",
				item.Title, item.Link, item.Description, item.Category, item.PubDate, item.Enclosure.URL, item.Enclosure.Type)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Добавлена новая новость: %s\n", item.Title)
		}
	}
}
