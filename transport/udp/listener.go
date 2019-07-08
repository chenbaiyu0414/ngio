package udp

import (
	"net"
	"ngio/internal"
	"ngio/logger"
	"ngio/option"
	"ngio/transport"
)

type listener struct {
	addr        *net.UDPAddr
	ch          *channel
	opts        *option.Options
	log         logger.Logger
	initializer func(channel internal.IChannel)
}

func NewListener(network, laddr string, opts *option.Options, initializer func(channel internal.IChannel)) (*listener, error) {
	addr, err := net.ResolveUDPAddr(network, laddr)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, option.ErrOptionIsNil
	}

	return &listener{
		addr:        addr,
		ch:          nil,
		opts:        opts,
		log:         logger.DefaultLogger(),
		initializer: initializer,
	}, nil
}

func (lsn *listener) Serve() error {
	if lsn.addr == nil {
		return transport.ErrBindAddrIsNil
	}

	conn, err := net.ListenUDP(lsn.addr.Network(), lsn.addr)
	if err != nil {
		return err
	}

	lsn.log.Infof("[network: %v, local: %v] listening", conn.LocalAddr().Network(), conn.LocalAddr())

	if err := option.SetupUDPOptions(conn, lsn.opts); err != nil {
		lsn.log.Errorf("[network: %v, local: %v] set socket option\r\n %v", conn.LocalAddr().Network(), conn.LocalAddr(), err)
		if closeErr := conn.Close(); closeErr != nil {
			lsn.log.Errorf("[network: %v, local: %v] close\r\n %v", conn.LocalAddr().Network(), conn.LocalAddr(), closeErr)
		}
		return err
	}

	lsn.ch = newChannel(conn)
	if lsn.initializer != nil {
		lsn.initializer(lsn.ch)
	}

	return lsn.ch.Serve()
}

func (lsn *listener) Shutdown() {
	if lsn.ch == nil {
		return
	}

	if !lsn.ch.IsActive() {
		lsn.log.Warn("close udp listener repeated")
		return
	}

	lsn.log.Infof("[network: %v, local: %v] stop listening", lsn.ch.LocalAddress().Network(), lsn.ch.LocalAddress())

	lsn.ch.Close()

	lsn.log.Infof("[network: %v, local: %v] listen stopped", lsn.ch.LocalAddress().Network(), lsn.ch.LocalAddress())

}
