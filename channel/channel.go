package channel

import (
	"net"
)

type Initializer func(ch Channel)

type Channel interface {
	IsActive() bool
	Pipeline() *Pipeline
	LocalAddress() net.Addr
	RemoteAddress() net.Addr
	Attributes() Attributes
	Serve() <-chan struct{}
	Write(msg interface{})
	Close() error
}
