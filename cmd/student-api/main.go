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
	"github.com/kushvardhan/Students-Api/pkg/http/handlers/student"
	"github.com/kushvardhan/Students-Api/pkg/storage/sqlite"
)

func main(){
	// load config

	cfg := config.MustLoad();

	// db setup
	stor, err := sqlite.New(cfg)
	if err != nil{
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux();

	router.HandleFunc("POST /api/students", student.New(stor))
	router.HandleFunc("GET /api/students/{id}", student.GetById(stor));
	router.HandleFunc("GET /api/students", student.GetList(stor));

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