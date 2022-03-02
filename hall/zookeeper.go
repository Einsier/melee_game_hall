package hall

import "melee_game_hall/configs"

/**
*@Author Sly
*@Date 2022/2/25
*@Version 1.0
*@Description:
 */

var ZooKeeper = NewZooKeeperProxy(configs.ZookeeperAddr)

type ZooKeeperProxy struct {
	addr string
}

//GetGameServer 通过注册中心获取GameServer的地址,如果有错误,返回错误类型,注意ip为类似"192.168.1.1",port类似":8000"
func (z *ZooKeeperProxy) GetGameServer() (string, string, string, error) {
	//todo 改成从注册中心取出
	return "127.0.0.1", ":8000", ":8001", nil
}

func NewZooKeeperProxy(addr string) *ZooKeeperProxy {
	return &ZooKeeperProxy{addr: addr}
}
