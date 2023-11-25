package main

import (
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"log"
	"net/http"
)

const link = "https://www.rediff.com/news/images10.html"

// TODO сделать рекурсивный парсинг

func handler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch URL, status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	div := doc.Find("#wrapper.mainwrapper")

	divs := div.Find("div.nbox")

	var data []struct {
		Image string
		Title string
	}

	divs.Each(func(i int, s *goquery.Selection) {
		image, _ := s.Find("img").Attr("src")
		title := s.Find("h2").Text()
		data = append(data, struct {
			Image string
			Title string
		}{Image: image, Title: title})
	})

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Rendered Page</title>
	<style>
		body {
			display: flex;
			flex-direction: column;
			align-items: center;
			height: 100vh;
			margin: 0;
			background-color: #f0f0f0;
		}
		.nbox {
			text-align: center;
		}
	</style>
	</head>
	<body>
	{{range .}}
	<div class="nbox">
		<img src="{{.Image}}">
		<h2>{{.Title}}</h2>
	</div>
	{{end}}
	</body>
	</html>
	`

	pageTmpl, err := template.New("page").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = pageTmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
