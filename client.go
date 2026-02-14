package wsroutes

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) readLoop(routes map[string]HandlerFunc) {
	if handler, ok := routes["/connect"]; ok {
		handler(c, nil)
	}
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		var event string
		json.Unmarshal(m["event"], &event)

		if handler, ok := routes[event]; ok {
			handler(c, m["data"])
		}
	}

	if handler, ok := routes["/disconnect"]; ok {
		handler(c, nil)
	}
}
