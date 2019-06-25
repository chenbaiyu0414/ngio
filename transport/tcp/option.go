package tcp

import (
	"ngio"
	"time"
)

type Option func(ch *channel) error

func WithKeepAlive(keepalive bool) Option {
	return func(ch *channel) error {
		return ch.rwc.SetKeepAlive(keepalive)
	}
}

func WithKeepAlivePeriod(d time.Duration) Option {
	return func(ch *channel) error {
		return ch.rwc.SetKeepAlivePeriod(d)
	}
}

func WithNoDelay(noDelay bool) Option {
	return func(ch *channel) error {
		return ch.rwc.SetNoDelay(noDelay)
	}
}

func WithLinger(sec int) Option {
	return func(ch *channel) error {
		return ch.rwc.SetLinger(sec)
	}
}

func WithWriteBufferSize(size int) Option {
	return func(ch *channel) error {
		return ch.rwc.SetWriteBuffer(size)
	}
}

func WithReadBufferSize(size int) Option {
	return func(ch *channel) error {
		return ch.rwc.SetReadBuffer(size)
	}
}

func WithDeadline(t time.Time) Option {
	return func(ch *channel) error {
		return ch.rwc.SetDeadline(t)
	}
}

func WithReadDeadline(t time.Time) Option {
	return func(ch *channel) error {
		return ch.rwc.SetReadDeadline(t)
	}
}

func WithWriteDeadline(t time.Time) Option {
	return func(ch *channel) error {
		return ch.rwc.SetWriteDeadline(t)
	}
}

func WithInitializer(initializer ngio.Initializer) Option {
	return func(ch *channel) error {
		initializer(ch)
		return nil
	}
}

func applyOptions(opts []Option, ch *channel) error {
	for _, opt := range opts {
		if err := opt(ch); err != nil {
			return err
		}
	}

	return nil
}
