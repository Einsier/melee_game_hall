package entity

import (
	"melee_game_hall/api/client"
	"melee_game_hall/api/database"
	"sync"
)

/**
*@Author Sly
*@Date 2022/2/23
*@Version 1.0
*@Description:
 */

//PlayerInfo 玩家的个人信息等
type PlayerInfo struct {
	InfoLock sync.Mutex

	PlayerId  int32
	NickName  string
	GameCount int32 //参与游戏数
	KillNum   int32 //总击杀数
	MaxKill   int32 //最高单局击杀数
}

func PlayerInfoFromDB(info *database.PlayerInfo) *PlayerInfo {
	return &PlayerInfo{
		PlayerId:  info.PlayerId,
		NickName:  info.NickName,
		GameCount: info.GameCount,
		KillNum:   info.KillNum,
		MaxKill:   info.MaxKill,
	}
}

func (pInfo *PlayerInfo) ToClient() *client.PlayerInfo {
	return &client.PlayerInfo{
		PlayerId:  int32(pInfo.PlayerId),
		NickName:  pInfo.NickName,
		GameCount: int32(pInfo.GameCount),
		KillNum:   int32(pInfo.KillNum),
		MaxKill:   int32(pInfo.MaxKill),
	}
}
