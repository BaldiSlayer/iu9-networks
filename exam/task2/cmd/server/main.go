package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func handleConnection(conn net.Conn) {
	//defer conn.Close()

	fmt.Println(2)

	// Принимаем URL от клиента
	fmt.Println(conn)
	url, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(3)
	fmt.Println("Получен URL:", url)

	url = strings.TrimRight(url, "\n")
	// Загружаем JSON с указанного URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Преобразуем содержимое ответа в строку
	bodyString := string(body)

	fmt.Println(bodyString)
	_, err = conn.Write(append([]byte(bodyString), '\n'))
	if err != nil {
		fmt.Println("Error sending data to client:", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен. Ожидание подключений...")

	for {
		conn, err := listener.Accept()
		fmt.Println(1)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
