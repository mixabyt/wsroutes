package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mixabyt/wsroutes"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func t(c *wsroutes.EventHandler, msg []byte) {
	fmt.Println(string(msg))
	c.Conn.WriteMessage(websocket.TextMessage, []byte("test"))
}

func main() {

	ws := wsroutes.New("/ws", ":8080", upgrader)
	ws.OnConnect(func(client *wsroutes.EventHandler, bytes []byte) {
		fmt.Println("user connected")
	})
	ws.OnDisconnect(func(client *wsroutes.EventHandler, bytes []byte) {
		fmt.Println("user disconnected")
	})

	ws.On("/hello", nil)
	ws.On("/test", t)

	http.Handle("/ws", ws)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
