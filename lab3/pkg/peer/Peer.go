package peer

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"iu9-networks/lab3/models"
	"log"
	"net"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Peer struct {
	Name      string
	Info      models.Node
	Neighbors map[string]models.Node
	conn      *websocket.Conn
}

var tpl *template.Template

// StartWebServer - метод "поднятия" веб сервера, отвечающего за интерфейс
func (p *Peer) StartWebServer() {
	tpl, _ = template.ParseFiles("../index.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, struct{ NodeName string }{NodeName: p.Name})
	})

	http.HandleFunc("/ws/"+p.Name, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		p.conn = conn

		for {
			// Чтение сообщений из веб-сокета
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Ошибка чтения сообщения из веб-сокета:", err)
				return
			}

			msg := models.Line{}
			err = json.Unmarshal(bytes, &msg)
			if err != nil {
				fmt.Println("Ошибка десериализации JSON:", err)
				continue
			}

			p.drawLine(msg)
		}
	})

	fmt.Println("Starting web server_2 on port ", p.Info.HtmlServerPort)
	http.ListenAndServe(":"+p.Info.HtmlServerPort, nil)
}

// drawLine - функция "отрисовки" линии
// передает информацию о нарисовании линии всем соседям
func (p *Peer) drawLine(msg models.Line) {
	// Отправка полученного сообщения всем узлам
	p.SendMessage(msg)
	// и себе не забыть нарисовать
	p.conn.WriteJSON(msg)
}

// StartSocket - стартуем сокет
func (p *Peer) StartSocket() {
	ln, _ := net.Listen("tcp", ":"+p.Info.Port)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go p.handleConnection(conn)
	}
}

// handleConnection - обработка входящих соединений
func (p *Peer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	var message models.Line
	json.Unmarshal(buf[:n], &message)

	// теперь просто по вебсокету отправляем

	if p.conn != nil {
		err := p.conn.WriteJSON(message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("Received message: %+v\n", message)
}

// SendMessage - посылка сообщения всем соседям
func (p *Peer) SendMessage(message models.Line) {
	for _, neighbor := range p.Neighbors {
		conn, err := net.Dial("tcp", ":"+neighbor.Port)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer conn.Close()

		jsonMessage, _ := json.Marshal(message)
		_, err = conn.Write(jsonMessage)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
