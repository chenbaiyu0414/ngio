package option

import "net"

func SetupTCPOptions(conn *net.TCPConn, opts *Options) (err error) {
	if opts.TCPKeepAlive == true {
		if err = conn.SetKeepAlive(true); err != nil {
			return
		}
	}

	if opts.TCPKeepAlivePeriod != 0 {
		if err = conn.SetKeepAlivePeriod(opts.TCPKeepAlivePeriod); err != nil {
			return
		}
	}

	if !opts.TCPNoDelay {
		// tcp nodelay is true by default. see src/net/tcpsock.newTCPConn:195
		if err = conn.SetNoDelay(opts.TCPNoDelay); err != nil {
			return
		}
	}

	if opts.TCPLinger >= 0 {
		if err = conn.SetLinger(opts.TCPLinger); err != nil {
			return
		}
	}

	if opts.ReadBuffer > 0 {
		if err = conn.SetReadBuffer(opts.ReadBuffer); err != nil {
			return
		}
	}

	if opts.WriteBuffer > 0 {
		if err = conn.SetWriteBuffer(opts.WriteBuffer); err != nil {
			return
		}
	}

	return
}

func SetupUDPOptions(conn *net.UDPConn, opts *Options) (err error) {
	if opts.ReadBuffer > 0 {
		if err = conn.SetReadBuffer(opts.ReadBuffer); err != nil {
			return
		}
	}

	if opts.WriteBuffer > 0 {
		if err = conn.SetWriteBuffer(opts.WriteBuffer); err != nil {
			return
		}
	}

	return
}
