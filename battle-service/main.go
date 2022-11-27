package main

import (
	battleHelper "battle-service/battle_hepler"
	"battle-service/servers"
	"battle-service/tools"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net"
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
	battleHelper.ReadData()
}

func serverStart() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		tools.Logger.Panicf("empty grpc port")
	}
	s := servers.GRpcServers()
	tools.Logger.Infof("server start: %s", grpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go serverStart()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)
	// Block until we receive our signal.
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
}
