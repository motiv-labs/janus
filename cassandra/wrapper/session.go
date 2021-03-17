package wrapper

// Initializer is a common interface for functionality to start a new session
type Initializer interface {
	NewSession() (Holder, error)
}

// Holder allows to store a close sessions
type Holder interface {
	GetSession() SessionInterface
	CloseSession()
}

// SessionInterface is an interface to wrap gocql methods used in Motiv
type SessionInterface interface {
	Query( stmt string, values ...interface{}) QueryInterface
	Close()
}

type QueryInterface interface {
	Exec() error
	Scan( dest ...interface{}) error
	Iter() IterInterface
	PageState(state []byte, ) QueryInterface
	PageSize(n int, ) QueryInterface
}

type IterInterface interface {
	Scan( dest ...interface{}) bool
	WillSwitchPage() bool
	PageState() []byte
	Close() error
	ScanAndClose( handle func() bool, dest ...interface{}) error
}
