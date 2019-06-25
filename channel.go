package ngio

import (
	"net"
)

type Initializer func(Channel)

type Channel interface {
	IsActive() bool
	Pipeline() *Pipeline
	LocalAddress() net.Addr
	RemoteAddress() net.Addr
	Attributes() Attributes
	Serve() <-chan error
	Write(msg interface{})
	Close()
}
