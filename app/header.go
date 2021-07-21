package app

import (
	"log"
	"time"
)

func ShowNotice() {
	log.SetFlags(0)
	log.Println("Wormholes - A self hosted link shortener")
	log.Printf("Copyright © %d Mohit Singh", time.Now().Year())
	log.Println("Licensed under GNU AFFERO GENERAL PUBLIC LICENSE 3.0")
	log.Println()
}
