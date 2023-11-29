package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Category    string `xml:"category"`
	Author      string `xml:"author"`
}

func getRSSData() (RSS, error) {
	resp, err := http.Get("https://news.rambler.ru/rss/Guadeloupe/")
	if err != nil {
		return RSS{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RSS{}, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return RSS{}, err
	}

	return rss, nil
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	page := `
<html>
<head>
    <title>My Website</title>
</head>
<body>
    <h1>Welcome to My Website</h1>
    <ul>
        <li><a href="/">Home</a></li>
        <li><a href="/about">About</a></li>
        <li><a href="/rss">RSS</a></li>
    </ul>
</body>
</html>`

	io.WriteString(w, page)
}

func getRSSPage(w http.ResponseWriter, r *http.Request) {
	rss, err := getRSSData()
	if err != nil {
		http.Error(w, "Error fetching RSS data", http.StatusInternalServerError)
		return
	}

	// Вывод информации из RSS на страницу
	fmt.Fprintf(w, "<h1>%s</h1>\n", rss.Channel.Title)
	fmt.Fprintf(w, "<p>%s</p>\n", rss.Channel.Description)

	for _, item := range rss.Channel.Items {
		fmt.Fprintf(w, "<h2><a href='%s'>%s</a></h2>\n", item.Link, item.Title)
		fmt.Fprintf(w, "<p>%s</p>\n", item.Description)
		fmt.Fprintf(w, "<p>Category: %s</p>\n", item.Category)
		fmt.Fprintf(w, "<p>Author: %s</p>\n", item.Author)
		fmt.Fprintf(w, "<p>Published Date: %s</p>\n", item.PubDate)
	}
}

func getAbout(w http.ResponseWriter, r *http.Request) {
	// Создаем HTML шаблон с изображением
	tmpl := `<html>
<head>
  <title>About</title>
</head>
<body>
  <h1>Привет, это лаба Алексея Лисова.</h1>
  <img src="/static/me.jpg" alt="Image Alt Text">
</body>
</html>`

	t, err := template.New("about").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("error parsing template: %s\n", err)
		return
	}

	// Выполняем шаблон и записываем результат в ResponseWriter
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("error executing template: %s\n", err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из параметров запроса
	data := r.URL.Query()

	// Вывод данных в консоль
	fmt.Println("Received GET request to /test with data:")
	for key, values := range data {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Отправка ответа клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received and logged.\n"))
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/rss", getRSSPage)
	http.HandleFunc("/about", getAbout)
	http.HandleFunc("/test", testHandler)

	// Устанавливаем обработчик для статических файлов в каталоге "static"
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
