package wsroutes

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
)

type HandlerFunc func(*Client, []byte)
type WsRoutes struct {
	httpRoute string
	port      string
	routes    map[string]HandlerFunc
	upgrader  websocket.Upgrader
}

func New(httpRoute, port string, upgrader websocket.Upgrader) *WsRoutes {
	return &WsRoutes{
		httpRoute: httpRoute,
		port:      port,
		routes:    make(map[string]HandlerFunc),
		upgrader:  upgrader,
	}
}

func (ws *WsRoutes) validateRoute(route string) bool {
	matched, _ := regexp.MatchString(`^/[a-zA-Z]+(/[a-zA-Z]+)*$`, route)
	//log.Fatal(err, matched, route)
	return matched

}

func (ws *WsRoutes) On(route string, callback HandlerFunc) {
	if !ws.validateRoute(route) {
		panic("invalid route path: " + route)
	}

	ws.routes[route] = callback
}

func (ws *WsRoutes) OnConnect(callback HandlerFunc) {
	_, exist := ws.routes["/connect"]
	if exist {
		panic("connect route exist")
	}
	ws.routes["/connect"] = callback
}
func (ws *WsRoutes) OnDisconnect(callback HandlerFunc) {
	_, exist := ws.routes["/disconnect"]
	if exist {
		panic("disconnect route exist")
	}
	ws.routes["/disconnect"] = callback
}

func (ws *WsRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{Conn: conn}

	go client.readLoop(ws.routes)
}
