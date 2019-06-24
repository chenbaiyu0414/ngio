package transport

type Server interface {
	Serve(network, localAddr string) error
	Shutdown()
}
