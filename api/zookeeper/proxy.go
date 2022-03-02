package zookeeper

/**
*@Author Sly
*@Date 2022/2/25
*@Version 1.0
*@Description:
 */

type Proxy interface {
	//GetGameServer 通过注册中心获取GameServer的地址,如果有错误,返回错误类型,注意ip为类似"192.168.1.1",port类似":8000"
	GetGameServer() (ip, rpcPort, clientPort string, e error)
}
