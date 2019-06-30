package listener

import (
	"crypto/tls"
	"net"
	"ngio/channel"
	"ngio/logger"
	"ngio/option"
	"strings"
)

type TCPListener struct {
	addr        *net.TCPAddr
	listener    *net.TCPListener
	opts        *option.Options
	log         logger.Logger
	initializer channel.Initializer
}

func NewTCPListener(network, laddr string, opts *option.Options, initializer channel.Initializer) (*TCPListener, error) {
	addr, err := net.ResolveTCPAddr(network, laddr)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, option.ErrOptionIsNil
	}

	return &TCPListener{
		addr:        addr,
		listener:    nil,
		opts:        opts,
		log:         logger.DefaultLogger(),
		initializer: initializer,
	}, nil
}

func (lsn *TCPListener) Serve() error {
	if lsn.addr == nil {
		return ErrBindAddrIsNil
	}

	listener, err := net.ListenTCP(lsn.addr.Network(), lsn.addr)
	if err != nil {
		return err
	}

	lsn.log.Infof("[network: %v, local: %v] listening", listener.Addr().Network(), listener.Addr())

	lsn.listener = listener

	for {
		conn, err := lsn.listener.AcceptTCP()
		if err != nil {
			// forwardly close will return err "use of closed network connection"
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			} else {
				lsn.log.Errorf("tcp accept\r\n %v", err)
				lsn.Shutdown()
				return err
			}
		}

		lsn.log.Debugf("[network: %v, remote: %v] accepted", conn.RemoteAddr().Network(), conn.RemoteAddr())

		if err = option.SetupTCPOptions(conn, lsn.opts); err != nil {
			lsn.log.Errorf("[network: %v, remote: %v] set socket option\r\n %v", conn.RemoteAddr().Network(), conn.RemoteAddr(), err)
			if closeErr := conn.Close(); closeErr != nil {
				lsn.log.Errorf("[network: %v, remote: %v] close\r\n %v", conn.RemoteAddr().Network(), conn.RemoteAddr(), closeErr)
			}
			continue
		}

		var ch *channel.TCPChannel
		if lsn.opts.TLSConfig != nil {
			ch = channel.NewTCPChannel(tls.Server(conn, lsn.opts.TLSConfig), lsn.opts.WriteDeadlinePeriod, lsn.opts.ReadDeadlinePeriod)
		} else {
			ch = channel.NewTCPChannel(conn, lsn.opts.WriteDeadlinePeriod, lsn.opts.ReadDeadlinePeriod)
		}

		if lsn.initializer != nil {
			lsn.initializer(ch)
		}

		go func() {
			_ = ch.Serve()
		}()
	}
}

func (lsn *TCPListener) Shutdown() {
	if lsn.listener == nil {
		return
	}

	lsn.log.Infof("[network: %v, local: %v] stop listening", lsn.listener.Addr().Network(), lsn.listener.Addr())

	// close listener and the serve loop will return
	if err := lsn.listener.Close(); err != nil {
		lsn.log.Errorf("[network: %v, local: %v] stop listening\r\n %v", lsn.listener.Addr().Network(), lsn.listener.Addr(), err)
	}

	lsn.log.Infof("[network: %v, local: %v] listen stopped", lsn.listener.Addr().Network(), lsn.listener.Addr())
}
