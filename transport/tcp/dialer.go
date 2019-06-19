package tcp

import (
	"net"
	"ngio/channel"
)

func Dial(network, localAddr, remoteAddr string, opt ...Option) (ch channel.Channel, err error) {
	raddr, err := net.ResolveTCPAddr(network, remoteAddr)
	if err != nil {
		return nil, err
	}

	laddr, err := net.ResolveTCPAddr(network, localAddr)
	if err != nil && localAddr != "" {
		return nil, err
	}

	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	tcpConn, err := net.DialTCP(network, laddr, raddr)
	if err != nil {
		return nil, err
	}

	applyOptions(&opts, tcpConn)

	ch = newChannel(tcpConn)

	opts.initializer(ch)

	return
}
