package header

import (
	"fmt"
	"log"
	"time"
)

// Show a basic notice about program on start.
func Show(name string) {
	log.SetFlags(0)
	log.Println()
	log.Println(fmt.Sprintf("Wormholes | %s Service", name))
	log.Printf("Copyright © %d Mohit Singh", time.Now().Year())
	log.Println("Licensed under GNU AFFERO GENERAL PUBLIC LICENSE 3.0")
	log.Println()
}
