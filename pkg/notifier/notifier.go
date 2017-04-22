package notifier

import (
	"encoding/json"

	"github.com/containous/traefik/log"
	"github.com/hellofresh/janus/pkg/store"
)

type NotificationCommand string

const (
	NoticeApiUpdated         NotificationCommand = "ApiUpdated"
	NoticeApiRemoved         NotificationCommand = "ApiRemoved"
	NoticeApiAdded           NotificationCommand = "ApiAdded"
	NoticeOAuthServerUpdated NotificationCommand = "OAuthUpdated"
	NoticeOAuthServerRemoved NotificationCommand = "OAuthRemoved"
	NoticeOAuthServerAdded   NotificationCommand = "OAuthAdded"
)

// RequireReload checks if a given command requires reload
func RequireReload(cmd NotificationCommand) bool {
	switch cmd {
	case NoticeApiUpdated, NoticeApiRemoved, NoticeApiAdded, NoticeOAuthServerUpdated, NoticeOAuthServerRemoved, NoticeOAuthServerAdded:
		return true
	default:
		return false
	}
}

// Notification is a type that encodes a message published to a pub sub channel (shared between implementations)
type Notification struct {
	Command   NotificationCommand `json:"command"`
	Payload   string              `json:"payload"`
	Signature string              `json:"signature"`
}

// Notifier will use redis pub/sub channels to send notifications
type Notifier struct {
	publisher store.Publisher
	channel   string
}

// New creates a new instance of Notifier
func New(publisher store.Publisher, channel string) *Notifier {
	return &Notifier{publisher, channel}
}

// Notify will send a notification to a channel
func (r *Notifier) Notify(notification Notification) bool {
	toSend, err := json.Marshal(notification)
	if err != nil {
		log.Error("Problem marshalling notification: ", err)
		return false
	}
	log.Debug("Sending notification", notification)
	if err := r.publisher.Publish(r.channel, toSend); err != nil {
		log.Error("Could not send notification: ", err)
		return false
	}
	return true
}
