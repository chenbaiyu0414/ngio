package channel

import (
	"fmt"
	"sync"
)

type Pipeline struct {
	ch           Channel
	head         *HeadContext
	tail         *TailContext
	handlerNames map[string]struct{}
	mu           sync.Mutex
}

func NewPipeline(ch Channel) *Pipeline {
	pipeline := &Pipeline{
		ch:           ch,
		handlerNames: make(map[string]struct{}, 8),
	}

	pipeline.head = NewHeadContext(pipeline, "head")
	pipeline.tail = NewTailContext(pipeline, "tail")

	pipeline.head.SetNext(pipeline.tail)
	pipeline.tail.SetPrev(pipeline.head)

	return pipeline
}

func (pl *Pipeline) Append(name string, handler interface{}) {
	pl.mu.Lock()
	if _, ok := pl.handlerNames[name]; !ok {
		pl.handlerNames[name] = struct{}{}
	} else {
		panic(fmt.Errorf(`[name: %s] repeat handler name. "head" and "tail" are retained by default`, name))
	}
	pl.mu.Unlock()

	ctx := NewDefaultContext(pl, name, handler)

	prev := pl.tail.Prev()

	prev.SetNext(ctx)
	ctx.SetPrev(prev)
	ctx.SetNext(pl.tail)
	pl.tail.SetPrev(ctx)
}

func (pl *Pipeline) Channel() Channel {
	return pl.ch
}

func (pl *Pipeline) FireActiveHandler() {
	pl.head.FireActiveHandler()
}

func (pl *Pipeline) FireInActiveHandler() {
	pl.head.FireInActiveHandler()
}

func (pl *Pipeline) FireReadHandler(msg interface{}) {
	pl.head.FireReadHandler(msg)
}

func (pl *Pipeline) FireWriteHandler(msg interface{}) {
	pl.tail.Write(msg)
}

func (pl *Pipeline) FireRecoverHandler(v interface{}) {
	pl.head.FireRecoverHandler(v)
}
