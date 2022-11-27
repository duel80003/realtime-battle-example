package servers

import (
	proto "battle-service/proto/gen"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	matchingMap     *MatchingMap
	once            sync.Once
	battleIDMap     *BattleIDAndRoles
	battleIDMapOnce sync.Once
	fakePlayerData  = map[string]*Player{
		"player_1": {ID: "player_1", Score: 1000, Role: "1001"},
		"player_2": {ID: "player_2", Score: 1000, Role: "1002"},
	}
)

type MatchServer struct{}

func (m MatchServer) Match(stream proto.MatchService_MatchServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		res := &proto.MatchResponse{}

		ticker := time.NewTicker(1 * time.Second)
		getMatchingMap().Add(in.GetPlayerId())
		for {
			select {
			case <-ticker.C:
				res.BattleId, res.P1Id, res.P2Id = getMatchingMap().FindMatch()
				if res.BattleId != "" {
					logger.Debugf("match done: %+v", res)
					err = stream.Send(res)
					return err
				}
			}
		}
	}
}

type Player struct {
	ID    string
	Score int32
	Role  string
}

type MatchingMap struct {
	mux     sync.Mutex
	players []*Player
}

func getMatchingMap() *MatchingMap {
	once.Do(func() {
		matchingMap = &MatchingMap{}
	})
	return matchingMap
}

func (m *MatchingMap) Add(playerID string) {
	logger.Infof("MatchingMap add new player id: %s", playerID)
	m.players = append(m.players, fakePlayerData[playerID])
}

func (m *MatchingMap) FindMatch() (battleId, p1, p2 string) {
	if len(m.players) < 2 {
		return
	}
	m.mux.Lock()
	m.mux.Unlock()
	for i := 0; i < len(m.players); i++ {
		for j := i + 1; j < len(m.players); j++ {
			scoreGap := m.players[i].Score - m.players[j].Score
			logger.Debugf("scoreGap: %d", scoreGap)
			if scoreGap >= -15 && scoreGap <= 15 {
				battleId = fmt.Sprintf("%d", time.Now().Unix())
				p1 = m.players[i].ID
				p2 = m.players[j].ID
				getBattleIDAndRoles().Add(battleId, p1, p2)
				m.players = nil
				return
			}
		}
	}
	return
}

type BattleIDAndRoles struct {
	mux sync.Mutex
	m   map[string][]string
}

func getBattleIDAndRoles() *BattleIDAndRoles {
	battleIDMapOnce.Do(func() {
		battleIDMap = &BattleIDAndRoles{
			m: make(map[string][]string),
		}
	})
	return battleIDMap
}

func (b *BattleIDAndRoles) Add(battleId, p1, p2 string) {
	b.mux.Lock()
	b.mux.Unlock()
	b.m[battleId] = []string{p1, p2}
}

func (b *BattleIDAndRoles) Get(battleID string) []string {
	return b.m[battleID]
}

func (b *BattleIDAndRoles) Delete(battleID string) {
	delete(b.m, battleID)
}
