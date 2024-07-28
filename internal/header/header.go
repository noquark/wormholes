package header

import (
	"log"
	"time"
)

// Show a basic notice about program on start.
func Show() {
	log.SetFlags(0)
	log.Println()
	log.Printf("Wormholes")
	log.Printf("Copyright © %d No Quark Labs", time.Now().Year())
	log.Println("Licensed under GNU AFFERO GENERAL PUBLIC LICENSE 3.0")
	log.Println()
}
