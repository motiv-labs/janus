package request

// ContextKey is used to create context keys that are concurrent safe
type ContextKey string

func (c ContextKey) String() string {
	return "janus." + string(c)
}
