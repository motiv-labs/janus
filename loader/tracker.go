package loader

// Tracker holds the channels to keep track of changes that will require a reload of actions
type Tracker struct {
	stopChan chan bool
	changed  chan bool
}

// NewTracker creates a new instance of Tracker
func NewTracker() *Tracker {
	return &Tracker{make(chan bool), make(chan bool)}
}

// Change tells the goroutine that something changed.
func (t *Tracker) Change() {
	t.changed <- true
}

// Changed gets the stop channel for this tracker.
// Reading from this channel will block while the task is running, and will
// unblock once the task has stopped (because the channel gets closed).
func (t *Tracker) Changed() <-chan bool {
	return t.changed
}

// Stop tells the goroutine to stop.
func (t *Tracker) Stop() {
	t.stopChan <- true
	close(t.stopChan)
	close(t.changed)
}

// StopChan gets the stop channel for this tracker.
// Reading from this channel will block while the task is running, and will
// unblock once the task has stopped (because the channel gets closed).
func (t *Tracker) StopChan() <-chan bool {
	return t.stopChan
}
