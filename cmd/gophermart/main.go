package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"

	"ivanmyagkov/gofermart/internal/config"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/server"
	"ivanmyagkov/gofermart/internal/storage"
)

func init() {
	err := env.Parse(&config.EnvVar)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&config.Flags.A, "a", config.EnvVar.RunAddress, "server address")
	flag.StringVar(&config.Flags.D, "d", config.EnvVar.DatabaseURI, "database uri")
	flag.StringVar(&config.Flags.R, "f", config.EnvVar.AccrualSystemAddress, "accrual system address")
	flag.Parse()
}
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	var err error
	//config
	cfg := config.NewConfig(config.Flags.A, config.Flags.D, config.Flags.R)

	//db
	var db interfaces.DB
	db, err = storage.NewDB(cfg.GetDatabaseURI(), ctx)
	if err != nil {
		log.Fatalf("Failed to create db %e", err)
	}
	srv := server.InitSrv(db)

	go func() {

		<-signalChan

		log.Println("Shutting down...")

		cancel()
		if err := srv.Shutdown(ctx); err != nil && err != ctx.Err() {
			srv.Logger.Fatal(err)
		}

		if err = db.Close(); err != nil {
			log.Println("Failed db...", err)
		}
	}()

	serverAddress := cfg.GetRunAddress()
	if err := srv.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		srv.Logger.Fatal(err)
	}
}
