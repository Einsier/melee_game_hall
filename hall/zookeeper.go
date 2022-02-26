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

func (z *ZooKeeperProxy) GetGameServer() (string, string, error) {
	return "0.0.0.0", ":7142", nil
}

func NewZooKeeperProxy(addr string) *ZooKeeperProxy {
	return &ZooKeeperProxy{addr: addr}
}
