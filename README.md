# NGIO

`NGIO` is a netty like network I/O framework.

**The API is not stable and unoptimized yet. Don't use it in product environment.**

# Features
#### Socket
- [x] TCP
- [ ] UDP
- [ ] WebSocket

#### Codec
- [x] DelimiterBasedFrameDecoder
- [x] LengthFieldBasedFrameDecoder
- [x] LengthFieldPrepender
- [x] LineBasedFrameDecoder
- [ ] SSL/TLS
- [ ] HTTP
- [ ] Protobuf
- [ ] ......

# How to Use
### Server

```go
func main() {
	srv := tcp.NewServer(
		tcp.WithNoDelay(true),
		tcp.WithKeepAlive(true),
		tcp.WithInitializer(func(ch ngio.Channel) {
			ch.Pipeline().Append(codec.NewByteToMessageDecoderAdapter(
				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))

			ch.Pipeline().Append(new(EchoHandler))
		}))

	go func() {
		if err := srv.Serve("tcp4", "localhost:9863"); err != nil {
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
func main() {
	client, err := tcp.Dial("tcp4", "", "localhost:9863",
		tcp.WithNoDelay(true),
		tcp.WithKeepAlive(true),
		tcp.WithInitializer(func(ch ngio.Channel) {
			ch.Pipeline().Append(codec.NewByteToMessageDecoderAdapter(
				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))

			ch.Pipeline().Append(new(EchoHandler))
		}))

	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch

	if err := client.Close(); err != nil {
		logger.Errorf("close client: %v", err)
	}
}
```

### Handler

```go
type EchoHandler struct {}

func (*EchoHandler) ChannelRead(ctx ngio.Context, msg interface{}) {
	bf, ok := msg.(buffer.ByteBuffer)
	if !ok {
		logger.Errorf("msg is not buffer.ByteBuffer")
		return
	}

	received := string(bf.GetBytes(bf.ReaderIndex(), bf.ReadableBytes()))

	logger.Infof("received: %s", received)

	ctx.Write(bf)
}

func (*EchoHandler) ChannelInActive(ctx ngio.Context) {
	logger.Infof("inactive")
}

func (*EchoHandler) ChannelActive(ctx ngio.Context) {
	logger.Infof("active")
}
```