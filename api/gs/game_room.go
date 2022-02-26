package gs

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

const MaxNormalGamePlayerNum = int32(10) //normal版本的最大玩家人数

//GameType 为了防止包的依赖,在这里不直接让gs包的GameType和client包的GameType用type保持一致
//应该手动保持一致.具体的做法是找到自动生成的client.pb.go,然后从中找到GameType并手动复制
type GameType int32

const (
	NormalGameType GameType = 0
)

var GameTypeMaxPlayer = map[GameType]int{
	NormalGameType: 10,
}

type RoomInfo struct {
	Ip     string
	Port   string
	RoomId int32
}

type RoomConnectionInfo struct {
	Id int32
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
