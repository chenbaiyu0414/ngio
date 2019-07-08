package transport

import "errors"

var (
	ErrBindAddrIsNil = errors.New("listener local addr is nil")
)

type Listener interface {
	Serve() error
	Shutdown()
}
