package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/caarlos0/env/v6"
	"golang.org/x/sync/errgroup"

	"ivanmyagkov/gofermart/internal/client"
	"ivanmyagkov/gofermart/internal/config"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/server"
	"ivanmyagkov/gofermart/internal/storage"
	"ivanmyagkov/gofermart/internal/workerpool"
)

func main() {
	err := env.Parse(&config.EnvVar)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&config.Flags.A, "a", config.EnvVar.RunAddress, "server address")
	flag.StringVar(&config.Flags.D, "d", config.EnvVar.DatabaseURI, "database uri")
	flag.StringVar(&config.Flags.R, "r", config.EnvVar.AccrualSystemAddress, "accrual system address")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	//config
	cfg := config.NewConfig(config.Flags.A, config.Flags.D, config.Flags.R)

	//db
	var db interfaces.DB
	db, err = storage.NewDB(cfg.GetDatabaseURI(), ctx)
	if err != nil {
		log.Fatalf("Failed to create db %e", err)
	}
	qu := make(chan string, 100)
	orders, err := db.SelectNewOrders()
	if err != nil {
		log.Println(err)
	}
	if len(orders) > 0 {
		for _, order := range orders {
			qu <- order
		}
	}

	g, _ := errgroup.WithContext(ctx)

	client := client.NewAccrualClient(config.Flags.R, db, qu)

	srv := server.InitSrv(db, qu)
	for i := 1; i <= runtime.NumCPU(); i++ {
		worker := workerpool.NewWorker(qu, client, ctx)
		g.Go(worker.Do)
	}

	go func() {

		<-signalChan

		log.Println("Shutting down...")

		if err := srv.Shutdown(ctx); err != nil && err != ctx.Err() {
			srv.Logger.Fatal(err)
		}

		if err = db.Close(); err != nil {
			log.Println("Failed db...", err)
		}
		cancel()
		close(qu)
		err = g.Wait()
		if err != nil {
			log.Println("err-group...", err)
		}

	}()

	serverAddress := cfg.GetRunAddress()
	if err := srv.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		srv.Logger.Fatal(err)
	}
}
