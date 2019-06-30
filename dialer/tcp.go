package dialer

import (
	"crypto/tls"
	"net"
	"ngio/channel"
	"ngio/logger"
	"ngio/option"
)

type TCPDialer struct {
	laddr, raddr *net.TCPAddr
	ch           *channel.TCPChannel
	opts         *option.Options
	log          logger.Logger
	initializer  channel.Initializer
}

func NewTCPDialer(network, laddr, raddr string, opts *option.Options, initializer channel.Initializer) (*TCPDialer, error) {
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

	return &TCPDialer{
		laddr:       localAddr,
		raddr:       remoteAddr,
		opts:        opts,
		log:         logger.DefaultLogger(),
		initializer: initializer,
	}, nil
}

func (dal *TCPDialer) Dial() error {
	if dal.raddr == nil {
		return ErrDialAddrIsNil
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
		dal.ch = channel.NewTCPChannel(tls.Client(conn, dal.opts.TLSConfig), dal.opts.WriteDeadlinePeriod, dal.opts.ReadDeadlinePeriod)
	} else {
		dal.ch = channel.NewTCPChannel(conn, dal.opts.WriteDeadlinePeriod, dal.opts.ReadDeadlinePeriod)
	}

	if dal.initializer != nil {
		dal.initializer(dal.ch)
	}

	return dal.ch.Serve()
}

func (dal *TCPDialer) Close() {
	if dal.ch == nil {
		return
	}

	if !dal.ch.IsActive() {
		//dal.log.Warnf("[network: %v, local: %v, remote: %v] repeat close", dal.ch.RemoteAddress().Network(), dal.ch.LocalAddress(), dal.ch.RemoteAddress())
		return
	}

	dal.ch.Close()
}
