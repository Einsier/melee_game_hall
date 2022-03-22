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
var dbProxyAddrFlag = flag.String("dbProxyAddr", "localhost:32002", "set the database proxy's addr")

func main() {
	flag.Parse()
	_ = hall.NewHall("0.0.0.0" + *clientGrpcPortFlag)
	configs.GameServerRpcAddr = *gsRpcAddrFlag
	configs.DBProxyAddr = *dbProxyAddrFlag
	logger.Info("hall开始运行")
	time.Sleep(100 * time.Minute)
}
