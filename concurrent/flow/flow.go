/*
Package flow provides a Flow implementation to perform background processing
and notify your subscribers through an Emitter.

Create a new instance of Flow, example:

	flw := flow.New(func(emitter concurrent.Emitter) {
		// ...
		result, err := repository.Get(id)
		if err != nil {
			emitter.OnError(err)
			return
		}

		emitter.OnNext(result)
		emitter.OnComplete()
	})

Subscribing to a Flow (Safe Concurrency), example:

	flw.SubscribeOutboxEvent(
		func(data interface) {
			// OnNext
			// ...
		},
		func(err error) {
			// OnError
			// ...
		},
		func(ok bool) {
			// OnComplete
			if ok {
				// ...
		})

*/
package flow

// Flow performs background processing and notifies your subscribers via an Emitter.
type Flow struct {
	subscribe func(emitter Emitter)
}

// New create new instance of Flow.
func New(onSubscribe func(emitter Emitter)) *Flow {
	return &Flow{ subscribe: onSubscribe }
}

// SubscribeOutboxEvent records callbacks for onNext, onError and onComplete.
// When subscribing to a Flow, processing is performed in the background and callbacks are notified via signals.
// Safe Concurrency.
func (f *Flow) Subscribe(onNext func(data interface{}), onError func(err error), onComplete func(ok bool)) {
	emitter := newEmitter()
	done := make(chan struct{})

	go func(complete chan<- struct{}) {
		loop:
		for {
			select {
			case data, ok := <-emitter.data:
				if ok {
					onNext(data)
				}
			case err, ok := <-emitter.err:
				if ok {
					onError(err)
				}
			case ok := <-emitter.complete:
				close(emitter.data)
				close(emitter.err)

				onComplete(ok)
				complete <- struct{}{}
				break loop
			}
		}
	}(done)

	go func() {
		f.subscribe(emitter)
	}()

	<-done
}

// SubscribeOnNext registers callbacks for onNext.
// When subscribing to a Flow, processing is performed in the background and callbacks are notified via signals.
func (f *Flow) SubscribeOnNext(onNext func(data interface{})) {
	f.Subscribe(onNext, func(err error) {}, func(ok bool) {})
}

// SubscribeOnError registers callbacks for onError.
// When subscribing to a Flow, processing is performed in the background and callbacks are notified via signals.
func (f *Flow) SubscribeOnError(onError func(err error)) {
	f.Subscribe(func(data interface{}) {}, onError, func(ok bool) {})
}

// SubscribeOnComplete registers callbacks for onComplete.
// When subscribing to a Flow, processing is performed in the background and callbacks are notified via signals.
func (f *Flow) SubscribeOnComplete(onComplete func(ok bool)) {
	f.Subscribe(func(data interface{}) {}, func(err error) {}, onComplete)
}
