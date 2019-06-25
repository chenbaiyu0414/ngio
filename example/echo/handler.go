package echo

import (
	"ngio"
	"ngio/buffer"
	"ngio/internal/logger"
)

type Handler struct {
}

func (*Handler) ChannelRead(ctx ngio.Context, msg interface{}) {
	bf, ok := msg.(buffer.ByteBuffer)
	if !ok {
		logger.Errorf("msg is not buffer.ByteBuffer")
		return
	}

	received := bf.GetBytes(bf.ReaderIndex(), bf.ReadableBytes())

	logger.Infof("received: %s", string(received))

	ctx.Write(bf)
}

func (*Handler) ChannelInActive(ctx ngio.Context) {
	logger.Infof("inactive")
}

func (*Handler) ChannelActive(ctx ngio.Context) {
	logger.Infof("active")
}
