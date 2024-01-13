package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var mu sync.Mutex
var connections = make(map[*websocket.Conn]bool)
var json_st string

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}

	mu.Lock()
	connections[conn] = true
	mu.Unlock()

	_, message, err := conn.ReadMessage() // Чтение сообщения из веб-сокета
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}
	fmt.Printf("Received message from client: %s\n", message) // Вывод сообщения в консоль

	// Send the current number to the newly connected client
	//if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", num))); err != nil {
	//	fmt.Println("Error writing message:", err)
	//	return
	//}

	fmt.Println(1)
	tcp_conn, err := net.Dial("tcp", "185.255.133.113:8888")
	if err != nil {
		fmt.Println("Ошибка соединения с сервером:", err)
		return
	}
	defer tcp_conn.Close()

	fmt.Println(1)

	message = append(message, '\n')
	// Отправляем данные на сервер
	_, err = tcp_conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки сообщения:", err)
		return
	}

	fmt.Println(4)

	// Принимаем ответ от сервера
	response, err := bufio.NewReader(tcp_conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка приема сообщения:", err)
		return
	}
	fmt.Println("Получен ответ от сервера:", response)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

	defer conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":8798", nil)
}
