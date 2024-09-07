package test

import (
	"testing"

	"github.com/darkphotonKN/starlight-cargo/internal/transport"
	"github.com/stretchr/testify/assert"
)

func TestTCP(t *testing.T) {
	// set up new tcp connection to simulate a client attemtping to connect
	opts := transport.Opts{
		ListenAddr: 3999,
	}

	tcp := transport.NewTCPTransport(opts).(*transport.TCPTransport)
	// test listen address is correct
	assert.Equal(t, tcp.Opts.ListenAddr, 3999)

	// test that there is no error from ListenAndAccept
	assert.Nil(t, tcp.ListenAndAccept())
}
