package hall

import (
	"errors"
	"melee_game_hall/api/client"
	"melee_game_hall/api/gs"
	"melee_game_hall/hall/codec"
	"melee_game_hall/hall/entity"
	"melee_game_hall/plugins/logger"
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
}

func NewHall() *Hall {
	//todo init grpc
	return &Hall{
		ClientGrpcImpl: ClientGrpcImpl{},
		players:        make(map[int32]*entity.HallPlayer),
		pLock:          sync.RWMutex{},
		rooms:          make(map[uint64]*entity.RoomInfo),
		rLock:          sync.Mutex{},
		queue:          make(map[entity.GameType]map[int32]struct{}),
		qLock:          sync.Mutex{},
	}
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
		h.queue[gameType][pId] = struct{}{}
	} else {
		//如果不为空则将玩家加入
		h.queue[gameType][pId] = struct{}{}
		logger.Infof("玩家%d开始排队\n", pId)
	}
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

		//从注册中心拿到一个game_server的地址
		ip, port, err := ZooKeeper.GetGameServer()
		if err != nil {
			//如果zookeeper或者gs出现了问题,那么向所有排队的玩家的客户端发送错误报告
			logger.Errorf("Zookeeper加载gs发生异常:%v", err)
			msg := codec.NewStartQueueResponse(false, nil, client.QueueErrorType_CanNotStartGameServer)
			h.SendHallPlayersByPlayerId(players, msg)
			//有问题的话把玩家状态由排队改成空闲
			h.ChangePlayerStatusByPlayerId(players, entity.PlayerIdle)
			return err
		} else {
			//如果成功从注册中心拿到game_server地址,给全部玩家发送game_server的地址
			pInfo := make([]*entity.PlayerInfo, maxPlayer)
			for i := 0; i < maxPlayer; i++ {
				pInfo[i] = h.GetHallPlayer(players[i]).PInfo
			}
			err, rInfo := CreateGameRoom(ip, port, gameType, pInfo)
			if err != nil {
				//如果开启game_room失败,返回错误信息
				msg := codec.NewStartQueueResponse(false, nil, client.QueueErrorType_CanNotStartGameRoom)
				h.SendHallPlayersByPlayerId(players, msg)
				h.ChangePlayerStatusByPlayerId(players, entity.PlayerIdle)
				return err
			}
			//如果正常开启game_room,给排队的玩家返回对局信息
			msg := codec.NewStartQueueResponse(true, rInfo, 0)
			h.SendHallPlayersByPlayerId(players, msg)
			h.ChangePlayerStatusByPlayerId(players, entity.PlayerInGame)
			return nil
		}
	}
	h.qLock.Unlock()
	return nil
}

func (h *Hall) ChangePlayerStatusByPlayerId(pId []int32, status entity.PlayerStatus) {
	for i := 0; i < len(pId); i++ {
		h.GetHallPlayer(pId[i]).SetStatus(status)
	}
}

//DeleteHallPlayer 删除playerId为 pId 的在大厅的玩家.不会判断当前有没有使用玩家的信息等.所以使用 GetHallPlayer 这个api的时候
//需要判断取出的玩家是不是nil
func (h *Hall) DeleteHallPlayer(pId int32) {
	hallPlayer := h.GetHallPlayer(pId)
	if hallPlayer == nil {
		return
	}
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.players, pId)
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
