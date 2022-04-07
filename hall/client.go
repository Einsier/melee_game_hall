package hall

import (
	"errors"
	"melee_game_hall/api/client"
	"melee_game_hall/api/database"
	"melee_game_hall/hall/codec"
	"melee_game_hall/hall/entity"
	"melee_game_hall/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/2/18
*@Version 1.0
*@Description:和客户端的grpc接口的实现
 */

func (h *Hall) Serve(stream client.Client_ServeServer) error {
	var pInfo *entity.PlayerInfo
	var hallPlayer *entity.HallPlayer
	var hToC *client.HToC
	var pId int32
	//var rInfo *entity.RoomInfo
	for {
		//鉴定用户身份
		info, err := stream.Recv()

		//如果用户传来的首个包不是登录,那么直接简单的退出
		if err != nil || info.MsgType != client.CToHType_LoginReq || info.LoginRequest == nil {
			return errors.New("用户初始连接有误")
		}

		//判断用户合法性
		req := info.LoginRequest
		phone := req.GetPhoneNum()
		password := req.GetPassword()
		dbReq := &database.IsAccountLegalRequest{
			Phone:    phone,
			Password: password,
		}

		//通过数据库连接模块查找信息
		dbResp := DB.IsAccountLegal(dbReq)
		if dbResp == nil {
			logger.Errorf("数据库服务初始化有误")
			return errors.New("数据库服务初始化有误")
		}
		if !dbResp.Ok {
			//用户连接失败
			htoC := codec.NewLoginResponse(false, dbResp.Error, nil)
			err = stream.Send(htoC)
			if err != nil {
				logger.Errorf("hall服务器连接有误")
				return errors.New("hall服务器连接有误")
			}
		} else {
			//用户连接成功,dbResp中的PlayerId保存着玩家id
			pId = dbResp.PlayerId
			dbPlayerInfo := DB.SearchPlayerInfo(&database.SearchPlayerInfoRequest{PlayerId: pId})
			if dbPlayerInfo == nil {
				logger.Errorf("hall服务器内部错误")
				return errors.New("hall服务器内部错误")
			}
			pInfo = entity.PlayerInfoFromDB(dbPlayerInfo.Info)
			hallPlayer = entity.NewHallPlayer(pInfo, stream)
			hToC = codec.NewLoginResponse(true, 0, pInfo)
			err = stream.Send(hToC)
			if err != nil {
				logger.Errorf("hall服务器连接有误")
				return errors.New("hall服务器连接有误")
			}
			break
		}
	}
	if h.GetHallPlayer(pId) != nil {
		htoC := codec.NewLoginResponse(false, database.WrongPassword, nil)
		_ = stream.Send(htoC)
		logger.Errorf("playerId:%d重复登录", pId)
		return errors.New("玩家重复登录")
	}
	//代码执行到这里,用户已经登录成功,存储该用户,表示其上线
	h.AddHallPlayer(hallPlayer)
	logger.Infof("playerId:%d的player进入大厅", pId)

	//等待用户做事情
	for {
		info, err := stream.Recv()
		if err != nil {
			logger.Errorf("来自playerId:%d的hall服务器连接有误:%s,已在大厅中删除该玩家", pId, err.Error())
			h.DeleteHallPlayer(pId)
			return err
		}

		//处理客户端发来的信息
		switch info.MsgType {
		case client.CToHType_LoginReq:
			//收到已经登录的玩家的登录请求,不校验,回复一个登录成功请求,把玩家在大厅中的状态改为idle状态
			hToC = codec.NewLoginResponse(true, 0, pInfo)
			err = stream.Send(hToC)
			h.GetHallPlayer(pId).SetStatus(entity.PlayerIdle)
		case client.CToHType_QuitHallReq:
			//退出大厅请求,删除该用户在大厅的存根并return
			h.DeleteHallPlayer(pId)
			return nil
		case client.CToHType_GetPlayerInfoReq:
			//因为保持了一条长连接,不用校验用户发来的内容
			dbPlayerInfo := DB.SearchPlayerInfo(&database.SearchPlayerInfoRequest{PlayerId: pId})
			if dbPlayerInfo == nil {
				logger.Errorf("hall服务器内部错误:SearchPlayerInfo为nil")
				return errors.New("hall服务器内部错误:SearchPlayerInfo为nil")
			}
			pInfo = entity.PlayerInfoFromDB(dbPlayerInfo.Info)
			hallPlayer = entity.NewHallPlayer(pInfo, stream)
			hToC = codec.NewLoginResponse(true, 0, pInfo)
			err = stream.Send(hToC)
			if err != nil {
				logger.Errorf("发送NewLoginResponse时服务器连接有误")
				return errors.New("发送NewLoginResponse时服务器连接有误")
			}
		case client.CToHType_StartQueuingReq:
			if info.StartQueueRequest == nil {
				logger.Errorf("playerId:%d发来的StartQueueRequest有误", pId)
				continue
			}
			_ = h.JoinQueue(entity.GameType(info.StartQueueRequest.GameType), pId)
		case client.CToHType_StopQueuingReq:
			if info.StopQueuingRequest == nil {
				logger.Errorf("playerId:%d发来的StopQueuingRequest有误", pId)
			}
			ok := h.StopQueuing(pId, entity.GameType(info.StopQueuingRequest.GameType))
			err := stream.Send(codec.NewStopQueuingResponse(ok))
			if err != nil {
				logger.Errorf("hall服务器发送StopQueuing时连接有误")
				return errors.New("hall服务器发送StopQueuing时连接有误")
			}
		}
	}
}
