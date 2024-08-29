package transport

import (
	"fmt"
	"net"

	"github.com/darkphotonKN/starlight-cargo/internal/types"
	"github.com/google/uuid"
)

type TCPTransport struct {
	Conn     net.Conn
	Peers    map[string]types.Peer
	listener net.Listener
	opts     Opts
}

type Opts struct {
	Port int
}

// factory function to  TCPTransport
func NewTCPTransport(opts Opts) types.Transport {
	return &TCPTransport{
		opts:  opts,
		Peers: make(map[string]types.Peer),
	}
}

type ConnState string

const (
	connected    ConnState = "connected"
	disconnected ConnState = "disconnected"
)

// Represents a single node on a network
type Peer struct {
	ID    string
	Addr  net.Addr
	State ConnState
	Conn  net.Conn
}

func (p *Peer) Connect() {
	p.State = connected
}

func (p *Peer) Disconnect() {
	p.State = disconnected
}

func NewPeer(addr net.Addr, conn net.Conn) types.Peer {
	return &Peer{ID: uuid.New().String(), Addr: addr, Conn: conn, State: disconnected}
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
		// save remote address location of the connection
		newPeer := NewPeer(conn.RemoteAddr(), conn)

		t.Peers[newPeer.(Peer).ID] = newPeer
		go t.handleConnection(conn)
	}
}

/**
* Starts individual goroutines to serve incoming messages.
**/
func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer conn.Close()

	MAX_FILE_SIZE := 2048

	for {
		// handle messages with this new peer new connection
		var buf = make([]byte, MAX_FILE_SIZE)

		bufLen, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("Error reading incoming message: %s\n", err)
			return
		}

		msg := buf[:bufLen]
		fmt.Printf("Message received for connection %s", msg)

		for key, val := range t.Peers {
			fmt.Printf("Peer: %+v\n", val)
			fmt.Printf("Key: %+v\n\n", key)
		}
	}
}
