package flow

import (
	"errors"
	"testing"
)

func TestFlow_Subscribe(t *testing.T) {
	var tests = []struct {
		a, b int
		want int
		ok bool
	} {
		{4, 2, 2, true},
		{10, 0, 0, false},
		{10, 5, 2, true},
		{0, 5, 0, true},
	}

	for _, tc := range tests {
		div(tc.a, tc.b).Subscribe(
			func(data interface{}) {
				if data.(int) != tc.want {
					t.Errorf("div(%d, %d) = %d, want = %d", tc.a, tc.b, data, tc.want)
				}
			},
			func(err error) {
				if !errors.Is(err, ErrDivisionByZero) {
					t.Errorf("div(%d, %d) = %s, want = %s", tc.a, tc.b, err, ErrDivisionByZero)
				}
			},
			func(ok bool) {
				if ok != tc.ok {
					t.Errorf("div(%d, %d) = %t, want = %t", tc.a, tc.b, ok, tc.ok)
				}
			})
	}
}

func TestFlow_SubscribeOnNext(t *testing.T) {
	a, b, want := 4, 2, 2

	div(a, b).SubscribeOnNext(
		func(data interface{}) {
			if data.(int) != 2 {
				t.Errorf("div(%d, %d) = %v, want = %d", a, b, data, want)
			}
		})
}

func TestFlow_SubscribeOnError(t *testing.T) {
	a, b := 4, 0

	div(a, b).SubscribeOnError(
		func(err error) {
			if !errors.Is(err, ErrDivisionByZero) {
				t.Errorf("div(%d, %d) = %s, want = %s", a, b, err, ErrDivisionByZero)
			}
		})
}

func TestFlow_SubscribeOnComplete(t *testing.T) {
	a, b := 4, 2

	div(a, b).SubscribeOnComplete(
		func(ok bool) {
			if !ok {
				t.Errorf("div(%d, %d) = %t, want = %t", a, b, ok, true)
			}
		})
}

func TestFlow_SubscribeParallel(t *testing.T) {
	var tests = []struct {
		name string
		a, b int
		want int
		ok bool
	} {
		{"test_1", 4, 2, 2, true},
		{"test_2", 10, 0, 0, false},
		{"test_3", 10, 5, 2, true},
		{"test_4", 0, 5, 0, true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			div(tc.a, tc.b).Subscribe(
				func(data interface{}) {
					if data.(int) != tc.want {
						t.Errorf("div(%d, %d) = %d, want = %d", tc.a, tc.b, data, tc.want)
					}
				},
				func(err error) {
					if !errors.Is(err, ErrDivisionByZero) {
						t.Errorf("div(%d, %d) = %s, want = %s", tc.a, tc.b, err, ErrDivisionByZero)
					}
				},
				func(ok bool) {
					if ok != tc.ok {
						t.Errorf("div(%d, %d) = %t, want = %t", tc.a, tc.b, ok, tc.ok)
					}
				})
		})
	}
}

var ErrDivisionByZero = errors.New("division by zero is not defined")

func div(a, b int) *Flow {
	return New(func(emitter Emitter) {
		if b == 0 {
			emitter.OnError(ErrDivisionByZero)
			return
		}

		emitter.OnNext(a / b)
		emitter.OnComplete()
	})
}