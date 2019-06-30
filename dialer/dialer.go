package dialer

import "errors"

var (
	ErrDialAddrIsNil = errors.New("dialer remote addr is nil")
)

type Dialer interface {
	Dial() error
	Close()
}
