package notifier

import (
	"time"

	log "github.com/sirupsen/logrus"
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
	log.Debug("Listening for change events")
	logWithFields := log.WithFields(log.Fields{
		"prefix": "pub-sub",
	})

	go func() {
		for {
			if err := n.subscriber.Subscribe(DefaultChannel, fn); err != nil {
				log.Info("n.subscriber.Subscribe")
				logWithFields.
					WithError(err).
					Error("Connection failed, reconnect in 10s")

				time.Sleep(10 * time.Second)

				logWithFields.Warning("Reconnecting")
			}
		}
	}()
}
