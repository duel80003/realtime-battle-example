package controllers

import "game-service/tools"

const TextMessage = 1

var (
	logger = tools.Logger
)

type RequestCode int32

const (
	RequestCodeMatch RequestCode = iota + 1
	RequestCodeBattle
)

type Request struct {
	PlayerID string `json:"playerId"`
}

type MatchResponse struct {
	BattleID  string `json:"battle_id"`
	PlayerID1 string `json:"playerId1"`
	PlayerID2 string `json:"playerId2"`
}

type BattleData struct {
	Text string `json:"text"`
}

type BattleResult struct {
	Winner string `json:"winner"`
}
