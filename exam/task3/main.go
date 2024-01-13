package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Item struct {
	Name    string
	GPS     string
	Address string
	Tel     string
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "http://pstgu.yss.su/iu9/mobiledev/lab4_yandex_map/?x=var05"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Ошибка при выполнении GET-запроса", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении ответа", http.StatusInternalServerError)
		return
	}

	var data []Item
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
