package main

import (
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/iu9-networks/lab5/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const port = ":8080"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleClient(conn *websocket.Conn) {
	// закрываем соединение при заверщении нашей функции
	// все, что написано после defer, выполняется при завершении функции, в которой
	// находится этот defer, если их несколько, то они выполняются в порядке стека
	// последняя в коде функции будет выполнена первой
	defer conn.Close()

	fmt.Println("Соединение открыто")
	// запускаем бесконенчый цикл обработки сообщений от клиента
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// если нам выбило такую ошибку, то это значит, что просто от нас отключился клиент
			if err.Error() == "websocket: close 1006 (abnormal closure): unexpected EOF" {
				fmt.Println("Соединение закрыто")
				break
			}
			log.Println("Error reading message:", err)

			break
		}

		// читаем сообщение, которое нам прислал пользователь и с помощью функции
		// json.Unmarshal переводим из формата json в котором нам прислал его клиент
		// в нашу гошную структуру
		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			break
		}

		// выводим то, что получили, это чисто для дебага, можешь убрать
		fmt.Println(message)

		// вызываем нашу функцию "подсчета" результата
		response, err := models.ResultFunction(message)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			break
		}

		// кастуем обратно в json для дальнейшей отправки результата на клиент
		res, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			break
		}

		// передаем наше сообщение на клиент как текст
		err = conn.WriteMessage(websocket.TextMessage, res)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	// тут мы собственно создаем "ручку" для обработки запросов на присоединение от клиентов
	// ручка - обработчик каких-то определенных событий
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		go handleClient(conn)
	})

	// поднимаем наш сервер
	log.Fatal(http.ListenAndServe(port, nil))
}
