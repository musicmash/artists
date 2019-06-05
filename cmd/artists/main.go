package main

import (
	"flag"

	"github.com/musicmash/artists/internal/api"
	"github.com/musicmash/artists/internal/config"
	"github.com/musicmash/artists/internal/db"
	"github.com/musicmash/artists/internal/log"
)

func init() {
	configPath := flag.String("config", "/etc/musicmash/artists.yaml", "Path to artists.yaml config")
	flag.Parse()

	if err := config.InitConfig(*configPath); err != nil {
		panic(err)
	}
	if config.Config.Log.Level == "" {
		config.Config.Log.Level = "info"
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)

	db.DbMgr = db.NewMainDatabaseMgr()
}

func main() {
	log.Info("Starting artists service...")
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
