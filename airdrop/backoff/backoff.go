package backoff

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var globalBackoffOptions = []backoff.ExponentialBackOffOpts{
	func(b *backoff.ExponentialBackOff) {
		b.InitialInterval = 1 * time.Second
	},
	func(b *backoff.ExponentialBackOff) {
		b.MaxInterval = 32 * time.Second
	},
	func(b *backoff.ExponentialBackOff) {
		b.Multiplier = 2
	},
}

func NewBackoff(ctx context.Context) *backoff.ExponentialBackOff {
	return backoff.NewExponentialBackOff(globalBackoffOptions...)
}
