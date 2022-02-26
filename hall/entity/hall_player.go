package entity

import (
	"melee_game_hall/api/client"
	"sync"
)

/**
*@Author Sly
*@Date 2022/2/25
*@Version 1.0
*@Description:
 */

type PlayerStatus int

const (
	PlayerIdle    PlayerStatus = 1 //啥也没干
	PlayerInGame  PlayerStatus = 2 //正在游戏中
	PlayerQueuing PlayerStatus = 3 //排队中
)

//HallPlayer 在大厅中的玩家实体
type HallPlayer struct {
	PlayerId int32 //玩家id

	sLock  sync.Mutex
	status PlayerStatus //玩家当前在大厅中的状态

	PInfo *PlayerInfo               //玩家信息
	rInfo *RoomInfo                 //如果如果玩家正在游戏,此字段保存正在游戏的信息
	Conn  client.Client_ServeServer //用于联系玩家的grpc连接
}

//SetStatus 互斥的切换玩家状态
func (hp *HallPlayer) SetStatus(status PlayerStatus) {
	hp.sLock.Lock()
	defer hp.sLock.Unlock()
	hp.status = status
}

//GetStatus 获取玩家状态
func (hp *HallPlayer) GetStatus() PlayerStatus {
	hp.sLock.Lock()
	defer hp.sLock.Unlock()
	return hp.status
}
