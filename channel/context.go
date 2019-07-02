package channel

import (
	"bytes"
	"fmt"
	"ngio/logger"
	"strconv"
)

type Context struct {
	name           string
	next, prev     *Context
	pipeline       *Pipeline
	handlerAdapter *HandlerAdapter
	log            logger.Logger
}

func NewContext(name string, handler interface{}, pipeline *Pipeline) *Context {
	switch handler.(type) {
	case ActiveHandler, InActiveHandler, ReadHandler, WriteHandler, ErrorHandler, *tailHandler:
	default:
		panic(fmt.Errorf(`invalid handler type. name: "%s"`, name))
	}

	return &Context{
		name:           name,
		pipeline:       pipeline,
		handlerAdapter: NewHandlerAdapter(name, handler),
		log:            logger.DefaultLogger(),
	}
}

func (ctx *Context) FireActiveHandler() {
	next := ctx.findInboundContext(Active)

	if next == nil {
		ctx.log.Warnf("[%v] fire next active handler failed: context after current not contains active handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire active", ctx, next)
	next.handlerAdapter.ChannelActive(next)
}

func (ctx *Context) FireInActiveHandler() {
	next := ctx.findInboundContext(InActive)

	if next == nil {
		ctx.log.Warnf("[%v] fire next inactive handler failed: context after current not contains inactive handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire inactive", ctx, next)
	next.handlerAdapter.ChannelInActive(next)
}

func (ctx *Context) FireReadHandler(msg interface{}) {
	next := ctx.findInboundContext(Read)

	if next == nil {
		ctx.log.Warnf("[%v] fire next read handler failed: context after current not contains read handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire read", ctx, next)
	next.handlerAdapter.ChannelRead(next, msg)
}

func (ctx *Context) Write(msg interface{}) {
	next := ctx.findOutboundContext(Write)

	if next == nil {
		// current context is the lasted outbound context.
		ctx.handlerAdapter.Write(ctx, msg)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire write", ctx, next)
	next.handlerAdapter.Write(next, msg)
}

func (ctx *Context) FireChannelErrorHandler(err error) {
	next := ctx.findInboundContext(HandleError)

	if next == nil {
		ctx.log.Warnf("[%v] fire next error handler failed: context after current not contains error handler. error: %v", ctx, err)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire handle error", ctx, next)
	next.handlerAdapter.HandleError(next, err)
}

func (ctx *Context) Next() *Context {
	return ctx.next
}

func (ctx *Context) Prev() *Context {
	return ctx.prev
}

func (ctx *Context) Pipeline() *Pipeline {
	return ctx.pipeline
}

func (ctx *Context) Name() string {
	return ctx.name
}

func (ctx *Context) String() string {
	buf := bytes.Buffer{}

	buf.WriteString(`context: "`)
	buf.WriteString(ctx.name)
	buf.WriteString(`", channel id: `)
	buf.WriteString(strconv.FormatInt(int64(ctx.pipeline.ch.Id()), 10))

	return buf.String()
}

func (ctx *Context) findInboundContext(flag Flag) *Context {
	for next := ctx.Next(); next != nil; next = next.next {
		if next.handlerAdapter.flag&flag == flag {
			return next
		}
	}

	return nil
}

func (ctx *Context) findOutboundContext(flag Flag) *Context {
	for prev := ctx.prev; prev != nil; prev = prev.prev {
		if prev.handlerAdapter.flag&flag == flag {
			return prev
		}
	}

	return nil
}

func interceptError(ctx *Context) {
	if v := recover(); v != nil {
		var e error

		if vErr, ok := v.(error); ok {
			e = vErr
		} else {
			e = fmt.Errorf("%v", v)
		}

		// if current context contains error handler, then invoke current context's error handler.
		// if not, invoke the next context which contains error handler.
		// if there's no context contains error handler after current context, log it.

		// when invoke error handle, the param ctx should be current context, it specified
		// which context that error occurred.
		if ctx.handlerAdapter.flag&HandleError == HandleError {
			ctx.handlerAdapter.HandleError(ctx, e)
		} else if next := ctx.findInboundContext(HandleError); next != nil {
			next.handlerAdapter.HandleError(ctx, e)
		} else {
			ctx.log.Errorf("unhandled error: %v", e)
		}
	}
}
