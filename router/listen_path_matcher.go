package router

import "regexp"

const (
	matchRule string = `(\/\*(.+)?)`
)

type ListenPathMatcher struct {
	reg *regexp.Regexp
}

func NewListenPathMatcher() *ListenPathMatcher {
	return &ListenPathMatcher{regexp.MustCompile(matchRule)}
}

func (l *ListenPathMatcher) Match(listenPath string) bool {
	return l.reg.MatchString(listenPath)
}

func (l *ListenPathMatcher) Extract(listenPath string) string {
	return l.reg.ReplaceAllString(listenPath, "")
}
