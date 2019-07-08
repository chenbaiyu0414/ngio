package echo

import (
	"ngio"
	"ngio/buffer"
	"ngio/logger"
)

type Handler struct {
	log logger.Logger
}

func NewHandler() *Handler {
	return &Handler{
		log: logger.DefaultLogger(),
	}
}

func (handler *Handler) ChannelRead(ctx ngio.ChannelContext, msg interface{}) {
	bf, ok := msg.(buffer.ByteBuffer)
	if !ok {
		handler.log.Errorf("msg is not buffer.ByteBuffer")
		return
	}

	received := bf.GetBytes(bf.ReaderIndex(), bf.ReadableBytes())

	handler.log.Infof("received: %s", string(received))

	ctx.Write(bf)
}

func (handler *Handler) ChannelInActive(ctx ngio.ChannelContext) {
	handler.log.Infof("inactive")
}

func (handler *Handler) ChannelActive(ctx ngio.ChannelContext) {
	handler.log.Infof("active")
}

func (handler *Handler) HandleError(ctx ngio.ChannelContext, err error) {
	handler.log.Errorf("unexpected unhandled error: %v", err)
}
