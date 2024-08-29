package types

// Represents the transport interface that
// requires types satisfy in order to be viable as a
// transport
type Transport interface {
	ListenAndAccept() error
}

// Represents each node in a network
type Peer interface {
}
