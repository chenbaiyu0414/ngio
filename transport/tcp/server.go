package tcp

import (
	"net"
	"ngio/channel"
	"ngio/internal/logger"
	"sync"
)

type Server struct {
	closeC   chan struct{}
	lis      *net.TCPListener
	channels sync.Map
	wg       sync.WaitGroup
	opts     []Option
}

func NewServer(opts ...Option) *Server {
	return &Server{
		opts:   opts,
		closeC: make(chan struct{}, 1),
	}
}

func (srv *Server) Serve(network, localAddr string) error {
	laddr, err := net.ResolveTCPAddr(network, localAddr)
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
		rwc, err := lis.AcceptTCP()
		if err != nil {
			select {
			case <-srv.closeC:
				return nil
			default:
				return err
			}
		}

		ch := newChannel(rwc)

		if err := applyOptions(srv.opts, ch); err != nil {
			logger.Errorf("apply options to channel failed: %v", err)
			_ = ch.Close()
			continue
		}

		go func() {
			srv.wg.Add(1)
			defer srv.wg.Done()

			srv.channels.Store(ch, struct{}{})
			defer srv.channels.Delete(ch)

			if err := <-ch.Serve(); err != nil {
				logger.Errorf("close channel: %v", err)
			}
		}()
	}
}

func (srv *Server) Shutdown() {
	srv.closeC <- struct{}{}

	if err := srv.lis.Close(); err != nil {
		logger.Errorf("close server listener: %v", err)
	}

	srv.channels.Range(func(key, value interface{}) bool {
		_ = key.(channel.Channel).Close()
		return true
	})

	srv.wg.Wait()

	close(srv.closeC)
}
