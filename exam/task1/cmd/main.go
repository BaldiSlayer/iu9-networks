package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
var number int

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}

	mu.Lock()
	connections[conn] = true
	num := number // Get the current number
	mu.Unlock()

	// Send the current number to the newly connected client
	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", num))); err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

	// defer conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			if num, err := strconv.Atoi(input); err == nil {
				// Update the global number and send it to all connected clients
				mu.Lock()
				number = num
				fmt.Println(connections)
				for client := range connections {
					if err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", number))); err != nil {
						fmt.Println("Error writing message:", err)
					}
				}
				mu.Unlock()
			} else {
				fmt.Println("Invalid input. Please enter a valid number.")
			}
		}
	}()

	http.ListenAndServe(":8798", nil)
}
