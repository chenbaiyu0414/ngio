# NGIO

`NGIO` is a netty like network I/O framework.

**The API is not stable and the performance is not optimized yet. Don't use it in product environment at present.**

# Usage

```go
package main

import (
	"ngio"
	"ngio/transport/tcp"
	"ngio/channel"
	"ngio/codec"
	"ngio/internal/logger"
	"ngio/buffer"
	"os"
	"os/signal"
)

func main() {
	srv := tcp.NewServer(
		tcp.WithNoDelay(true),
		tcp.WithKeepAlive(true),
		tcp.WithInitializer(func(ch channel.Channel) {
			ch.Pipeline().Append(codec.NewByteToMessageDecoderWrapper(codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))
			ch.Pipeline().Append(&handler{})
		}))

	go func() {
		if err := srv.Serve("tcp4", "localhost:9863"); err != nil {
			t.Error(err)
			return
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch

	srv.Shutdown()
}

type handler struct {
}

func (*handler) ChannelRead(ctx channel.Context, msg interface{}) {
	r, ok := msg.(buffer.ByteBuffer)
	if !ok {
		logger.Errorf("msg is not buffer.ByteBuffer")
		return
	}

	text := string(r.ReadBytes(r.ReadableBytes()))

	logger.Infof("received: %s", text)

	b := buffer.NewByteBuf(make([]byte, 1024))
	b.WriteBytes([]byte(text))

	ctx.Write(b)
}

func (*handler) ChannelInActive(ctx channel.Context) {
	logger.Infof("inactive")
}

func (*handler) ChannelActive(ctx channel.Context) {
	logger.Infof("active")
}
```

# 