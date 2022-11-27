package battle_hepler

const (
	NormalAttack = "普攻"
	Standby      = "待機"
)

type RoleState int32

type GameInfo struct {
	IsDead         bool      `json:"-"`
	State          RoleState `json:"-"`
	NextAttackTime int64     `json:"-"`
}

type Role struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	HP     int32    `json:"hp"`
	Atk    int32    `json:"atk"`
	Luck   int32    `json:"luck"`
	Skills []*Skill `json:"skills"`
	GameInfo
}
