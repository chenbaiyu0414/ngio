package tcp

import (
	"crypto/tls"
	"net"
	"ngio/internal"
	"ngio/logger"
	"ngio/option"
	"ngio/transport"
)

type dialer struct {
	laddr, raddr *net.TCPAddr
	ch           *channel
	opts         *option.Options
	log          logger.Logger
	initializer  func(channel internal.IChannel)
}

func NewDialer(network, laddr, raddr string, opts *option.Options, initializer func(channel internal.IChannel)) (*dialer, error) {
	remoteAddr, err := net.ResolveTCPAddr(network, raddr)
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveTCPAddr(network, laddr)
	if laddr != "" && err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, option.ErrOptionIsNil
	}

	return &dialer{
		laddr:       localAddr,
		raddr:       remoteAddr,
		opts:        opts,
		log:         logger.DefaultLogger(),
		initializer: initializer,
	}, nil
}

func (dal *dialer) Dial() error {
	if dal.raddr == nil {
		return transport.ErrDialAddrIsNil
	}

	conn, err := net.DialTCP(dal.raddr.Network(), dal.laddr, dal.raddr)
	if err != nil {
		return err
	}

	dal.log.Infof("[network: %v, local: %v, remote: %v] dialed", conn.RemoteAddr().Network(), conn.LocalAddr(), conn.RemoteAddr())

	if err := option.SetupTCPOptions(conn, dal.opts); err != nil {
		dal.log.Errorf("[network: %v, local: %v, remote: %v] set socket option\r\n %v", conn.RemoteAddr().Network(), conn.LocalAddr(), conn.RemoteAddr(), err)
		if closeErr := conn.Close(); closeErr != nil {
			dal.log.Errorf("[network: %v, local: %v, remote: %v] close\r\n %v", conn.RemoteAddr().Network(), conn.LocalAddr(), conn.RemoteAddr(), closeErr)
		}
		return err
	}

	if dal.opts.TLSConfig != nil {
		dal.ch = newChannel(tls.Client(conn, dal.opts.TLSConfig), dal.opts.WriteDeadlinePeriod, dal.opts.ReadDeadlinePeriod)
	} else {
		dal.ch = newChannel(conn, dal.opts.WriteDeadlinePeriod, dal.opts.ReadDeadlinePeriod)
	}

	if dal.initializer != nil {
		dal.initializer(dal.ch)
	}

	return dal.ch.Serve()
}

func (dal *dialer) Close() {
	if dal.ch == nil {
		return
	}

	if !dal.ch.IsActive() {
		//dal.log.Warnf("[network: %v, local: %v, remote: %v] repeat close", dal.ch.RemoteAddress().Network(), dal.ch.LocalAddress(), dal.ch.RemoteAddress())
		return
	}

	dal.ch.Close()
}
