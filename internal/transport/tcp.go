package transport

import (
	"bytes"
	"fmt"
	"net"

	"github.com/darkphotonKN/starlight-cargo/internal/types"
	"github.com/google/uuid"
)

type TCPTransport struct {
	Peers    map[string]*Peer
	listener net.Listener
	opts     Opts
}

type Opts struct {
	Port int
}

// factory function to TCPTransport
func NewTCPTransport(opts Opts) types.Transport {
	return &TCPTransport{
		opts:  opts,
		Peers: make(map[string]*Peer),
	}
}

type ConnState string

const (
	CONNECTED    ConnState = "connected"
	DISCONNECTED ConnState = "disconnected"
)

// Represents a single node on a network
type Peer struct {
	ID    string
	Addr  net.Addr
	State ConnState
	Conn  net.Conn
}

func (p *Peer) Connect() {
	p.State = CONNECTED
}

func (p *Peer) Disconnect() {
	p.Conn.Close()
	p.State = DISCONNECTED
}

func NewPeer(addr net.Addr, conn net.Conn) types.Peer {
	return &Peer{ID: uuid.New().String(), Addr: addr, Conn: conn, State: DISCONNECTED}
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

		fmt.Printf("New Connection\n\n\n")

		if err != nil {
			fmt.Printf("Error when attempting to accept connection")
			// continue loop, serve next incoming connection
			continue
		}

		// serve individual connections
		// save remote address location of the connection
		newPeer := NewPeer(conn.RemoteAddr(), conn)

		peer, ok := newPeer.(*Peer)

		if !ok {
			fmt.Printf("Error when attempting to assert peer from interface.")
			return
		}

		t.Peers[peer.ID] = peer
		peer.Connect() // update to connected state
		go t.handleConnection(conn)
	}
}

const (
	CMD_UPLOAD   string = "upload"
	CMD_DOWNLOAD string = "download"
	CMD_MESSAGE  string = "message"
)

/**
* Starts individual goroutines to serve incoming messages.
**/
func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer conn.Close()

	MAX_MSG_SIZE := 2048

	for {
		// handle messages with this new peer new connection
		var buf = make([]byte, MAX_MSG_SIZE)

		bufLen, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("Error reading incoming message: %s\n", err)
			return
		}

		msg := buf[:bufLen]

		fmt.Printf("%s", string(msg))

		// read command from buffer
		command, payload := t.parseCommand(msg)

		switch command {
		case CMD_UPLOAD:
			fmt.Printf("Received upload command. Payload: %s", payload)

		case CMD_DOWNLOAD:
			fmt.Printf("Received download command. Payload: %s", payload)

		case CMD_MESSAGE:
			fmt.Printf("Received message command. Payload: %s", payload)
		default:
			fmt.Printf("No matching command.")
			// stop function and disconnect
			return
		}

		t.broadcastToAll(msg)
	}
} /**
* Parses a command from the tcp peer-to-peer msg.
* Assumes space has the the pre-defined meaning of separating command from payload.
**/
func (t *TCPTransport) parseCommand(msg []byte) (string, []byte) {

	// NOTE: force only one split, with space being the separator of command and payload
	cmdAndPayload := bytes.SplitN(msg, []byte(" "), 2)

	return string(cmdAndPayload[0]), cmdAndPayload[1]
}

/**
* Broadcasts message to all users.
**/
func (t *TCPTransport) broadcastToAll(msg []byte) {
	fmt.Println("Broadcasting to all..")

	// broadcast message to all connected peers
	for _, peer := range t.Peers {
		fmt.Println("Looping peers.. current peer:", peer)
		if peer.State == CONNECTED {
			_, err := peer.Conn.Write([]byte(fmt.Sprintf("peer %s: %s", peer.ID, msg)))

			if err != nil {
				fmt.Println("Error when writing message to peer:", err)
			}
		}
	}

}
