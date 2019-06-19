package tcp

import (
	"net"
	"ngio/channel"
	"ngio/internal/logger"
	"sync"
)

var defaultServerOptions = tcpOptions{}

type Server struct {
	opts          *tcpOptions
	closeC        chan struct{}
	activeChannel map[channel.Channel]struct{}
	lis           *net.TCPListener
	mu            sync.Mutex
	wg            sync.WaitGroup
}

func NewServer(opt ...Option) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Server{
		opts:          &opts,
		closeC:        make(chan struct{}, 1),
		activeChannel: make(map[channel.Channel]struct{}),
	}
}

func (srv *Server) Serve(network, addr string) error {
	laddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return err
	}

	lis, err := net.ListenTCP(network, laddr)
	if err != nil {
		return err
	}

	srv.lis = lis
	logger.Infof("server listening on %v", lis.Addr())

	for {
		tcpConn, err := lis.AcceptTCP()
		if err != nil {
			select {
			case <-srv.closeC:
				return nil
			default:
				return err
			}
		}

		applyOptions(srv.opts, tcpConn)

		ch := newChannel(tcpConn)

		srv.opts.initializer(ch)

		srv.wg.Add(1)
		go func() {
			// ch.Serve will blocked until ch closed
			<-ch.Serve()

			// delete ch from channel map
			srv.mu.Lock()
			delete(srv.activeChannel, ch)
			srv.mu.Unlock()

			srv.wg.Done()
		}()
	}

}

func (srv *Server) Shutdown() {
	srv.closeC <- struct{}{}
	srv.lis.Close()

	go func() {
		srv.mu.Lock()
		for ch := range srv.activeChannel {
			ch.Close()
		}
		srv.mu.Unlock()
	}()

	srv.wg.Wait()

	close(srv.closeC)
}

func applyOptions(options *tcpOptions, conn *net.TCPConn) {
	conn.SetKeepAlive(options.keepalive)
	conn.SetKeepAlivePeriod(options.keepalivePeriod)
	conn.SetNoDelay(options.noDelay)
	conn.SetLinger(options.linger)
	conn.SetWriteBuffer(options.writeBufferSize)
	conn.SetReadBuffer(options.readBufferSize)
	conn.SetDeadline(options.deadline)
	conn.SetWriteDeadline(options.writeDeadline)
	conn.SetReadDeadline(options.readDeadline)
}
