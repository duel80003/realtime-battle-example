package battle_hepler

type Skill struct {
	ID              string `json:"ID"`
	Name            string `json:"name"`
	Damage          int32  `json:"damage"`
	FreezingTime    int64  `json:"freezingTime"`
	NextReleaseTime int64  `json:"_"`
}
