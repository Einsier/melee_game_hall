package hall

import (
	"errors"
	"fmt"
	"melee_game_hall/api/gs"
	"melee_game_hall/hall/entity"
	"melee_game_hall/plugins/logger"
	"net/rpc"
)

/**
*@Author Sly
*@Date 2022/2/25
*@Version 1.0
*@Description:用于rpc调用game_server,或者让game_server rpc调用自己
 */

func callGsRpc(gsAddr, rpcName string, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", gsAddr)
	if err != nil {
		return err
	}

	err = c.Call(rpcName, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func CreateGameRoom(gsIP, gsPort string, gameType entity.GameType, info []*entity.PlayerInfo) (error, *entity.RoomInfo) {
	switch gameType {
	case entity.GameType(gs.NormalGameType):
		//开启普通房间
		createReq := new(gs.CreateNormalGameRequest)
		pIds := make([]*gs.PlayerInfo, len(info))
		for i := 0; i < len(info); i++ {
			pIds[i] = &gs.PlayerInfo{PlayerId: info[i].PlayerId}
		}
		createReq.PlayerInfo = pIds
		ret := new(gs.CreateNormalGameResponse)
		err := callGsRpc(gsIP+gsPort, "GameServer.CreateNormalGameRoom", createReq, ret)
		if err != nil {
			return err, nil
		} else if !ret.Ok || ret.ConnectionInfo == nil {
			return errors.New("调用game_server创建房间时出错"), nil
		} else {
			logger.Infof("从地址为:%s的GameServer处开启房间%d", gsIP+gsPort, ret.ConnectionInfo.Id)
			return nil, entity.RoomInfoFromGS(&gs.RoomInfo{
				Ip:     gsIP,
				Port:   gsPort,
				RoomId: ret.ConnectionInfo.Id,
			})
		}
	default:
		logger.Errorf("出现了不在列表中的游戏类型:%d", gameType)
		return fmt.Errorf("出现了不在列表中的游戏类型:%d", gameType), nil
	}
}
