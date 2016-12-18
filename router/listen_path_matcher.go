package router

import "regexp"

type ListenPathMatcher struct {
	reg *regexp.Regexp
}

func NewListenPathMatcher() *ListenPathMatcher {
	return &ListenPathMatcher{regexp.MustCompile(`(\/\*(.+)?)`)}
}

func (l *ListenPathMatcher) Match(listenPath string) bool {
	return l.reg.MatchString(listenPath)
}

func (l *ListenPathMatcher) Extract(listenPath string) string {
	return l.reg.ReplaceAllString(listenPath, "")
}
