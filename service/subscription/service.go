package subscription

import (
	"github.com/rudiarta/warung_pintar/model"
	"github.com/rudiarta/warung_pintar/util"
)

type Service interface {
	WriteBroadcast()
}

type SubscriptionService struct {
	Conn     *model.Connection
	NameOrId string
}

func NewSubscriptionService(c *model.Connection) SubscriptionService {
	return SubscriptionService{c, util.GenerateUUID()}
}
