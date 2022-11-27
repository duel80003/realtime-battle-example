package servers

import (
	proto "battle-service/proto/gen"
	"battle-service/tools"
	"google.golang.org/grpc"
)

var (
	logger = tools.Logger
)

func GRpcServers() *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterBattleServiceServer(s, &BattleServer{})
	proto.RegisterMatchServiceServer(s, MatchServer{})
	return s
}
