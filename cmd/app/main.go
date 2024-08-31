package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darkphotonKN/starlight-cargo/internal/transport"
)

func main() {
	fmt.Println(`  
 _______ _________ _______  _______  _       _________ _______          _________   _______  _______  _______  _______  _______ 
(  ____ \\__   __/(  ___  )(  ____ )( \      \__   __/(  ____ \|\     /|\__   __/  (  ____ \(  ___  )(  ____ )(  ____ \(  ___  )
| (    \/   ) (   | (   ) || (    )|| (         ) (   | (    \/| )   ( |   ) (     | (    \/| (   ) || (    )|| (    \/| (   ) |
| (_____    | |   | (___) || (____)|| |         | |   | |      | (___) |   | |     | |      | (___) || (____)|| |      | |   | |
(_____  )   | |   |  ___  ||     __)| |         | |   | | ____ |  ___  |   | |     | |      |  ___  ||     __)| | ____ | |   | |
      ) |   | |   | (   ) || (\ (   | |         | |   | | \_  )| (   ) |   | |     | |      | (   ) || (\ (   | | \_  )| |   | |
/\____) |   | |   | )   ( || ) \ \__| (____/\___) (___| (___) || )   ( |   | |     | (____/\| )   ( || ) \ \__| (___) || (___) |
\_______)   )_(   |/     \||/   \__/(_______/\_______/(_______)|/     \|   )_(     (_______/|/     \||/   \__/(_______)(_______)
`)

	time.Sleep(time.Millisecond * 1400)
	fmt.Printf("Welcome to Starlight Cargo - Your Galactic File Management System!\n\n")
	time.Sleep(time.Millisecond * 1400)
	fmt.Printf("Initializing the interstellar transport layer...\n\n")

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
