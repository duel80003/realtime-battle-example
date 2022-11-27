package servers

import (
	battleHelper "battle-service/battle_hepler"
	proto "battle-service/proto/gen"
	"fmt"
	"io"

	"time"
)

type BattleServer struct{}

func (b BattleServer) RealtimeBattle(stream proto.BattleService_RealtimeBattleServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		players := getBattleIDAndRoles().Get(in.BattleId)
		logger.Infof("RealtimeBattle, player ids : %v", players)
		if len(players) != 2 {
			return fmt.Errorf("invalid players")
		}
		if in.Close {
			logger.Infof("battle %s end", in.BattleId)
			getBattleIDAndRoles().Delete(in.BattleId)
			return nil
		}
		battleAlgorithm(in.BattleId, players, stream)
	}
}

func battleAlgorithm(battleID string, players []string, stream proto.BattleService_RealtimeBattleServer) {
	defer logger.Infof("battleAlgorithm: %s done", battleID)
	rid1 := fakePlayerData[players[0]].Role
	rid2 := fakePlayerData[players[1]].Role
	role1 := battleHelper.GetRole(rid1)
	role2 := battleHelper.GetRole(rid2)
	roles := []*battleHelper.Role{
		&role1,
		&role2,
	}
	m := make(map[string]string)
	m[rid1] = players[0]
	m[rid2] = players[1]
	for battleHelper.GameContinue(roles) {
		time.Sleep(time.Millisecond * 100)
		res := battleHelper.StartBattle(roles)
		res.BattleId = battleID
		if res.IsEnd {
			if roles[0].IsDead {
				res.Winner = m[roles[1].ID]
			} else {
				res.Winner = m[roles[0].ID]
			}
		}
		stream.Send(res)
	}

}
