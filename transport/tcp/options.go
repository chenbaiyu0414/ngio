package tcp

import (
	"ngio/channel"
	"time"
)

type Option interface {
	apply(options *tcpOptions)
}

type tcpOptions struct {
	keepalive       bool
	keepalivePeriod time.Duration
	noDelay         bool
	linger          int
	writeBufferSize int
	readBufferSize  int
	deadline        time.Time
	readDeadline    time.Time
	writeDeadline   time.Time
	initializer     channel.Initializer
}

type tcpOptionApplier struct {
	f func(options *tcpOptions)
}

func newTcpOptionApplier(f func(options *tcpOptions)) *tcpOptionApplier {
	return &tcpOptionApplier{f: f}
}

func (toa *tcpOptionApplier) apply(options *tcpOptions) {
	toa.f(options)
}

func WithKeepAlive(keepalive bool) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.keepalive = keepalive
	})
}

func WithKeepAlivePeriod(d time.Duration) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.keepalivePeriod = d
	})
}

func WithNoDelay(noDelay bool) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.noDelay = noDelay
	})
}

func WithLinger(sec int) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.linger = sec
	})
}

func WithWriteBufferSize(size int) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.writeBufferSize = size
	})
}

func WithReadBufferSize(size int) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.readBufferSize = size
	})
}

func WithDeadline(t time.Time) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.deadline = t
	})
}

func WithReadDeadline(t time.Time) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.readDeadline = t
	})
}

func WithWriteDeadline(t time.Time) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.writeDeadline = t
	})
}

func WithInitializer(initializer channel.Initializer) Option {
	return newTcpOptionApplier(func(options *tcpOptions) {
		options.initializer = initializer
	})
}
