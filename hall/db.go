package hall

import "melee_game_hall/api/database"

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

var DB *DBProxy

type DBProxy struct {
}

//IsAccountLegal 用于发送账号密码,返回验证信息等
func (D *DBProxy) IsAccountLegal(req *database.IsAccountLegalRequest) *database.IsAccountLegalResponse {
	//如果有错误,让返回值为nil
	return nil
}

func (D *DBProxy) SearchPlayerInfo(req *database.SearchPlayerInfoRequest) *database.SearchPlayerInfoResponse {
	return nil
}

func (D *DBProxy) UpdatePlayerInfo(req *database.UpdatePlayerInfoRequest) *database.UpdatePlayerInfoResponse {
	return nil
}

func (D *DBProxy) AddSingleGameInfo(req *database.AddSingleGameInfoRequest) *database.AddSingleGameInfoResponse {
	return nil
}
