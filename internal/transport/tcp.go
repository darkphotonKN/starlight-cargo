package transport

import (
	"fmt"
	"net"

	"github.com/darkphotonKN/starlight-cargo/internal/types"
)

type TCPTransport struct {
	Conn     net.Conn
	Peers    map[string]Peer
	listener net.Listener
	opts     Opts
}

type Opts struct {
	Port int
}

// factory function to start TCPTransports
func NewTCPTransport(opts Opts) types.Transport {
	return &TCPTransport{
		opts: opts,
	}
}

// Represents a single node on a network
type Peer struct {
}

/**
* Starts a TCP server and listens for connections.
**/
func (t *TCPTransport) ListenAndAccept() error {
	port := fmt.Sprintf(":%d", t.opts.Port)

	var err error
	t.listener, err = net.Listen("tcp", port)

	if err != nil {
		return err
	}

	// start goroutine to concurrently manage incoming connections
	go t.AcceptLoop()

	// end of function, no error so return nil
	return nil
}

/**
* Reads incoming connections and spins up goroutines for each for message handling.
**/
func (t *TCPTransport) AcceptLoop() {

	for {

		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("Error when attempting to accept connection")
			// continue loop, serve next incoming connection
			continue
		}

		// serve individual connections
		go t.handleConnection(conn)

	}
}

/**
* Starts individual goroutines to serve incoming messages.
**/
func (t *TCPTransport) handleConnection(conn net.Conn) {
	MAX_FILE_SIZE := 2048

	for {
		var buf = make([]byte, MAX_FILE_SIZE)

		bufLen, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("Error reading incoming message: %s\n", err)
		}

		msg := buf[:bufLen]
		fmt.Printf("Message received for connection %s", msg)
	}

}
