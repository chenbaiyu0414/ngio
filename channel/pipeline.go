package channel

type Pipeline struct {
	ch   Channel
	head *HeadContext
	tail *TailContext
}

func NewPipeline(ch Channel) *Pipeline {
	pipeline := &Pipeline{
		ch: ch,
	}

	pipeline.head = NewHeadContext(pipeline)
	pipeline.tail = NewTailContext(pipeline)

	pipeline.head.SetNext(pipeline.tail)
	pipeline.tail.SetPrev(pipeline.head)

	return pipeline
}

func (pl *Pipeline) Append(handler interface{}) {
	ctx := NewDefaultContext(pl, handler)

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
