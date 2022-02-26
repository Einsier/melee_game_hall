package database

/**
*@Author Sly
*@Date 2022/2/18
*@Version 1.0
*@Description:
 */

//IsAccountLegalRequest 登录
type IsAccountLegalRequest struct {
	Phone    string
	Password string
}

type IsAccountLegalResponse struct {
	Ok       bool
	Error    ErrorType //失败
	PlayerId int32
}

type ErrorType int

const (
	PhoneNotExist  ErrorType = 1
	WrongPassword  ErrorType = 2
	PlayerNotExist ErrorType = 3
	DBInnerError   ErrorType = 4
)

type PlayerInfo struct {
	PlayerId  int32
	NickName  string
	GameCount int32 //参与游戏数
	KillNum   int32 //总击杀数
	MaxKill   int32 //最高单局击杀数
}

//SearchPlayerInfoRequest 查找用户信息
type SearchPlayerInfoRequest struct {
	PlayerId int32
}

type SearchPlayerInfoResponse struct {
	Ok    bool
	Error ErrorType   //失败
	Info  *PlayerInfo //成功
}

//UpdatePlayerInfoRequest 玩家信息更新
type UpdatePlayerInfoRequest struct {
	Info *PlayerInfo
}

type UpdatePlayerInfoResponse struct {
	Ok    bool
	Error ErrorType
	Info  *PlayerInfo //如果失败,返回原来的info
}

//SingleGameInfo 单局游戏结算信息
type SingleGameInfo struct {
	Players   []int32 //参加游戏的玩家id
	StartTime int64   //游戏开始时间
	EndTime   int64   //游戏结束时间
}

type AddSingleGameInfoRequest struct {
	Info *SingleGameInfo
}

type AddSingleGameInfoResponse struct {
	Ok    bool
	Error ErrorType
}
