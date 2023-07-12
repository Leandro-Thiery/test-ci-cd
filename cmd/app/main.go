package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-ci-cd/cmd/config"
	"time"

	"github.com/gin-gonic/gin"
)

func setup(cfg *config.Config) {
	g := gin.Default()
	server := &http.Server{
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              fmt.Sprintf(":%d", 8080),
		Handler:           g,
	}

	go func() {
		log.Printf("listening at port %d", cfg.GetConfig().GetServerConfig().Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Print("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 seconds.
	<-ctx.Done()
	log.Println("timeout of 2 seconds.")
	log.Println("Server exiting")
}

func main() {
	cfg := config.InitConfig()
	setup(cfg)
}
