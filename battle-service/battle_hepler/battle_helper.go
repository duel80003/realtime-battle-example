package battle_hepler

import (
	proto "battle-service/proto/gen"
	"math/rand"
	"time"
)

func StartBattle(roles []*Role) (res *proto.RealtimeBattleResponse) {
	res = &proto.RealtimeBattleResponse{}
	res.RoleInfos = []*proto.RoleInfo{
		Attack(roles[0], roles[1]),
		Attack(roles[1], roles[0]),
	}
	res.IsEnd = roles[0].IsDead || roles[1].IsDead
	return
}

func GameContinue(roles []*Role) bool {
	for _, v := range roles {
		if v.IsDead {
			return false
		}
	}
	return true
}

func Attack(attacker, defender *Role) *proto.RoleInfo {
	if attacker.IsDead {
		return nil
	}
	name, dmg := getDamage(attacker)
	isMissed := !isHit(defender.Luck)
	if !isMissed {
		defender.HP -= dmg
		defender.IsDead = isDead(defender.HP)
	}
	logger.Debugf("%s attacking %s, %s, dmg: %d", attacker.ID, defender.ID, name, dmg)
	return &proto.RoleInfo{
		RoleId:         attacker.ID,
		RoleName:       attacker.Name,
		TargetId:       defender.ID,
		TargetRoleName: defender.Name,
		SkillName:      name,
		Damage:         dmg,
		IsMissed:       isMissed,
	}
}

func randSkill(skills []*Skill) *Skill {
	index := rand.Intn(len(skills) - 0)
	return skills[index]
}

func isHit(defenderLuck int32) bool {
	value := rand.Intn(100-1) + 1
	return int32(value) < (100 - defenderLuck)
}

func isDead(hp int32) bool {
	return hp <= 0
}

func inCoolDown(nextReleaseTime int64) bool {
	return nextReleaseTime > time.Now().UnixMilli()
}

func computeNextReleaseTime(t int64) int64 {
	return time.Now().UnixMilli() + t*1000
}

func getDamage(attacker *Role) (string, int32) {
	skill := randSkill(attacker.Skills)
	logger.Debugf("skill NextReleaseTime: %d", skill.NextReleaseTime)
	if inCoolDown(skill.NextReleaseTime) {
		if inCoolDown(attacker.NextAttackTime) {
			return Standby, 0
		}
		attacker.NextAttackTime = computeNextReleaseTime(1)
		return NormalAttack, attacker.Atk
	}
	skill.NextReleaseTime = computeNextReleaseTime(skill.FreezingTime)
	logger.Debugf("now %d, next release time %d", time.Now().UnixMilli(), skill.NextReleaseTime)
	return skill.Name, skill.Damage
}
