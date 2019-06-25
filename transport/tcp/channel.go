package tcp

import (
	"io"
	"net"
	"ngio"
	"ngio/buffer"
	"sync"
)

// channel is a connection between server and client
type channel struct {
	isActive      bool
	rwc           *net.TCPConn
	closeC        chan struct{}
	quitC         chan error
	writeC        chan buffer.ByteBuffer
	wg            sync.WaitGroup
	pipeline      *ngio.Pipeline
	recvAllocator *buffer.RecvByteBufAllocator
	attributes    ngio.Attributes
}

func newChannel(rwc *net.TCPConn) *channel {
	ch := &channel{
		rwc:           rwc,
		closeC:        make(chan struct{}),
		quitC:         make(chan error, 1),
		writeC:        make(chan buffer.ByteBuffer, 16),
		recvAllocator: buffer.NewRecvByteBufAllocator(buffer.DefaultMinimum, buffer.DefaultMaximum, buffer.DefaultInitial),
		attributes:    ngio.NewDefaultAttributes(),
	}

	ch.pipeline = ngio.NewPipeline(ch)
	return ch
}

func (ch *channel) IsActive() bool {
	return ch.isActive
}

func (ch *channel) Pipeline() *ngio.Pipeline {
	return ch.pipeline
}

func (ch *channel) LocalAddress() net.Addr {
	if ch.rwc == nil {
		return nil
	}

	return ch.rwc.LocalAddr()
}

func (ch *channel) RemoteAddress() net.Addr {
	if ch.rwc == nil {
		return nil
	}

	return ch.rwc.RemoteAddr()
}

func (ch *channel) Attributes() ngio.Attributes {
	return ch.attributes
}

func (ch *channel) Serve() <-chan error {
	go ch.read()
	go ch.write()

	ch.isActive = true
	ch.pipeline.FireActiveHandler()

	return ch.quitC
}

func (ch *channel) read() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.closeC:
			return
		default:
			buf := ch.recvAllocator.Allocate()
			n, err := buf.ReadFrom(ch.rwc)

			if err == nil {
				ch.recvAllocator.Record(n)
				ch.pipeline.FireReadHandler(buf)
				continue
			}

			if err == io.EOF || err == io.ErrUnexpectedEOF {
				ch.Close()
			}

			return
		}
	}
}

func (ch *channel) write() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.closeC:
			return
		default:
			for buf := range ch.writeC {
				unwritten := int64(buf.ReadableBytes())

				for unwritten > 0 {
					n, err := buf.WriteTo(ch.rwc)

					if err == nil {
						unwritten -= n
						continue
					}

					if err == io.EOF || err == io.ErrUnexpectedEOF {
						ch.Close()
					}

					return
				}
			}
		}
	}
}

func (ch *channel) Write(msg interface{}) {
	if !ch.isActive {
		return
	}

	if buf, ok := msg.(buffer.ByteBuffer); ok {
		ch.writeC <- buf
	}
}

func (ch *channel) Close() {
	if !ch.isActive {
		return
	}

	ch.isActive = false
	ch.pipeline.FireInActiveHandler()

	go func() {
		// broadcast close signal
		close(ch.closeC)
		close(ch.writeC)

		err := ch.rwc.Close()

		ch.wg.Wait()

		ch.quitC <- err
	}()
}
