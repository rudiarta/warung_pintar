package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rudiarta/warung_pintar/util"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Response struct {
	Status string `json:"status"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type message struct {
	AllMessage []string `json:"all-message"`
}

var text = &message{}

type subscription struct {
	conn     *connection
	nameOrId string
}

func (s *subscription) WriteBroadcast() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, []byte("Message Start---")); err != nil {
				return
			}
			for _, m := range message.AllMessage {
				if err := c.ws.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
					return
				}
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, []byte("---Message End")); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

type connection struct {
	ws   *websocket.Conn
	send chan message
}

type hub struct {
	register   chan subscription
	unregister chan subscription
	brodcast   chan string
	pool       map[string]*connection
}

var NewHub = hub{
	register:   make(chan subscription),
	unregister: make(chan subscription),
	brodcast:   make(chan string),
	pool:       make(map[string]*connection),
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			h.pool[s.nameOrId] = s.conn
		case s := <-h.unregister:
			delete(h.pool, s.nameOrId)
			close(s.conn.send)
		case m := <-h.brodcast:
			text.AllMessage = append(text.AllMessage, m)
			connections := h.pool
			for _, c := range connections {
				select {
				case c.send <- *text:
				}
			}
		}
	}
}

func WebSocketHandle(h hub) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}

		c := &connection{ws: ws, send: make(chan message)}
		s := subscription{c, util.GenerateUUID()}
		h.register <- s
		go s.WriteBroadcast()
	}
}

func SendMessageHandle(h hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.brodcast <- r.FormValue("message")
		w.Header().Set("Content-Type", "application/json")
		result := Response{
			Status: "success",
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetAllMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(text)

}

func main() {
	go NewHub.run()

	handleWebsocket := WebSocketHandle(NewHub)
	http.HandleFunc("/", handleWebsocket)

	handleSendMessage := SendMessageHandle(NewHub)
	http.HandleFunc("/send-message", handleSendMessage)

	http.HandleFunc("/get-all-message", GetAllMessage)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Printf("Error: %v", err.Error())
	}
}
