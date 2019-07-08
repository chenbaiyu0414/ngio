package ngio

import (
	"ngio/internal"
	"ngio/option"
	"ngio/transport"
	"ngio/transport/tcp"
	"ngio/transport/udp"
)

type Client struct {
	network, laddr, raddr string
	dialer                transport.Dialer
	opts                  *option.Options
	initializer           func(channel Channel)
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

func (clt *Client) Channel(initializer func(channel Channel)) *Client {
	clt.initializer = initializer
	return clt
}

func (clt *Client) Dial() (err error) {
	switch clt.network {
	case "tcp", "tcp4", "tcp6":
		clt.dialer, err = tcp.NewDialer(clt.network, clt.laddr, clt.raddr, clt.opts, func(channel internal.IChannel) {
			clt.initializer(channel)
		})
	case "udp", "udp4", "udp6":
		clt.dialer, err = udp.NewDialer(clt.network, clt.laddr, clt.raddr, clt.opts, func(channel internal.IChannel) {
			clt.initializer(channel)
		})
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
