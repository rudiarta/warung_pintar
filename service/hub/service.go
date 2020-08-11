package hub

import (
	"github.com/rudiarta/warung_pintar/model"
	"github.com/rudiarta/warung_pintar/service/subscription"
)

type Service interface {
	Run()
}

func NewHub() HubService {
	return HubService{
		Register:   make(chan subscription.SubscriptionService),
		Unregister: make(chan subscription.SubscriptionService),
		Brodcast:   make(chan string),
		Pool:       make(map[string]*model.Connection),
	}
}

type HubService struct {
	Register   chan subscription.SubscriptionService
	Unregister chan subscription.SubscriptionService
	Brodcast   chan string
	Pool       map[string]*model.Connection
}
