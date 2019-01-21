package observability

import "context"

type reqIDKeyType int

const reqIDKey reqIDKeyType = iota

// RequestIDToContext puts a request ID to context for future use
func RequestIDToContext(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		panic("Can not put request ID to empty context")
	}

	return context.WithValue(ctx, reqIDKey, requestID)
}

// RequestIDFromContext tries to extract request ID from context if present, otherwise returns empty string
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		panic("Can not get request ID from empty context")
	}

	if requestID, ok := ctx.Value(reqIDKey).(string); ok {
		return requestID
	}

	return ""
}
