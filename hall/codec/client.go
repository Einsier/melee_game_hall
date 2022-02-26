package codec

import (
	"melee_game_hall/api/client"
	"melee_game_hall/api/database"
	"melee_game_hall/hall/entity"
)

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

//NewLoginResponse 通过数据库的校验返回信息构造发给Client的顶层message
func NewLoginResponse(ok bool, err database.ErrorType, info *entity.PlayerInfo) *client.HToC {
	e := client.LoginErrorType(0)
	if !ok {
		//如果有错
		switch err {
		case database.PhoneNotExist:
			e = client.LoginErrorType_PhoneNotExist
		case database.WrongPassword:
			e = client.LoginErrorType_WrongPassword
		case database.DBInnerError:
			e = client.LoginErrorType_DBInnerError
		case database.PlayerNotExist:
			e = client.LoginErrorType_AccountNotExist
		default:
			e = client.LoginErrorType_HallInnerError
		}
		return &client.HToC{
			MsgType: client.HToCType_LoginResp,
			LoginResponse: &client.LoginResponse{
				Ok:    ok,
				Error: e,
				PInfo: nil,
			},
		}
	} else {
		return &client.HToC{
			MsgType: client.HToCType_LoginResp,
			LoginResponse: &client.LoginResponse{
				Ok:    ok,
				Error: 0,
				PInfo: info.ToClient(),
			},
		}
	}
}

func NewStartQueueResponse(ok bool, info *entity.RoomInfo, errorType client.QueueErrorType) *client.HToC {
	if !ok {
		//如果排队失败,返回失败原因,设置Ok为false
		return &client.HToC{
			MsgType: client.HToCType_StartQueuingResp,
			StartQueueResponse: &client.StartQueueResponse{
				Ok:                   false,
				GameServerConnection: nil,
				Error:                errorType,
			},
		}
	} else {
		return &client.HToC{
			MsgType: client.HToCType_StartQueuingResp,
			StartQueueResponse: &client.StartQueueResponse{
				Ok: true,
				GameServerConnection: &client.GameServerConnection{
					Ip:     info.Ip,
					Port:   info.Port,
					RoomId: info.RoomId,
				},
				Error: 0,
			},
		}
	}
}

func NewStopQueuingResponse(ok bool) *client.HToC {
	return &client.HToC{
		MsgType:             client.HToCType_StopQueuingResp,
		StopQueuingResponse: &client.StopQueuingResponse{Ok: ok},
	}
}
