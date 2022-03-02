package hall

import (
	"melee_game_hall/api/database"
	"strconv"
)

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

var DB *DBProxy

type DBProxy struct {
}

var TestPlayerArrangeIndex = int32(0)

//IsAccountLegal 通过数据库,根据账户密码验证用户身份,如果服务器连接有问题返回nil,其余返回非空的 IsAccountLegalResponse
func (D *DBProxy) IsAccountLegal(req *database.IsAccountLegalRequest) *database.IsAccountLegalResponse {
	//如果有错误,让返回值为nil
	TestPlayerArrangeIndex++
	return &database.IsAccountLegalResponse{
		Ok:       true,
		Error:    0,
		PlayerId: TestPlayerArrangeIndex,
	}
}

//SearchPlayerInfo 通过数据库,根据playerId查询用户数据,如果服务器连接有问题返回nil,其余返回非空的 SearchPlayerInfoResponse
func (D *DBProxy) SearchPlayerInfo(req *database.SearchPlayerInfoRequest) *database.SearchPlayerInfoResponse {
	return &database.SearchPlayerInfoResponse{
		Ok:    true,
		Error: 0,
		Info: &database.PlayerInfo{
			PlayerId:  req.PlayerId,
			NickName:  "player" + strconv.Itoa(int(req.PlayerId)),
			GameCount: req.PlayerId * 2,
			KillNum:   req.PlayerId * 3,
			MaxKill:   req.PlayerId * 4,
		},
	}
}

func (D *DBProxy) UpdatePlayerInfo(req *database.UpdatePlayerInfoRequest) *database.UpdatePlayerInfoResponse {
	return nil
}

func (D *DBProxy) AddSingleGameInfo(req *database.AddSingleGameInfoRequest) *database.AddSingleGameInfoResponse {
	return nil
}
