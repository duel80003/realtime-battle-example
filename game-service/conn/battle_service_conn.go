package conn

import (
	"game-service/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"os"
	"sync"
	"time"
)

var (
	battleServiceConn     *grpc.ClientConn
	battleServiceConnOnce sync.Once
	logger                = tools.Logger
)

func InitBattleGRPCConn() {
	battleServiceConnOnce.Do(func() {
		addr := os.Getenv("BATTLE_SERVICE")
		logger.Infof("battle grpc addr %s", addr)
		conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second * 3,
			PermitWithoutStream: true,
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Panicf("atuh battle service connection error: %s", err)
		}
		battleServiceConn = conn
	})
}

func CloseBattleServiceGrpcConn() {
	err := battleServiceConn.Close()
	if err != nil {
		logger.Errorf("closing auth grpc connection error: %s", err)
		return
	}
	logger.Infof("auth grpc disconnected")
}

func GetBattleServiceGRpcConn() *grpc.ClientConn {
	return battleServiceConn
}
