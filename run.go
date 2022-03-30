package main

import (
	"flag"
	"melee_game_hall/configs"
	"melee_game_hall/hall"
	"melee_game_hall/plugins/logger"
	"time"
)

/**
*@Author Sly
*@Date 2022/2/27
*@Version 1.0
*@Description:
 */

var clientGrpcPortFlag = flag.String("clientGrpcPort", ":9000", "set the port of grpc in order to communicate with clients")
var gsRpcAddrFlag = flag.String("gsRpcAddr", "localhost:8000", "set the addr of rpc in order to communicate with game server")
var dbProxyAddrFlag = flag.String("dbProxyAddr", "42.192.200.194:32002", "set the database proxy's addr")
var etcdAddrFlag = flag.String("etcdAddr", "42.192.200.194:2379", "set the address of etcd")
var testFlag = flag.Bool("t", false, "if this is a local test")

func ParseFlags() {
	flag.Parse()
	if *testFlag {
		//如果当前是本机测试
		*clientGrpcPortFlag = ":9000"
		*gsRpcAddrFlag = "localhost:8000"
		*etcdAddrFlag = "42.192.200.194:2379"
		*dbProxyAddrFlag = "1.116.109.113:1234"
	}
}

func main() {
	ParseFlags()
	_ = hall.NewHall("0.0.0.0" + *clientGrpcPortFlag)
	configs.GameServerRpcAddr = *gsRpcAddrFlag
	configs.DBProxyAddr = *dbProxyAddrFlag
	configs.EtcdAddr = *etcdAddrFlag
	hall.EtcdCli = hall.NewEtcdCli()
	hall.DB = hall.NewDBProxy(configs.DBProxyAddr)
	logger.Info("hall开始运行")
	time.Sleep(100 * time.Minute)
}
