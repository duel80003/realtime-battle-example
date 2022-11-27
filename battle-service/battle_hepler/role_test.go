package battle_hepler

import (
	"battle-service/tools"
	"github.com/joho/godotenv"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load("../.env")
	if err != nil {
		panic("env file lost")
	}
	tools.LogInit()
	setFilePathForTest()
	ReadData()
}

func TestAttack(t *testing.T) {
	t.Logf("test start")
	roles := []*Role{
		GetRole("1001"),
		GetRole("1002"),
	}

	for GameContinue(roles) {
		time.Sleep(time.Millisecond * 100)
		res := StartBattle(roles)
		if res.IsEnd {
			t.Logf("test end")
		}
	}
	t.Logf("role %s, state; %v", roles[0].Name, roles[0].IsDead)
	t.Logf("role %s, state; %v", roles[1].Name, roles[1].IsDead)
}
