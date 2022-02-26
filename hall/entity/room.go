package entity

import (
	"hash/crc32"
	"melee_game_hall/api/gs"
)

/**
*@Author Sly
*@Date 2022/2/23
*@Version 1.0
*@Description:
 */

type RoomStatus int

const (
	RoomIdleStatus RoomStatus = 1
)

type RoomInfo struct {
	Id     uint64     //这个room在hall上的id
	Ip     string     //room_server的ip
	Port   string     //room_server的port
	RoomId int32      //room_server的room_id
	status RoomStatus //这个room的状态
}

func RoomInfoFromGS(info *gs.RoomInfo) *RoomInfo {
	ri := new(RoomInfo)
	ri.Id = CountIdFromGS(info)
	ri.Ip = info.Ip
	ri.Port = info.Port
	ri.RoomId = info.RoomId
	ri.status = RoomIdleStatus
	return ri
}

//CountIdFromGS 通过game_server的ip,port,结合game_server的roomId计算出一个某个room在hall这里的id
func CountIdFromGS(info *gs.RoomInfo) uint64 {
	key := []byte(info.Ip + info.Port)
	h := crc32.ChecksumIEEE(key)
	return uint64(h)<<32 + uint64(info.RoomId)
}
