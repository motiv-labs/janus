package metrics

import (
	"context"

	stats "github.com/hellofresh/stats-go"
)

type statsKeyType int

const statsKey statsKeyType = iota

// NewContext returns a context that has a stats Client
func NewContext(ctx context.Context, client stats.Client) context.Context {
	return context.WithValue(ctx, statsKey, client)
}

// WithContext returns a stats Client with as much context as possible
func WithContext(ctx context.Context) stats.Client {
	ctxStats, ok := ctx.Value(statsKey).(stats.Client)
	if !ok {
		ctxStats, _ := stats.NewClient("noop://", "")
		return ctxStats
	}
	return ctxStats
}
