package gs

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

//GameType 为了防止包的依赖,在这里不直接让gs包的GameType和client包的GameType用type保持一致
//应该手动保持一致.具体的做法是找到自动生成的client.pb.go,然后从中找到GameType并手动复制
type GameType int32

const (
	NormalGameType GameType = 0
)

//GameTypeMaxPlayer 各种游戏类型的玩家数目,用于匹配对应数目的玩家进入游戏
var GameTypeMaxPlayer = map[GameType]int{
	NormalGameType: 3, //普通游戏
}

type RoomInfo struct {
	Ip     string //gs的ip
	Port   string //用于跟Client连接的kcp/tcp端口
	RoomId int32  //房间id
}

//GameAccountInfo 对局结算信息
type GameAccountInfo struct {
	StartTime        int64
	EndTime          int64
	PlayerAccountMap map[int32]*PlayerAccountInfo
}

//PlayerAccountInfo 玩家结算信息
type PlayerAccountInfo struct {
	Id        int32 //玩家id
	KillNum   int32 //击杀数
	AliveTime int64 //生存时间
}

type RoomConnectionInfo struct {
	Id         int32
	ClientAddr string
}

//PlayerInfo 用于normal_game一开始玩家进入游戏时的身份校验以及玩家的姓名展示
type PlayerInfo struct {
	PlayerId int32
	NickName string
}

type CreateNormalGameRequest struct {
	PlayerInfo []*PlayerInfo
	GameId     string //作为etcd的路径,游戏结束之后由server通知hall,方便hall落库
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
