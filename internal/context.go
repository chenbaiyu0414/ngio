package internal

import (
	"bytes"
	"fmt"
	"ngio/logger"
	"strconv"
)

type channelContext struct {
	name       string
	next, prev *channelContext
	pipeline   IChannelPipeline
	handler    IChannelHandler
	log        logger.Logger
}

func newContext(name string, handler interface{}, pipeline IChannelPipeline) *channelContext {
	switch handler.(type) {
	case ActiveHandler, InActiveHandler, ReadHandler, WriteHandler, ErrorHandler, nil: // nil for tail channelContext
	default:
		panic(fmt.Errorf(`invalid handler type. name: "%s"`, name))
	}

	return &channelContext{
		name:     name,
		pipeline: pipeline,
		handler:  NewChannelHandler(name, handler),
		log:      logger.DefaultLogger(),
	}
}

func (ctx *channelContext) FireChannelActiveHandler() {
	next := ctx.findInboundContext(Active)

	if next == nil {
		ctx.log.Warnf("[%v] fire next active handler failed: channelContext after current not contains active handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire active", ctx, next)
	next.handler.ChannelActive(next)
}

func (ctx *channelContext) FireChannelInActiveHandler() {
	next := ctx.findInboundContext(InActive)

	if next == nil {
		ctx.log.Warnf("[%v] fire next inactive handler failed: channelContext after current not contains inactive handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire inactive", ctx, next)
	next.handler.ChannelInActive(next)
}

func (ctx *channelContext) FireChannelReadHandler(msg interface{}) {
	next := ctx.findInboundContext(Read)

	if next == nil {
		ctx.log.Warnf("[%v] fire next read handler failed: channelContext after current not contains read handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire read", ctx, next)
	next.handler.ChannelRead(next, msg)
}

func (ctx *channelContext) Write(msg interface{}) {
	next := ctx.findOutboundContext(Write)

	if next == nil {
		// current channelContext is the lasted outbound channelContext.
		ctx.log.Warnf("[%v] fire next write handler failed: channelContext before current not contains write handler.", ctx)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire write", ctx, next)
	next.handler.Write(next, msg)
}

func (ctx *channelContext) FireChannelErrorHandler(err error) {
	next := ctx.findInboundContext(HandleError)

	if next == nil {
		ctx.log.Warnf("[%v] fire next error handler failed: channelContext after current not contains error handler. error: %v", ctx, err)
		return
	}

	defer interceptError(ctx)

	ctx.log.Debugf("[%v] => [%v] fire handle error", ctx, next)
	next.handler.HandleError(next, err)
}

func (ctx *channelContext) Next() IChannelContext {
	return ctx.next
}

func (ctx *channelContext) Prev() IChannelContext {
	return ctx.prev
}

func (ctx *channelContext) Pipeline() IChannelPipeline {
	return ctx.pipeline
}

func (ctx *channelContext) Name() string {
	return ctx.name
}

func (ctx *channelContext) String() string {
	buf := bytes.Buffer{}

	buf.WriteString(`channelContext: "`)
	buf.WriteString(ctx.name)
	buf.WriteString(`", channel id: `)
	buf.WriteString(strconv.FormatInt(int64(ctx.pipeline.Channel().Id()), 10))

	return buf.String()
}

func (ctx *channelContext) findInboundContext(flag Flag) *channelContext {
	for next := ctx.next; next != nil; next = next.next {
		if flag&flag == flag {
			return next
		}
	}

	return nil
}

func (ctx *channelContext) findOutboundContext(flag Flag) *channelContext {
	for prev := ctx.prev; prev != nil; prev = prev.prev {
		if flag&flag == flag {
			return prev
		}
	}

	return nil
}

func interceptError(ctx *channelContext) {
	if v := recover(); v != nil {
		var e error

		if vErr, ok := v.(error); ok {
			e = vErr
		} else {
			e = fmt.Errorf("%v", v)
		}

		// if current channelContext contains error handler, then invoke current channelContext's error handler.
		// if not, invoke the next channelContext which contains error handler.
		// if there's no channelContext contains error handler after current channelContext, log it.

		// when invoke error handle, the param ctx should be current channelContext, it specified
		// which channelContext that error occurred.
		if ctx.handler.Flag()&HandleError == HandleError {
			ctx.handler.HandleError(ctx, e)
		} else if next := ctx.findInboundContext(HandleError); next != nil {
			ctx.handler.HandleError(ctx, e)
		} else {
			ctx.log.Errorf("unhandled error: %v", e)
		}
	}
}
