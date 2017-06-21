package notifier

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// NotificationCommand represents a notification command
type NotificationCommand string

const (
	// NoticeAPIUpdated notifies when an API is updated
	NoticeAPIUpdated NotificationCommand = "ApiUpdated"
	// NoticeAPIRemoved notifies when an API is removed
	NoticeAPIRemoved NotificationCommand = "ApiRemoved"
	// NoticeAPIAdded notifies when an API is added
	NoticeAPIAdded NotificationCommand = "ApiAdded"
	// NoticeOAuthServerUpdated notifies when an OAuth server is updated
	NoticeOAuthServerUpdated NotificationCommand = "OAuthUpdated"
	// NoticeOAuthServerRemoved notifies when an OAuth server is removed
	NoticeOAuthServerRemoved NotificationCommand = "OAuthRemoved"
	// NoticeOAuthServerAdded notifies when an OAuth server is added
	NoticeOAuthServerAdded NotificationCommand = "OAuthAdded"
	// DefaultChannel represents the default channel's name
	DefaultChannel = "janus.cluster.notifications"
)

// Subscriber holds the basic methods to subscribe to a topic
type Subscriber interface {
	Subscribe(channel string, callback func(Notification)) error
}

// Publisher holds the basic methods to publish a message
type Publisher interface {
	Publish(topic string, data []byte) error
}

// Notification is a type that encodes a message published to a pub sub channel (shared between implementations)
type Notification struct {
	Command   NotificationCommand `json:"command"`
	Payload   string              `json:"payload"`
	Signature string              `json:"signature"`
}

// Notifier holds the basic methods to notify listeners
type Notifier interface {
	Notify(notification Notification) bool
}

// PublisherNotifier will use redis pub/sub channels to send notifications
type PublisherNotifier struct {
	publisher Publisher
	channel   string
}

// NewPublisherNotifier creates a new instance of Notifier
func NewPublisherNotifier(publisher Publisher, channel string) *PublisherNotifier {
	if channel == "" {
		channel = DefaultChannel
	}

	return &PublisherNotifier{publisher, channel}
}

// Notify will send a notification to a channel
func (r *PublisherNotifier) Notify(notification Notification) bool {
	toSend, err := json.Marshal(notification)
	if err != nil {
		log.WithError(err).Error("Problem marshalling notification")
		return false
	}

	log.WithField("type", notification.Command).Debug("Sending notification")
	if err := r.publisher.Publish(r.channel, toSend); err != nil {
		log.WithError(err).Error("Could not send notification")
		return false
	}
	return true
}

// RequireReload checks if a given command requires reload
func RequireReload(cmd NotificationCommand) bool {
	switch cmd {
	case NoticeAPIUpdated, NoticeAPIRemoved, NoticeAPIAdded, NoticeOAuthServerUpdated, NoticeOAuthServerRemoved, NoticeOAuthServerAdded:
		return true
	default:
		return false
	}
}
