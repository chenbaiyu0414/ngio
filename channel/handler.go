package channel

type Flag int

const (
	None Flag = 1 << iota
	Active
	InActive
	Read
	Write
	Recover
)

type Handler interface {
	ActiveHandler
	InActiveHandler
	ReadHandler
	WriteHandler
	RecoverHandler
	Flag() Flag
}

type ActiveHandler interface {
	ChannelActive(ctx Context)
}

type InActiveHandler interface {
	ChannelInActive(ctx Context)
}

type ReadHandler interface {
	ChannelRead(ctx Context, msg interface{})
}

type WriteHandler interface {
	Write(ctx Context, msg interface{})
}

type RecoverHandler interface {
	ChannelRecovered(ctx Context, v interface{})
}

type HandlerWrapper struct {
	activeHandler   ActiveHandler
	inActiveHandler InActiveHandler
	readHandler     ReadHandler
	writeHandler    WriteHandler
	recoverHandler  RecoverHandler
	flag            Flag
}

func NewHandlerWrapper(handler interface{}) *HandlerWrapper {
	wrapper := new(HandlerWrapper)

	if h, ok := handler.(ActiveHandler); ok {
		wrapper.activeHandler = h
		wrapper.flag |= Active
	}

	if h, ok := handler.(InActiveHandler); ok {
		wrapper.inActiveHandler = h
		wrapper.flag |= InActive
	}

	if h, ok := handler.(ReadHandler); ok {
		wrapper.readHandler = h
		wrapper.flag |= Read
	}

	if h, ok := handler.(WriteHandler); ok {
		wrapper.writeHandler = h
		wrapper.flag |= Write
	}

	if h, ok := handler.(RecoverHandler); ok {
		wrapper.recoverHandler = h
		wrapper.flag |= Recover
	}

	return wrapper
}

func (w *HandlerWrapper) ChannelActive(ctx Context) {
	if w.activeHandler != nil {
		w.activeHandler.ChannelActive(ctx)
	}
}

func (w *HandlerWrapper) ChannelInActive(ctx Context) {
	if w.inActiveHandler != nil {
		w.inActiveHandler.ChannelInActive(ctx)
	}
}

func (w *HandlerWrapper) ChannelRead(ctx Context, msg interface{}) {
	if w.readHandler != nil {
		w.readHandler.ChannelRead(ctx, msg)
	}
}

func (w *HandlerWrapper) Write(ctx Context, msg interface{}) {
	if w.writeHandler != nil {
		w.writeHandler.Write(ctx, msg)
	}
}

func (w *HandlerWrapper) ChannelRecovered(ctx Context, v interface{}) {
	if w.recoverHandler != nil {
		w.recoverHandler.ChannelRecovered(ctx, v)
	}
}

func (w *HandlerWrapper) Flag() Flag {
	return w.flag
}
