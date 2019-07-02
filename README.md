# NGIO

`NGIO` is a netty like network I/O framework.

**The API is not stable and unoptimized yet. Don't use it in product environment.**

# Features
#### Socket
- [x] TCP
- [x] UDP
- [ ] WebSocket

#### Codec
- [x] DelimiterBasedFrameDecoder
- [x] LengthFieldBasedFrameDecoder
- [x] LengthFieldPrepender
- [x] LineBasedFrameDecoder
- [x] SSL/TLS
- [ ] ~~HTTP~~
- [ ] Protobuf
- [ ] ......

# How to Use
### Server

```go
package main

import (
	"math"
	"ngio"
	"ngio/channel"
	"ngio/codec"
	"ngio/example/echo"
	"ngio/option"
	"os"
	"os/signal"
)

func main() {
	srv := ngio.NewServer("tcp4", "localhost:9863").
    		Option(option.TCPNoDelay(true)).
    		Option(option.TCPKeepAlive(true)).
    		Channel(func(ch channel.Channel) {
    			ch.Pipeline().AddLast("decoder", codec.NewByteToMessageDecoderAdapter(
    				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))
    
    			ch.Pipeline().AddLast("handler", echo.NewHandler())
    		})
    
    	go func() {
    		if err := srv.Serve(); err != nil {
    			panic(err)
    			return
    		}
    	}()
    
    	ch := make(chan os.Signal)
    	signal.Notify(ch, os.Kill, os.Interrupt)
    	<-ch
    
    	srv.Shutdown()
}
```

### Client

```go
package main

import (
	"math"
	"ngio"
	"ngio/channel"
	"ngio/codec"
	"ngio/example/echo"
	"ngio/option"
	"os"
	"os/signal"
)

func main() {
	clt := ngio.NewClient("tcp4", "", "localhost:9863").
    		Option(option.TCPNoDelay(true)).
    		Option(option.TCPKeepAlive(true)).
    		Channel(func(ch channel.Channel) {
    			ch.Pipeline().AddLast("encoder", codec.NewByteToMessageDecoderAdapter(
    				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))
    
    			ch.Pipeline().AddLast("handler", echo.NewHandler())
    		})
    
    	go func() {
    		if err := clt.Dial(); err != nil {
    			panic(err)
    			return
    		}
    	}()
    
    	ch := make(chan os.Signal)
    	signal.Notify(ch, os.Kill, os.Interrupt)
    	<-ch
    
    	clt.Close()
}

```

### Handler

```go
package echo

import (
	"ngio/buffer"
	"ngio/channel"
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

func (handler *Handler) ChannelRead(ctx *channel.Context, msg interface{}) {
	bf, ok := msg.(buffer.ByteBuffer)
	if !ok {
		handler.log.Errorf("msg is not buffer.ByteBuffer")
		return
	}

	received := bf.GetBytes(bf.ReaderIndex(), bf.ReadableBytes())

	handler.log.Infof("received: %s", string(received))

	ctx.Write(bf)
}

func (handler *Handler) ChannelInActive(ctx *channel.Context) {
	handler.log.Infof("inactive")
}

func (handler *Handler) ChannelActive(ctx *channel.Context) {
	handler.log.Infof("active")
}

func (handler *Handler) HandleError(ctx *channel.Context, err error) {
	handler.log.Errorf("unexpected unhandled error: %v", err)
}


```