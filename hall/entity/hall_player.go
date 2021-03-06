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

	QueueType GameType //如果当前状态为排队,此字段表示正在拍的队列的游戏类型

	sLock  sync.Mutex
	status PlayerStatus //玩家当前在大厅中的状态

	PInfo *PlayerInfo               //玩家信息
	rInfo *RoomInfo                 //如果如果玩家正在游戏,此字段保存正在游戏的信息
	Conn  client.Client_ServeServer //用于联系玩家的grpc连接

	ALock          sync.Mutex
	WaitingAccount map[string]struct{} //用于记录当前还有哪些对局结算信息没有被落库
	Quit           bool                //玩家是否断开连接
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

func NewHallPlayer(pInfo *PlayerInfo, server client.Client_ServeServer) *HallPlayer {
	return &HallPlayer{
		PlayerId:       pInfo.PlayerId,
		status:         PlayerIdle,
		PInfo:          pInfo,
		rInfo:          nil,
		Conn:           server,
		WaitingAccount: make(map[string]struct{}),
		Quit:           false,
	}
}
