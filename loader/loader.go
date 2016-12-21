package loader

// Loader holds the methods that need to be implemented by a struct that
// wants to load something
type Loader interface {
	Load()
}

// Listener holds the methods for listening changes
type Listener interface {
	ListenToChanges(tracker Tracker)
}
