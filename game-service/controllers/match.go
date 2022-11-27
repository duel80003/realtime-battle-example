package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	grpc "game-service/conn"
	proto "game-service/proto/gen"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

func MatchHandler(request *Request) {
	logger.Infof("MatchHandler: %+v", request)
	wait := doMatching(request)
	matchResult := <-wait
	doBattle(matchResult)
	logger.Infof("battle end.")
}

func doMatching(request *Request) chan *MatchResponse {
	battleMatchClient := proto.NewMatchServiceClient(grpc.GetBattleServiceGRpcConn())
	stream, err := battleMatchClient.Match(context.TODO())
	if err != nil {
		logger.Fatalf("client.RouteChat failed: %v", err)
	}
	wait := make(chan *MatchResponse)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}
			if err != nil {
				logger.Fatalf("client.doMatching failed: %v", err)
			}
			res := &MatchResponse{
				BattleID:  in.BattleId,
				PlayerID1: in.P1Id,
				PlayerID2: in.P2Id,
			}
			logger.Infof("match result: %+v", res)
			body, err := json.Marshal(res)
			if err != nil {
				logger.Errorf("received and marshal error: %s", err)
			}
			conn1 := getConnInfo(in.P1Id)
			conn2 := getConnInfo(in.P2Id)
			if conn1 != nil {
				conn1.WriteMessage(TextMessage, body)
			}
			if conn2 != nil {
				conn2.WriteMessage(TextMessage, body)
			}
			wait <- res
		}
	}()
	stream.Send(&proto.MatchRequest{
		PlayerId: request.PlayerID,
	})
	stream.CloseSend()
	return wait
}

func doBattle(match *MatchResponse) {
	logger.Infof("BattleHandler, match result: %+v", match)
	defer func() {
		deleteConn(match.PlayerID1)
		deleteConn(match.PlayerID2)
	}()
	battleServiceClient := proto.NewBattleServiceClient(grpc.GetBattleServiceGRpcConn())
	stream, err := battleServiceClient.RealtimeBattle(context.TODO())
	if err != nil {
		logger.Fatalf("client.RouteChat failed: %v", err)
	}
	wait := make(chan struct{})
	conn1 := getConnInfo(match.PlayerID1)
	conn2 := getConnInfo(match.PlayerID2)
	go func() {
		battleData := &BattleData{}
		defer CloseWSConn(conn1, conn2)
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}
			if err != nil {
				logger.Fatalf("client.doBattle failed: %v", err)
			}
			for _, v := range in.GetRoleInfos() {
				if v.SkillName == "待機" {
					continue
				}
				if in.IsEnd {
					stream.Send(&proto.RealtimeBattleRequest{
						BattleId: in.BattleId,
						Close:    true,
					})
					result := &BattleResult{
						Winner: in.Winner,
					}
					body, _ := json.Marshal(result)
					writeMessage(conn1, conn2, &body)
					return
				}
				if v.GetIsMissed() {
					logger.Infof("%s 迴避了 %s 的 %s", v.TargetRoleName, v.RoleName, v.SkillName)
					battleData.Text = fmt.Sprintf("%s 迴避了 %s 的 %s", v.TargetRoleName, v.RoleName, v.SkillName)
				} else {
					logger.Infof("%s 對 %s 使用: %s, 造成: %d 傷害", v.RoleName, v.TargetRoleName, v.SkillName, v.Damage)
					battleData.Text = fmt.Sprintf("%s 對 %s 使用: %s, 造成: %d 傷害", v.RoleName, v.TargetRoleName, v.SkillName, v.Damage)
				}
				body, _ := json.Marshal(battleData)
				writeMessage(conn1, conn2, &body)
			}

		}
	}()
	stream.Send(&proto.RealtimeBattleRequest{
		BattleId: match.BattleID,
	})
	<-wait
	stream.CloseSend()
}

func writeMessage(conn1, conn2 *websocket.Conn, body *[]byte) {
	if conn1 != nil {
		conn1.WriteMessage(TextMessage, *body)
	}
	if conn2 != nil {
		conn2.WriteMessage(TextMessage, *body)
	}
}

func CloseWSConn(conn1, conn2 *websocket.Conn) {
	logger.Infof("close connection")
	time.AfterFunc(time.Second*1, func() {
		if conn1 != nil {
			conn1.Close()
		}
		if conn2 != nil {
			conn2.Close()
		}
	})
}
