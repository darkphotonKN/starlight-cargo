package transport

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/darkphotonKN/starlight-cargo/internal/types"
	"github.com/google/uuid"
)

// TODO: temp enacting as the database of users
type User struct {
	email    string
	password string
}

var users = make(map[string]User)

type TCPTransport struct {
	Peers    map[string]*Peer
	listener net.Listener
	Opts     Opts
}

type Opts struct {
	ListenAddr uint
}

// factory function to TCPTransport
func NewTCPTransport(opts Opts) types.Transport {
	// preloading users
	users[uuid.NewString()] = User{
		email:    "darkphoton20@gmail.com",
		password: "123456",
	}

	return &TCPTransport{
		Opts:  opts,
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
	port := fmt.Sprintf(":%d", t.Opts.ListenAddr)

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
	defer func() {
		fmt.Printf("Closing connection for peer %s", conn.RemoteAddr())
		conn.Close()
	}()

	// -- authorize user loop --

	for {
		// email input handling
		conn.Write([]byte("Please enter email.\n"))

		// read response
		var emailResBuf = make([]byte, 128)
		n, err := conn.Read(emailResBuf)

		if err != nil {
			fmt.Printf("Error reading incoming message for email input: %s\n", err)
			continue
		}

		// only extract the readable part of the buffer
		trimmedEmailRes := strings.TrimSpace(string(emailResBuf[:n]))

		fmt.Println("Received email:", trimmedEmailRes)

		// check if user exists
		exists := false
		var existingUser User

		for _, user := range users {
			if user.email == trimmedEmailRes {
				exists = true
				existingUser = user
			}
		}

		if !exists {
			fmt.Println("Email was incorrect.", string(emailResBuf))
			conn.Write([]byte("Email was incorrect.\n"))

			// repeat and ask again
			continue
		}

		// password input handling
		conn.Write([]byte("Please enter password.\n"))

		var pwResBuf = make([]byte, 128)
		n, err = conn.Read(pwResBuf)

		if err != nil {
			fmt.Printf("Error reading incoming message for password input: %s\n", err)
			continue
		}

		// only extract the readable part of the buffer
		trimmedPwRes := strings.TrimSpace(string(pwResBuf[:n]))

		if trimmedPwRes != existingUser.password {
			fmt.Println("Password was incorrect.", string(trimmedPwRes))
			conn.Write([]byte("Password was incorrect.\n"))
			continue
		}

		// break out of loop if user exist
		break
	}

	fmt.Println("Auth passed")

	// -- user authenticated - handle command payload loop --

	MAX_MSG_SIZE := 2048

	for {
		fmt.Println("Starting connected read loop.")

		conn.Write([]byte("You have been connected. Please type a command.\n"))

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
			t.broadcastToAll(payload)

		default:
			fmt.Println("Unknown command.")
			// stop function and disconnect
			return
		}

	}
}

/**
* Parses a command from the tcp peer-to-peer msg.
* Assumes space has the the pre-defined meaning of separating command from payload.
**/
func (t *TCPTransport) parseCommand(msg []byte) (string, []byte) {

	// NOTE: force only one split, with space being the separator of command and payload
	cmdAndPayload := bytes.SplitN(msg, []byte(" "), 2)

	// check for message to fit the predefined protocol and hence have two parts
	if len(cmdAndPayload) < 2 {
		return string(cmdAndPayload[0]), nil

	}

	// first part is the command, second part the payload
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

func (t *TCPTransport) AuthenticateUser() {

}
