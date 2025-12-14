package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kushvardhan/Students-Api/pkg/config"
)

func main(){
	// load config

	cfg := config.MustLoad();

	// setup router
	router := http.NewServeMux();

	router.HandleFunc("GET /",func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"));

	})

	// server
	server := http.Server{
		Handler: router,
		Addr: cfg.Address,
	}

	slog.Info("server started", slog.String("address", cfg.Address));

	done := make(chan os.Signal,1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe();
		if err != nil{
			log.Fatal("Failed to start server");
		}
}()

<- done

slog.Info("Shutting down the server");

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

defer cancel()

if err := server.Shutdown(ctx); err != nil{
	slog.Error("Failed to shutdown server", slog.String("error", err.Error()));
}

slog.Info("sever shutdown successfully");

}