package ngio

import (
	"errors"
	"ngio/internal"
	"ngio/option"
	"ngio/transport"
	"ngio/transport/tcp"
	"ngio/transport/udp"
)

var (
	ErrUnsupportedNetwork = errors.New("unsupported network")
)

type Server struct {
	network, laddr string
	lsn            transport.Listener
	opts           *option.Options
	initializer    func(channel Channel)
}

func NewServer(network, laddr string) *Server {
	defaultOptions := &option.Options{
		TCPNoDelay: true, // tcp nodelay is true by default. see src/net/tcpsock.newTCPConn:195
		TCPLinger:  -1,   // tcp linger < 0 by default. see net.TCPConn:SetLinger() comments.
	}

	return &Server{
		network:     network,
		laddr:       laddr,
		lsn:         nil,
		opts:        defaultOptions,
		initializer: nil,
	}
}

func (srv *Server) Option(opts ...option.Option) *Server {
	for _, o := range opts {
		o.Apply(srv.opts)
	}

	return srv
}

func (srv *Server) Channel(initializer func(channel Channel)) *Server {
	srv.initializer = initializer
	return srv
}

func (srv *Server) Serve() (err error) {
	switch srv.network {
	case "tcp", "tcp4", "tcp6":
		srv.lsn, err = tcp.NewListener(srv.network, srv.laddr, srv.opts, func(channel internal.IChannel) {
			srv.initializer(channel)
		})
	case "udp", "udp4", "udp6":
		srv.lsn, err = udp.NewListener(srv.network, srv.laddr, srv.opts, func(channel internal.IChannel) {
			srv.initializer(channel)
		})
	//case "ip", "ip4", "ip6":
	default:
		err = ErrUnsupportedNetwork
	}

	if err != nil {
		return
	}

	return srv.lsn.Serve()
}

func (srv *Server) Shutdown() {
	srv.lsn.Shutdown()
}
