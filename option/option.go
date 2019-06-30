package option

import (
	"crypto/tls"
	"errors"
	"time"
)

var (
	ErrOptionIsNil = errors.New("option is nil")
)

type Options struct {
	TCPKeepAlive        bool
	TCPKeepAlivePeriod  time.Duration
	TCPNoDelay          bool
	TCPLinger           int
	WriteBuffer         int
	ReadBuffer          int
	ReadDeadlinePeriod  time.Duration
	WriteDeadlinePeriod time.Duration
	TLSConfig           *tls.Config
}

type Option interface {
	Apply(*Options)
}

type optionFunc struct {
	f func(*Options)
}

func newOptionFunc(f func(*Options)) *optionFunc {
	return &optionFunc{f: f}
}

func (ofn *optionFunc) Apply(o *Options) {
	ofn.f(o)
}

func TCPKeepAlive(keepalive bool) Option {
	return newOptionFunc(func(o *Options) {
		o.TCPKeepAlive = keepalive
	})
}

func TCPKeepAlivePeriod(d time.Duration) Option {
	return newOptionFunc(func(o *Options) {
		o.TCPKeepAlivePeriod = d
	})
}

func TCPNoDelay(noDelay bool) Option {
	return newOptionFunc(func(o *Options) {
		o.TCPNoDelay = noDelay
	})
}

func TCPLinger(linger int) Option {
	return newOptionFunc(func(o *Options) {
		o.TCPLinger = linger
	})
}

func ReadBuffer(size int) Option {
	return newOptionFunc(func(o *Options) {
		o.ReadBuffer = size
	})
}

func WriteBuffer(size int) Option {
	return newOptionFunc(func(o *Options) {
		o.WriteBuffer = size
	})
}

func WriteDeadlinePeriod(d time.Duration) Option {
	return newOptionFunc(func(o *Options) {
		o.WriteDeadlinePeriod = d
	})
}

func ReadDeadlinePeriod(d time.Duration) Option {
	return newOptionFunc(func(o *Options) {
		o.ReadDeadlinePeriod = d
	})
}

func TLS(tlsConfig *tls.Config) Option {
	return newOptionFunc(func(o *Options) {
		o.TLSConfig = tlsConfig
	})
}
