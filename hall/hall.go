package hall

import (
	"errors"
	"google.golang.org/grpc"
	"melee_game_hall/api/client"
	"melee_game_hall/api/gs"
	"melee_game_hall/configs"
	"melee_game_hall/hall/codec"
	"melee_game_hall/hall/entity"
	"melee_game_hall/plugins/logger"
	"net"
	"sync"
)

/**
*@Author Sly
*@Date 2022/2/18
*@Version 1.0
*@Description:1.0版本,暂时将game_server的地址写死,后期改用注册中心
 */

type ClientGrpcImpl struct {
	client.UnimplementedClientServer
}

type Hall struct {
	ClientGrpcImpl

	players map[int32]*entity.HallPlayer //当前在本大厅的玩家
	pLock   sync.RWMutex

	rooms map[uint64]*entity.RoomInfo //本大厅向game_server开启的rooms的存根
	rLock sync.Mutex

	queue map[entity.GameType]map[int32]struct{} //key为游戏类型,value为等待这个游戏类型开始的玩家.一种游戏一个排队序列
	qLock sync.Mutex

	clientGrpcAddr string //对客户端开放的grpc端口
}

//NewHall 初始化hall,并开启grpc服务
func NewHall(clientGrpcAddr string) *Hall {
	l, err := net.Listen("tcp", clientGrpcAddr)
	if err != nil {
		panic(":" + err.Error())
	}
	s := grpc.NewServer()

	//初始化hall
	hall := &Hall{
		ClientGrpcImpl: ClientGrpcImpl{},
		players:        make(map[int32]*entity.HallPlayer),
		pLock:          sync.RWMutex{},
		rooms:          make(map[uint64]*entity.RoomInfo),
		rLock:          sync.Mutex{},
		queue:          make(map[entity.GameType]map[int32]struct{}),
		qLock:          sync.Mutex{},
		clientGrpcAddr: clientGrpcAddr,
	}

	//将初始化的hall注册到grpc服务中
	client.RegisterClientServer(s, hall)

	logger.Infof("当前已在%s上开启客户端的grpc服务", clientGrpcAddr)
	go func() {
		e := s.Serve(l)
		if e != nil {
			panic("启动同客户端的grpc服务时出错:%s" + e.Error())
		}
	}()
	return hall
}

//AddHallPlayer 将HallPlayer添加到大厅中
func (h *Hall) AddHallPlayer(player *entity.HallPlayer) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	h.players[player.PlayerId] = player
}

//GetHallPlayer 从大厅中根据玩家id取出HallPlayer
func (h *Hall) GetHallPlayer(pId int32) *entity.HallPlayer {
	h.pLock.RLock()
	defer h.pLock.RUnlock()
	return h.players[pId]
}

//AddRoom 向大厅中添加一个通过本大厅开启的房间的房间信息(即game_room的信息)
func (h *Hall) AddRoom(info *entity.RoomInfo) {
	h.rLock.Lock()
	defer h.rLock.Unlock()
	h.rooms[info.Id] = info
}

//JoinQueue 玩家开始排队,凑齐对应人数(参见api/gs/game_room.go中的GameTypeMaxPlayer,保存每种游戏类型的最大支持玩家)
func (h *Hall) JoinQueue(gameType entity.GameType, pId int32) error {
	hallPlayer := h.GetHallPlayer(pId)
	if hallPlayer == nil {
		return errors.New("pId 非法")
	}

	//取出排队人数和游戏对局需要人数进行比较
	maxPlayer := gs.GameTypeMaxPlayer[gs.GameType(gameType)]
	h.qLock.Lock()

	//更改Player的状态必须要加锁进行
	hallPlayer.SetStatus(entity.PlayerQueuing)
	if h.queue[gameType] == nil {
		//如果当前队列为空,那么创建一下
		h.queue[gameType] = make(map[int32]struct{})
	}
	h.queue[gameType][pId] = struct{}{}
	logger.Infof("玩家%d开始排队\n", pId)
	if len(h.queue[gameType]) == maxPlayer {
		//如果凑齐人数的话,将所有正在排队的玩家取出,并且发给大厅服务器
		players := make([]int32, maxPlayer)
		temp := 0
		for key, _ := range h.queue[gameType] {
			//拷贝一下玩家id
			players[temp] = key
			temp++
		}
		delete(h.queue, gameType)
		h.qLock.Unlock()
		logger.Infof("玩家:%v 排队成功,准备进入游戏\n", players)

		pInfo := make([]*entity.PlayerInfo, maxPlayer)
		for i := 0; i < maxPlayer; i++ {
			p := h.GetHallPlayer(players[i])
			if p != nil {
				pInfo[i] = p.PInfo
			}
		}

		gameId := GetGameId()
		rInfo, err := CreateGameRoom(configs.GameServerRpcAddr, gameType, pInfo, gameId)
		if err != nil {
			//如果开启game_room失败,返回错误信息
			msg := codec.NewStartQueueResponse(false, nil, client.QueueErrorType_CanNotStartGameRoom)
			h.SendHallPlayersByPlayerId(players, msg)
			h.ChangePlayerStatusByPlayerId(players, entity.PlayerIdle)
			return err
		}
		//如果正常开启game_room,给排队的玩家返回对局信息
		msg := codec.NewStartQueueResponse(true, rInfo, 0)
		ok := h.SendHallPlayersByPlayerId(players, msg)
		if !ok {
			//如果出现了队列中的玩家失去连接等情况,给其它玩家发送失败,并且给gs发送删除game_room的请求
			_, _ = DestroyGameRoom(configs.GameServerRpcAddr, rInfo.RoomId)
			msg := codec.NewStartQueueResponse(false, nil, client.QueueErrorType_NoEnoughPlayer)
			h.SendHallPlayersByPlayerId(players, msg)
			h.ChangePlayerStatusByPlayerId(players, entity.PlayerIdle)
			return nil
		}

		//向etcd注册监听事件
		go h.ListenGameAccountEvent(gameId)
		logger.Infof("已注册gameId:%s 的etcd的监听事件", gameId)

		h.ChangePlayerStatusByPlayerId(players, entity.PlayerInGame)

		for _, pid := range players {
			h.AddPlayerWaitAccount(pid, gameId)
		}

		//前端收到之后应该有短暂进入游戏的动画展示用于拖延时间,因为需要通知gs开启房间
		err = StartGame(configs.GameServerRpcAddr, rInfo.RoomId)
		if err != nil {
			//如果开启房间失败,前端负责处理,重新连回大厅
			logger.Errorf("程序出错,未能StartGame")
		}
		return nil
	}
	h.qLock.Unlock()
	return nil
}

