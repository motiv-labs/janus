package router

import "regexp"

const (
	matchRule = `(\/\*(.+)?)`
)

// ListenPathMatcher is responsible for matching a listen path to a set of rules
type ListenPathMatcher struct {
	reg *regexp.Regexp
}

// NewListenPathMatcher creates a new instance ListenPathMatcher
func NewListenPathMatcher() *ListenPathMatcher {
	return &ListenPathMatcher{regexp.MustCompile(matchRule)}
}

// Match verifies if a listen path matches the given rule
func (l *ListenPathMatcher) Match(listenPath string) bool {
	return l.reg.MatchString(listenPath)
}

// Extract takes the usable part of the listen path based on the provided rule
func (l *ListenPathMatcher) Extract(listenPath string) string {
	return l.reg.ReplaceAllString(listenPath, "")
}
