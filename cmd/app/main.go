package main

import (
	"fmt"
	"log"

	"github.com/darkphotonKN/starlight-cargo/internal/transport"
)

func main() {
	fmt.Println("---------Starlight Cargo---------")

	opts := transport.Opts{
		Port: 3600,
	}

	tcpTransport := transport.NewTCPTransport(opts)
	err := tcpTransport.ListenAndAccept()

	if err != nil {
		log.Fatalf("Could not initiate tcp listeners. Err: %s\n", err)
	}

	// TODO: Remove Later
	select {}
}
