package main

import (
	"encoding/json"
	"github.com/Baldislayer/iu9-networks/lab5/models"
	"github.com/gorilla/websocket"
	"iu9-networks/lab5/configs"
	"log"
)

const (
	host = "185.139.70.64"
)

const path = "ws://" + host + configs.Port + "/ws"

func main() {
	// поключаемся к серверу посредством вебсокета
	conn, _, err := websocket.DefaultDialer.Dial(path, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server_2:", err)
	}

	// отложенная операция, при завершении функции main мы закроем соединение
	// подробнее про defer в файле server_2.go
	defer conn.Close()

	for {
		// ты спросишь меня наверное, почему я не сделал присвоение по ссылке
		// а сделал через return
		// дело в том, что так делают обычно в языке Golang
		// это показывает, что данные, которые мы передаем в качестве параметра - мутируют (как-то изменяются)
		// и человеку,который будет использовать твой код - будет намного легче это понять
		message := models.Message{}
		message = message.InputFromConsole()

		// переводим нашу структуру в jsong
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			continue
		}

		// отправляем как текстовое сообщение наш запрос на сервер по вебсокету
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			continue
		}

		// принимаем ответ
		_, res, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		// читаем ответ, которое нам прислал сервер и с помощью функции
		// json.Unmarshal переводим из формата json
		// в нашу гошную структуру
		var response models.Response
		err = json.Unmarshal(res, &response)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		// выводим все в консоль
		response.OutToConsole()
	}
}
