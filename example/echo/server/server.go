package main

import (
	"math"
	"ngio"
	"ngio/codec"
	"ngio/example/echo"
	"ngio/transport/tcp"
	"os"
	"os/signal"
)

func main() {
	srv := tcp.NewServer(
		tcp.WithNoDelay(true),
		tcp.WithKeepAlive(true),
		tcp.WithInitializer(func(ch ngio.Channel) {
			ch.Pipeline().Append(codec.NewByteToMessageDecoderAdapter(
				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))

			ch.Pipeline().Append(new(echo.Handler))
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
