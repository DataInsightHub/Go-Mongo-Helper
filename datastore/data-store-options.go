package datastore

import (
	"time"
)

type (
	DataStoreOptions interface {
		apply(*dataStoreOption)
	}
)

type (
	dataStoreOption struct {
		timeout time.Duration
		usePing bool
	}
)

type timeoutOption time.Duration

func (value timeoutOption) apply(o *dataStoreOption) {
	duration := time.Duration(value)
	if duration <= 0 {
		return
	}
	o.timeout = duration
}

func WithTimeoutOption(duration time.Duration) DataStoreOptions {
	return timeoutOption(duration)
}

type usePingOption bool

func (value usePingOption) apply(o *dataStoreOption) {
	o.usePing = bool(value)
}

func WithUsePingOption(usePing bool) DataStoreOptions {
	return usePingOption(usePing)
}
