package main

import (
	"context"
	"game-service/conn"
	"game-service/controllers"
	"game-service/tools"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load()
	if err != nil {
		tools.Logger.Warnf("env file lost")
	}
	tools.LogInit()
	conn.InitBattleGRPCConn()

}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", controllers.EntryHandler)
	router.HandleFunc("/ping", controllers.PingHandler)
	host := os.Getenv("HOST")
	srv := &http.Server{
		Addr: host,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	go func() {
		tools.Logger.Infof("game server starting")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	conn.CloseBattleServiceGrpcConn()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
