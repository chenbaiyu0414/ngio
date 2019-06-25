package main

import (
	"math"
	"ngio"
	"ngio/codec"
	"ngio/example/echo"
	"ngio/internal/logger"
	"ngio/transport/tcp"
	"os"
	"os/signal"
)

func main() {
	client, err := tcp.Dial("tcp4", "", "localhost:9863",
		tcp.WithNoDelay(true),
		tcp.WithKeepAlive(true),
		tcp.WithInitializer(func(ch ngio.Channel) {
			ch.Pipeline().Append(codec.NewByteToMessageDecoderAdapter(
				codec.NewLineBasedFrameDecoder(math.MaxUint8, true)))

			ch.Pipeline().Append(new(echo.Handler))
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
