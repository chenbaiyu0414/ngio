package transport

import "io"

type Client interface {
	io.Closer
}
