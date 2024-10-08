package main

import (
	"fmt"
	"log"
	"time"

	fileservice "github.com/darkphotonKN/starlight-cargo/internal/file_service"
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

	time.Sleep(time.Millisecond * 800)
	fmt.Printf("You have launched the server of the starlight cargo system.\n\n")
	time.Sleep(time.Millisecond * 800)
	fmt.Printf("Initializing the interstellar transport server and awaiting starship connections...\n\n")

	opts := transport.Opts{
		ListenAddr: 3600,
	}

	// creating file service for injection
	fs := fileservice.NewFileService()
	tcpTransport := transport.NewTCPTransport(opts, fs)
	err := tcpTransport.ListenAndAccept()

	if err != nil {
		log.Fatalf("Could not initiate tcp listeners. Err: %s\n", err)
	}

	// TODO: Remove Later
	select {}
}
