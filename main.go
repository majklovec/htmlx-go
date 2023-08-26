package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	tmpl := template.Must(template.New("timeTemplate").Parse(`<div id="time">{{.}}</div>`))

	for {
		data := time.Now().Format("2006-01-02 15:04:05")

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			return
		}

		time.Sleep(time.Second)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}
