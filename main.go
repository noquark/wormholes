package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/mohitsinghs/wormholes/app"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/links"
)

var port int
var cfgFile, ldb string

func main() {
	flag.IntVar(&port, "port", 3000, "Port to run")
	flag.StringVar(&ldb, "backend", "postgres", "Backend to store links")
	flag.StringVar(&cfgFile, "config", "", "Path to non-default config")
	conf, err := config.LoadDefault()
	if err != nil {
		log.Printf("Failed to read config : %v", err)
	}
	config.Merge("WH", conf)
	flag.Parse()

	pgconn := conf.Postgres.Connect()
	rdbconn := conf.Redis.Connect()

	var linkStore links.Store

	switch ldb {
	case "postgres":
		linkStore = links.NewPgStore(pgconn)
	case "redis":
		linkStore = links.NewRdbStore(rdbconn)
	default:
		log.Fatalln("Backend should be one of postgres and redis")
	}

	f := factory.New()
	f.TryRestore(linkStore.Ids)

	app := Setup(linkStore, f)

	go func() {
		ShowHeader(port)
		app.Listen(fmt.Sprintf(":%d", port))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	_ = <-ch
	if err := app.Shutdown(); err != nil {
		log.Printf("Error stopping server : %v", err.Error())
	}

	if err := f.Backup(); err != nil {
		log.Printf("Error during backup : %v", err.Error())
	}
}
