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
			ch.Pipeline().Append("encoder", codec.NewByteToMessageDecoderAdapter(
				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))

			ch.Pipeline().Append("handler", echo.NewHandler())
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
