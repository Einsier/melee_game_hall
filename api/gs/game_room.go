package gs

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

const MaxNormalGamePlayerNum = int32(2) //normal版本的最大玩家人数

//GameType 为了防止包的依赖,在这里不直接让gs包的GameType和client包的GameType用type保持一致
//应该手动保持一致.具体的做法是找到自动生成的client.pb.go,然后从中找到GameType并手动复制
type GameType int32

const (
	NormalGameType GameType = 0
)

var GameTypeMaxPlayer = map[GameType]int{
	NormalGameType: 2,
}

type RoomInfo struct {
	Ip     string //gs的ip
	Port   string //用于跟Client连接的kcp/tcp端口
	RoomId int32  //房间id
}

type RoomConnectionInfo struct {
	Id         int32
	ClientPort string
}

//PlayerInfo 暂时只有playerId,用于normal_game一开始玩家进入游戏时的身份校验
type PlayerInfo struct {
	PlayerId int32
}

type CreateNormalGameRequest struct {
	PlayerInfo []*PlayerInfo
}

type StartNormalGameRequest struct {
	RoomId int32
}

type CreateNormalGameResponse struct {
	Ok             bool
	ConnectionInfo *RoomConnectionInfo
}

type StartNormalGameResponse struct {
	Ok bool
}

type RoomStatus int

const (
	RoomInitStatus       RoomStatus = 1
	RoomStartStatus      RoomStatus = 2
	RoomDestroyingStatus RoomStatus = 3
)

type DestroyGameRoomRequest struct {
	RoomId int32
}
type DestroyGameRoomResponse struct {
	Status RoomStatus
	Ok     bool
}
