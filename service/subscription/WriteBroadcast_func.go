package subscription

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (s *SubscriptionService) WriteBroadcast() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Ws.WriteMessage(websocket.TextMessage, []byte("Message Start---")); err != nil {
				return
			}
			for _, m := range message.AllMessage {
				if err := c.Ws.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
					return
				}
			}
			if err := c.Ws.WriteMessage(websocket.TextMessage, []byte("---Message End")); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
