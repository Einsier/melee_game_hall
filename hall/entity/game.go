package entity

import (
	"melee_game_hall/api/client"
	"melee_game_hall/api/database"
)

/**
*@Author Sly
*@Date 2022/2/23
*@Version 1.0
*@Description:存放对局有关的实体
 */

//SingleGameAccount 存放游戏结算信息
type SingleGameAccount struct {
	Players   []int32 //参加游戏的玩家id
	StartTime int64   //游戏开始时间
	EndTime   int64   //游戏结束时间
}

func (info *SingleGameAccount) ToDB() *database.SingleGameInfo {
	return &database.SingleGameInfo{
		Players:   info.Players,
		StartTime: info.StartTime,
		EndTime:   info.EndTime,
	}
}

//GameType 游戏类型,和client.GameType保持一致
type GameType client.GameType
