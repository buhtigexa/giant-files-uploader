package main

import (
	"bugtigexa.giantfilesuploader.com/cmd"
	"context"
	"errors"
	"flag"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func waitForShutdown(cancel context.CancelFunc, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⚠️ Interrupt signal received, shutting down...")
	cancel()
	srv.Shutdown(context.Background())

}

func main() {

	port := flag.String("port", "3000", "port to serve on")
	host := flag.String("host", "localhost", "host to serve on")
	flag.Parse()

	streamServer := cmd.NewStreamServer(*host, "3001")

	ctx, cancel := context.WithCancel(context.Background())

	go streamServer.Start(ctx)

	apiServer := &http.Server{
		Addr:    *host + ":" + *port,
		Handler: streamServer.Routes(),
	}

	go waitForShutdown(cancel, apiServer)

	if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

}
