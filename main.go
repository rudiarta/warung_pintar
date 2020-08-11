package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rudiarta/warung_pintar/model"
	"github.com/rudiarta/warung_pintar/service/hub"
	"github.com/rudiarta/warung_pintar/service/subscription"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandle(h hub.HubService) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}

		c := &model.Connection{Ws: ws, Send: make(chan model.Message)}
		s := subscription.NewSubscriptionService(c)
		h.Register <- s
		go s.WriteBroadcast()
	}
}

func SendMessageHandle(h hub.HubService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Brodcast <- r.FormValue("message")
		w.Header().Set("Content-Type", "application/json")
		result := model.Response{
			Status: "success",
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetAllMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.Text)

}

func main() {
	h := hub.NewHub()
	go h.Run()

	handleWebsocket := WebSocketHandle(h)
	http.HandleFunc("/", handleWebsocket)

	handleSendMessage := SendMessageHandle(h)
	http.HandleFunc("/send-message", handleSendMessage)

	http.HandleFunc("/get-all-message", GetAllMessage)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Printf("Error: %v", err.Error())
	}
}
