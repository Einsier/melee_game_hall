package hall

import (
	"errors"
	"fmt"
	"melee_game_hall/api/gs"
	"melee_game_hall/hall/entity"
	"melee_game_hall/plugins/logger"
	"net/rpc"
	"strings"
)

/**
*@Author Sly
*@Date 2022/2/25
*@Version 1.0
*@Description:用于rpc调用game_server,或者让game_server rpc调用自己
 */

func callRpc(rpcAddr, rpcName string, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", rpcAddr)
	if err != nil {
		return err
	}

	err = c.Call(rpcName, args, reply)
	if err != nil {
		return err
	}
	return nil
}

//CreateGameRoom 通知game_server开启一个房间,如果建立成功,则将roomInfo返还给用户,否则返回错误
//todo 改成集群的模式
func CreateGameRoom(gameServerRpcAddr string, gameType entity.GameType, info []*entity.PlayerInfo, gameId string) (*entity.RoomInfo, error) {
	switch gameType {
	case entity.GameType(gs.NormalGameType):
		//开启普通房间
		createReq := new(gs.CreateNormalGameRequest)
		createReq.GameId = gameId
		pIds := make([]*gs.PlayerInfo, len(info))
		for i := 0; i < len(info); i++ {
			pIds[i] = &gs.PlayerInfo{PlayerId: info[i].PlayerId}
		}
		createReq.PlayerInfo = pIds
		ret := new(gs.CreateNormalGameResponse)
		err := callRpc(gameServerRpcAddr, "GameServer.CreateNormalGameRoom", createReq, ret)
		if err != nil {
			return nil, err
		} else if !ret.Ok || ret.ConnectionInfo == nil {
			return nil, errors.New("调用game_server创建房间时出错")
		} else {
			logger.Infof("从rpc地址为:%s的GameServer处开启房间%d,tcp连接为:%s", gameServerRpcAddr, ret.ConnectionInfo.Id, ret.ConnectionInfo.ClientAddr)
			sp := strings.Split(ret.ConnectionInfo.ClientAddr, ":")
			return entity.RoomInfoFromGS(&gs.RoomInfo{
				Ip:     sp[0],
				Port:   ":" + sp[1],
				RoomId: ret.ConnectionInfo.Id,
			}), nil
		}
	default:
		logger.Errorf("出现了不在列表中的游戏类型:%d", gameType)
		return nil, fmt.Errorf("出现了不在列表中的游戏类型:%d", gameType)
	}
}

func StartGame(gsIP, gsPort string, roomId int32) error {
	startGameRep := gs.StartNormalGameRequest{RoomId: roomId}
	ret := new(gs.StartNormalGameResponse)
	err := callRpc(gsIP+gsPort, "GameServer.StartNormalGame", startGameRep, ret)
	if err != nil {
		return err
	} else {
		if ret.Ok == false {
			return fmt.Errorf("StartGame失败")
		}
		return nil
	}
}

func DestroyGameRoom(gsIP, gsPort string, roomId int32) (gs.RoomStatus, error) {
	destroyGameRoomReq := gs.DestroyGameRoomRequest{RoomId: roomId}
	ret := new(gs.DestroyGameRoomResponse)
	err := callRpc(gsIP+gsPort, "GameServer.DestroyGameRoom", destroyGameRoomReq, ret)
	if err != nil {
		return 0, err
	} else {
		return ret.Status, nil
	}
}
