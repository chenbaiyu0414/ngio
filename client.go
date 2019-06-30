package ngio

import (
	"ngio/channel"
	"ngio/dialer"
	"ngio/option"
)

type Client struct {
	network, laddr, raddr string
	dialer                dialer.Dialer
	opts                  *option.Options
	initializer           channel.Initializer
}

func NewClient(network, laddr, raddr string) *Client {
	return &Client{
		network:     network,
		laddr:       laddr,
		raddr:       raddr,
		dialer:      nil,
		opts:        new(option.Options),
		initializer: nil,
	}
}

func (clt *Client) Option(opts ...option.Option) *Client {
	for _, o := range opts {
		o.Apply(clt.opts)
	}

	return clt
}

func (clt *Client) Channel(initializer channel.Initializer) *Client {
	clt.initializer = initializer
	return clt
}

func (clt *Client) Dial() (err error) {
	switch clt.network {
	case "tcp", "tcp4", "tcp6":
		clt.dialer, err = dialer.NewTCPDialer(clt.network, clt.laddr, clt.raddr, clt.opts, clt.initializer)
	case "udp", "udp4", "udp6":
		clt.dialer, err = dialer.NewUDPDialer(clt.network, clt.laddr, clt.raddr, clt.opts, clt.initializer)
	//case "ip", "ip4", "ip6":
	default:
		err = ErrUnsupportedNetwork
	}

	if err != nil {
		return
	}

	return clt.dialer.Dial()
}

func (clt *Client) Close() {
	clt.dialer.Close()
}