//ChangePlayerStatusByPlayerId 互斥的根据玩家id切片,更改玩家状态
func (h *Hall) ChangePlayerStatusByPlayerId(pId []int32, status entity.PlayerStatus) {
	for i := 0; i < len(pId); i++ {
		p := h.GetHallPlayer(pId[i])
		if p != nil {
			p.SetStatus(status)
		}
	}
}

//AddPlayerWaitAccount 为玩家添加对局等待事件
func (h *Hall) AddPlayerWaitAccount(pid int32, gameId string) {
	hp := h.GetHallPlayer(pid)
	if hp == nil {
		logger.Errorf("在为玩家添加等待对局信息时取到了不存在的用户.playerId:%d", pid)
		return
	}

	hp.ALock.Lock()
	defer hp.ALock.Unlock()
	hp.WaitingAccount[gameId] = struct{}{}
}

//DeleteHallPlayer 删除playerId为 pId 的在大厅的玩家.如果当前玩家正在排队,给队列加锁并且删除掉该玩家.
//不会判断当前有没有使用玩家的信息等.所以使用 GetHallPlayer 这个api的时候需要判断取出的玩家是不是nil
func (h *Hall) DeleteHallPlayer(pId int32) {
	hallPlayer := h.GetHallPlayer(pId)
	if hallPlayer == nil {
		//如果当前玩家已经被删除,那么直接返回
		return
	}
	h.qLock.Lock()
	//判断玩家是否在排队,如果在排队则从队列中删除该玩家
	if hallPlayer.GetStatus() == entity.PlayerQueuing {
		delete(h.queue[hallPlayer.QueueType], hallPlayer.PlayerId)
	}
	h.qLock.Unlock()

	//判断该玩家是否还有没有结算的信息,如果存在,那么暂时不删除,而是将Quit字段改为true
	hallPlayer.ALock.Lock()
	if len(hallPlayer.WaitingAccount) > 0 {
		hallPlayer.ALock.Unlock()
		hallPlayer.Quit = true
		return
	}
	hallPlayer.ALock.Unlock()

	//如果该玩家没有待结算的信息,直接删除该玩家
	h.DeleteHallPlayerAndPersist(pId)
}

//StopQueuing 取消排队
func (h *Hall) StopQueuing(pId int32, gameType entity.GameType) bool {
	//这里需要跟 JoinQueue 互斥,所以加锁.
	h.qLock.Lock()
	defer h.qLock.Unlock()

	//如果当前状态是Idle状态或者InGame状态(已经排队成功,没来得及取消),那么返回失败(false)
	//如果当前状态是Queuing状态,因为加锁,不会有新玩家加入同一个队列,也就不会凑齐人数发给game_server
	status := h.GetHallPlayer(pId).GetStatus()
	if status != entity.PlayerQueuing {
		return false
	}
	if h.queue[gameType] == nil {
		return false
	}
	_, ok := h.queue[gameType][pId]
	if !ok {
		return false
	}
	delete(h.queue[gameType], pId)
	h.GetHallPlayer(pId).SetStatus(entity.PlayerIdle)
	return true
}
