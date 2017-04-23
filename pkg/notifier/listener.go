package notifier

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// NotificationListener listens for
type NotificationListener struct {
	subscriber Subscriber
}

// NewNotificationListener creates a new instance of NotificationListener
func NewNotificationListener(subscriber Subscriber) *NotificationListener {
	return &NotificationListener{subscriber}
}

// Start starts listening for signals on the cluster
func (n *NotificationListener) Start(fn func(v Notification)) {
	log.Debug("Listening for changes")

	go func() {
		for {
			err := n.subscriber.Subscribe(DefaultChannel, fn)
			if err != nil {
				log.WithFields(log.Fields{
					"prefix": "pub-sub",
					"err":    err,
				}).Error("Connection failed, reconnect in 10s")

				time.Sleep(10 * time.Second)
				log.WithFields(log.Fields{
					"prefix": "pub-sub",
				}).Warning("Reconnecting")
			}
		}
	}()
}
