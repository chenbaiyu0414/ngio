package channel

import (
	"fmt"
	"net"
)

type Initializer func(ch Channel)

type Channel interface {
	Id() uint32
	IsActive() bool
	Pipeline() *Pipeline
	LocalAddress() net.Addr
	RemoteAddress() net.Addr
	Attributes() Attributes
	Serve() error
	Write(msg interface{})
	Close()
	fmt.Stringer
}
