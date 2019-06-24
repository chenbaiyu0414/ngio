package tcp

import (
	"net"
	"ngio/internal/logger"
)

func Dial(network, localAddr, remoteAddr string, opts ...Option) (clt *Client, err error) {
	raddr, err := net.ResolveTCPAddr(network, remoteAddr)
	if err != nil {
		return
	}

	laddr, err := net.ResolveTCPAddr(network, localAddr)
	if localAddr != "" && err != nil {
		return
	}

	rwc, err := net.DialTCP(network, laddr, raddr)
	if err != nil {
		return
	}

	clt = &Client{
		ch: newChannel(rwc),
	}

	if err = applyOptions(opts, clt.ch); err != nil {
		logger.Errorf("apply options to channel failed: %v", err)
		_ = clt.ch.Close()
		return
	}

	go func() {
		<-clt.ch.Serve()
	}()

	return
}

type Client struct {
	ch *tcpChannel
}

func (clt *Client) Close() error {
	return clt.ch.Close()
}
