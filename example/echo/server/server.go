package main

import (
	"math"
	"ngio"
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
		Channel(func(ch ngio.Channel) {
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
