package hall

/**
*@Author Sly
*@Date 2022/3/17
*@Version 1.0
*@Description:用于同etcd进行连接
 */

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"melee_game_hall/api/database"
	"melee_game_hall/api/gs"
	"melee_game_hall/configs"
	"melee_game_hall/metrics"
	"melee_game_hall/plugins/logger"
	"melee_game_hall/utils"
	"strconv"
	"time"
)

var EtcdCli *clientv3.Client

func NewEtcdCli() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{configs.EtcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("创建etcd client时出错!err:%s", err.Error())
	}

	kv := clientv3.NewKV(cli)
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	_, err = kv.Get(ctx, "/test")
	if err != nil {
		log.Fatalf("连接etcd时出错!err:%s", err.Error())
	}
	return cli
}

//GetGameId 获取一个用于注册监听事件的名称,为了不与其他对局重复,使用nano time和random结合的方式
func GetGameId() string {
	return time.Now().Format(time.StampNano) + "-" + strconv.Itoa(rand.Int())
}

//ListenGameAccountEvent 等待对局服务器结束游戏,并且将对应的对局信息存入到etcd中.
func (h *Hall) ListenGameAccountEvent(gameId string) {
	//最多等待1h
	ctx, _ := context.WithTimeout(context.TODO(), time.Hour)

	//将前缀和gameId拼成一起,注册监听事件
	ch := EtcdCli.Watch(ctx, configs.AccountPath+"/"+gameId)

	//在这里进行等待
	response := <-ch

	logger.Infof("监听到对局:%s已结束,开始持久化操作", gameId)
	metrics.GaugeGameRoomCount.Dec()
	if len(response.Events) == 0 {
		//如果取到的值无效,或者超出最大的等待时间,那么不再等待,直接返回
		logger.Errorf("%s的对局结算信息有误", gameId)
		return
	}
	//如果取到的值有效,那么进行落库
	rawInfo := response.Events[0].Kv.Value

	gInfo := new(gs.GameAccountInfo)
	err := json.Unmarshal(rawInfo, gInfo)
	if err != nil {
		logger.Errorf("%s的对局结算信息编解码有误", gameId)
		return
	}

	//落库对局信息
	dbGameInfo := new(database.SingleGameInfo)
	dbGameInfo.StartTime = gInfo.StartTime
	dbGameInfo.EndTime = gInfo.EndTime
	pIds := make([]int32, len(gInfo.PlayerAccountMap))
	i := 0
	for pid, _ := range gInfo.PlayerAccountMap {
		pIds[i] = pid
		i++
	}
	dbGameInfo.Players = pIds
	info := DB.AddSingleGameInfo(&database.AddSingleGameInfoRequest{Info: dbGameInfo})
	if info == nil {
		//如果数据库连接错误,不落库当前对局的信息
		logger.Errorf("上传%s的对局结算信息时数据库代理连接错误", gameId)
		return
	}

	//落库玩家信息
	for pid, pInfo := range gInfo.PlayerAccountMap {
		h.LogPlayerAccountInfo(pid, pInfo, gameId)
	}
}

//DeleteHallPlayerAndPersist 删除玩家并且将数据持久化到db
func (h *Hall) DeleteHallPlayerAndPersist(pid int32) {
	hp := h.GetHallPlayer(pid)
	if hp == nil {
		//如果该玩家已经没有了
		logger.Errorf("删除玩家并且落库信息时有nil玩家指针出现,playerId:%d", pid)
		return
	}

	//将该玩家的信息进行落库
	req := new(database.UpdatePlayerInfoRequest)
	req.Info = new(database.PlayerInfo)
	req.Info.PlayerId = pid
	req.Info.NickName = hp.PInfo.NickName

	hp.PInfo.InfoLock.Lock()
	//填充落库信息,虽然不用加锁,但是保险起见还是加上了...
	req.Info.KillNum = hp.PInfo.KillNum
	req.Info.MaxKill = hp.PInfo.MaxKill
	req.Info.GameCount = hp.PInfo.GameCount
	hp.PInfo.InfoLock.Unlock()

	resp := DB.UpdatePlayerInfo(req)
	if resp.Ok == false {
		logger.Errorf("落库id为%d的玩家的信息时出现错误,错误id:%v", pid, resp.Error)
	}

	//从大厅中删除该玩家
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.players, pid)
}

//LogPlayerAccountInfo 落库玩家信息
func (h *Hall) LogPlayerAccountInfo(pid int32, info *gs.PlayerAccountInfo, gameId string) {
	hp := h.GetHallPlayer(pid)
	if hp == nil {
		//检查空指针
		logger.Errorf("落库时有nil玩家指针出现,playerId:%d,account info:%+v", pid, info)
		return
	}

	//更新玩家信息
	hp.PInfo.InfoLock.Lock()
	hp.PInfo.KillNum += info.KillNum
	hp.PInfo.MaxKill = utils.Int32Max(hp.PInfo.MaxKill, info.KillNum)
	hp.PInfo.GameCount++
	hp.PInfo.InfoLock.Unlock()

	//从玩家等待结算中删除该对局
	hp.ALock.Lock()
	defer hp.ALock.Unlock()
	delete(hp.WaitingAccount, gameId)
	if len(hp.WaitingAccount) == 0 && hp.Quit == true {
		h.DeleteHallPlayerAndPersist(pid)
	}
}
