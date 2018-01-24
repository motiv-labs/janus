package router

import "regexp"

const (
	parameterMatchRule = `\{([^/}]+)\}`
)

// ListenPathParameterNameExtractor is responsible for extracting parameters name from the listen path
type ListenPathParameterNameExtractor struct {
	reg *regexp.Regexp
}

// NewListenPathParamNameExtractor creates a new instance ListenPathParameterNameExtractor
func NewListenPathParamNameExtractor() *ListenPathParameterNameExtractor {
	return &ListenPathParameterNameExtractor{regexp.MustCompile(parameterMatchRule)}
}

// Extract takes the usable part of the listen path and extracts parameter names
func (l *ListenPathParameterNameExtractor) Extract(listenPath string) []string {
	submatches := l.reg.FindAllStringSubmatch(listenPath, -1)
	result := make([]string, 0, len(submatches))

	for _, submatch := range submatches {
		result = append(result, submatch[1])
	}

	return result
}
