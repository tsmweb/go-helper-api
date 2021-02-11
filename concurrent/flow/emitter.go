package flow

// Emitter push signals to subscriber.
type Emitter struct {
	data     chan interface{}
	err      chan error
	complete chan bool
}

func newEmitter() Emitter {
	return Emitter{
		data:     make(chan interface{}),
		err:      make(chan error),
		complete: make(chan bool),
	}
}

// OnNext send data signals to subscriber.
func (e Emitter) OnNext(data interface{}) {
	e.data <- data
}

// OnError send errors signals to subscriber.
func (e Emitter) OnError(err error) {
	e.err <- err
	e.onComplete(false)
}

func (e Emitter) onComplete(success bool) {
	e.complete <- success
}

// OnComplete send the completed signal to subscriber.
func (e Emitter) OnComplete() {
	e.onComplete(true)
}
