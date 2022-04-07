package hall

import (
	"log"
	"melee_game_hall/api/database"
	"melee_game_hall/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

//如果测试(不校验/落库信息),把下面这行等号之后,一直到type DBProxy上面一行注释掉
var DB *DBProxy

func NewDBProxy(addr string) *DBProxy {
	dbp := &DBProxy{addr: addr}
	resp := &database.IsAccountLegalResponse{}
	err := callRpc(addr, "HallHandler.IsAccountLegal", &database.IsAccountLegalRequest{
		Phone:    "17306409322",
		Password: "123456",
	}, resp)
	if err != nil {
		log.Fatalln("数据库代理模块连接失败...")
	}
	logger.Infof("数据库代理连接成功\n")
	return dbp
}

type DBProxy struct {
	addr string
}

//var TestPlayerArrangeIndex = int32(0)

//IsAccountLegal 通过数据库,根据账户密码验证用户身份,如果服务器连接有问题返回nil,其余返回非空的 IsAccountLegalResponse
func (dbp *DBProxy) IsAccountLegal(req *database.IsAccountLegalRequest) *database.IsAccountLegalResponse {
	resp := &database.IsAccountLegalResponse{}
	err := callRpc(dbp.addr, "HallHandler.IsAccountLegal", req, resp)
	if err != nil {
		logger.Errorf("数据库代理模块出现连接不上的情况:", err)
		return nil
	}
	return resp

	//todo 测试用,不连数据库判断玩家合法性
	/*TestPlayerArrangeIndex++
	return &database.IsAccountLegalResponse{
		Ok:       true,
		Error:    0,
		PlayerId: TestPlayerArrangeIndex,
	}*/
}

//SearchPlayerInfo 通过数据库,根据playerId查询用户数据,如果服务器连接有问题返回nil,其余返回非空的 SearchPlayerInfoResponse
func (dbp *DBProxy) SearchPlayerInfo(req *database.SearchPlayerInfoRequest) *database.SearchPlayerInfoResponse {

	resp := &database.SearchPlayerInfoResponse{}
	err := callRpc(dbp.addr, "HallHandler.SearchPlayerInfo", req, resp)
	if err != nil {
		logger.Errorf("数据库代理模块出现连接不上的情况:", err)
		return nil
	}
	return resp

	//todo 测试用,不连数据库判断玩家合法性
	/*return &database.SearchPlayerInfoResponse{
		Ok:    true,
		Error: 0,
		Info: &database.PlayerInfo{
			PlayerId:  req.PlayerId,
			NickName:  "player" + strconv.Itoa(int(req.PlayerId)),
			GameCount: req.PlayerId,
			KillNum:   req.PlayerId,
			MaxKill:   req.PlayerId,
		},
	}*/
}

func (dbp *DBProxy) UpdatePlayerInfo(req *database.UpdatePlayerInfoRequest) *database.UpdatePlayerInfoResponse {
	resp := new(database.UpdatePlayerInfoResponse)
	err := callRpc(dbp.addr, "HallHandler.UpdatePlayerInfo", req, resp)
	if err != nil {
		logger.Errorf("数据库代理模块出现连接不上的情况:", err)
		return nil
	}
	logger.Infof("落库玩家消息:%+v", req.Info)
	return resp
	//todo 测试用
	/*return &database.UpdatePlayerInfoResponse{
		Ok:    true,
		Error: 0,
		Info:  nil,
	}*/
}

func (dbp *DBProxy) AddSingleGameInfo(req *database.AddSingleGameInfoRequest) *database.AddSingleGameInfoResponse {
	logger.Infof("落库对局消息:%+v", req.Info)
	resp := new(database.AddSingleGameInfoResponse)
	err := callRpc(dbp.addr, "HallHandler.AddSingleGameInfo", req, resp)
	if err != nil {
		logger.Errorf("数据库代理模块出现连接不上的情况:", err)
		return nil
	}
	return resp
	/*	logger.Infof("落库对局消息:%+v", req.Info)
		return &database.AddSingleGameInfoResponse{
			Ok:    true,
			Error: 0,
		}*/
}
