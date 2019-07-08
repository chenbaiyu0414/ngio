package internal

import (
	"ngio/logger"
)

type ChannelHandler struct {
	name            string
	activeHandler   ActiveHandler
	inActiveHandler InActiveHandler
	readHandler     ReadHandler
	writeHandler    WriteHandler
	errorHandler    ErrorHandler
	flag            Flag
	log             logger.Logger
}

func NewChannelHandler(name string, handler interface{}) *ChannelHandler {
	adapter := &ChannelHandler{
		name: name,
		flag: None,
		log:  logger.DefaultLogger(),
	}

	if h, ok := handler.(ActiveHandler); ok {
		adapter.activeHandler = h
		adapter.flag |= Active
	}

	if h, ok := handler.(InActiveHandler); ok {
		adapter.inActiveHandler = h
		adapter.flag |= InActive
	}

	if h, ok := handler.(ReadHandler); ok {
		adapter.readHandler = h
		adapter.flag |= Read
	}

	if h, ok := handler.(WriteHandler); ok {
		adapter.writeHandler = h
		adapter.flag |= Write
	}

	if h, ok := handler.(ErrorHandler); ok {
		adapter.errorHandler = h
		adapter.flag |= HandleError
	}

	return adapter
}

func (adapter *ChannelHandler) ChannelActive(ctx IChannelContext) {
	if adapter.activeHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel active", adapter.name)
		adapter.activeHandler.ChannelActive(ctx)
	}
}

func (adapter *ChannelHandler) ChannelInActive(ctx IChannelContext) {
	if adapter.inActiveHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel inactive", adapter.name)
		adapter.inActiveHandler.ChannelInActive(ctx)
	}
}

func (adapter *ChannelHandler) ChannelRead(ctx IChannelContext, msg interface{}) {
	if adapter.readHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel read", adapter.name)
		adapter.readHandler.ChannelRead(ctx, msg)
	}
}

func (adapter *ChannelHandler) Write(ctx IChannelContext, msg interface{}) {
	if adapter.writeHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel write", adapter.name)
		adapter.writeHandler.Write(ctx, msg)
	}
}

func (adapter *ChannelHandler) HandleError(ctx IChannelContext, err error) {
	if adapter.errorHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel handle error", adapter.name)
		adapter.errorHandler.HandleError(ctx, err)
	}
}

func (adapter *ChannelHandler) Flag() Flag {
	return adapter.flag
}
