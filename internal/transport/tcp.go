package transport

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/darkphotonKN/starlight-cargo/internal/auth"
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

		// start a go routine to concurrently check for idle connections
		go t.observeConnections()
	}
}

// constants for commands
const (
	CMD_UPLOAD   string = "upload"
	CMD_DOWNLOAD string = "download"
	CMD_MESSAGE  string = "message"
)

// constants for status
const (
	AUTHENTICATED string = "AUTHENTICATED"
)

/**
* Starts individual goroutines to serve incoming messages.
**/
func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer func() {
		fmt.Printf("Closing connection for peer %s", conn.RemoteAddr())
		conn.Close()
	}()

	// -- 1. authorize user loop --
	err := t.handleAuthorizationLoop(conn)

	if err != nil {
		fmt.Println(err)
		return
	}

	// -- 2. user authenticated - handle command payload loop --
	err = t.handleCommandLoop(conn)

	if err != nil {
		fmt.Println(err)
		return
	}
}

/**
* Loop that authorizes the user until they've entered a mistake three times.
**/
func (t *TCPTransport) handleAuthorizationLoop(conn net.Conn) error {
	attempts := 0

	for {
		// exit out of loop and hence close connection (defer of handleConnection) if user enters
		// the wrong email or password 3 times

		if attempts == 3 {
			t.sendErrorMessage(ErrorMsg{errorType: MSG_ERROR, customMsg: "Too many attempts."}, conn)
			return fmt.Errorf("Too many attempts.")
		}

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
			attempts++
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
			attempts++
			conn.Write([]byte("Password was incorrect.\n"))
			continue
		}

		// authentication passed, provide credentials then break out to next step
		accessToken, err := auth.GenerateJWT(existingUser.email, []byte(auth.SECRET_KEY))
		if err != nil {
			fmt.Println("Error when generating jwt.")
			continue
		}
		conn.Write([]byte(fmt.Sprintf("%s:%s\n", AUTHENTICATED, accessToken)))

		break
	}

	// no error, successfully passed authentication
	return nil
}

/**
* Authorized command loop to handle command operations.
**/
func (t *TCPTransport) handleCommandLoop(conn net.Conn) error {

	MAX_MSG_SIZE := 2048
	for {
		conn.Write([]byte("You are connected. Please type a command.\n"))

		// handle messages with this new peer new connection
		var buf = make([]byte, MAX_MSG_SIZE)

		bufLen, err := conn.Read(buf)

		// clean up connection if client has disconnected
		if err == net.ErrClosed {
			return err
		}

		if err != nil {
			fmt.Printf("Error reading incoming message: %s\n", err)
			// let user try again if message read errored
			continue
		}

		msg := buf[:bufLen]

		fmt.Printf("%s", string(msg))

		// read command from buffer
		accessToken, command, payload := t.parseCommand(msg)

		fmt.Println()
		fmt.Println("DEBUG accessToken:", accessToken)
		fmt.Println()

		// authorize access token
		_, err = auth.ValidateJWT(accessToken, []byte(auth.SECRET_KEY))

		// authorization failed break loop and disconnect user
		if err != nil {
			t.sendErrorMessage(ErrorMsg{errorType: AUTH_ERROR}, conn)

			return fmt.Errorf("Error when validating access token.")
		}

		switch command {
		case CMD_UPLOAD:
			fmt.Printf("Received upload command. Payload: %s", payload)

		case CMD_DOWNLOAD:
			fmt.Printf("Received download command. Payload: %s", payload)

		case CMD_MESSAGE:
			fmt.Printf("Received message command. Payload: %s", payload)
			conn.Write([]byte(fmt.Sprintf("Received payload: %s\n", payload)))

		default:
			fmt.Println("Unknown command.")
			// stop function and disconnect
			continue
		}
	}
}

type ErrorType int

const (
	AUTH_ERROR ErrorType = iota
	MSG_ERROR
)

type ErrorMsg struct {
	errorType ErrorType
	customMsg string
}

/**
* Sends a fixed messages back to client.
**/
func (t *TCPTransport) sendErrorMessage(e ErrorMsg, conn net.Conn) {
	msg := ""

	switch e.errorType {

	case AUTH_ERROR:
		if e.customMsg != "" {
			msg = e.customMsg
		} else {
			msg = "Access token was unauthorized."
		}

	case MSG_ERROR:
		if e.customMsg != "" {
			msg = e.customMsg
		} else {
			msg = "Access token was unauthorized."
		}
	}

	conn.Write([]byte(msg))
}

/**
* Parses a command from the tcp peer-to-peer msg.
* Assumes space has the the pre-defined meaning of separating command from payload.
**/
func (t *TCPTransport) parseCommand(msg []byte) (string, string, []byte) {

	// NOTE: force split into pre-defined structure of [accessToken] space [command] space [payload]
	msgPack := bytes.SplitN(msg, []byte(" "), 3)

	// check for message to fit the predefined protocol and hence have two parts
	if len(msgPack) < 3 {
		return string(msgPack[0]), "", nil

	}

	return string(msgPack[0]), string(msgPack[1]), msgPack[2]
}

/**
* Helper that broadcasts a message to all connected users.
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

/**
* TODO: Add actual check.
* Observe connections for idle connections and handle them.
**/
func (t *TCPTransport) observeConnections() {
	// for {
	// 	// check idle connections every ten seconds
	// 	time.Sleep(time.Second * 10)
	//
	// 	for index, peer := range t.Peers {
	// 		fmt.Printf("Current Peer no. %s: %+v", index, peer)
	// 	}
	// }
}
