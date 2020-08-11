package hub

import "github.com/rudiarta/warung_pintar/model"

func (h *HubService) Run() {
	for {
		select {
		case s := <-h.Register:
			h.Pool[s.NameOrId] = s.Conn
		case s := <-h.Unregister:
			delete(h.Pool, s.NameOrId)
			close(s.Conn.Send)
		case m := <-h.Brodcast:
			model.Text.AllMessage = append(model.Text.AllMessage, m)
			connections := h.Pool
			for _, c := range connections {
				select {
				case c.Send <- *model.Text:
				}
			}
		}
	}
}
