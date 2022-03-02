package main

import (
	"flag"
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

func main() {
	flag.Parse()
	_ = hall.NewHall("0.0.0.0" + *clientGrpcPortFlag)
	logger.Info("hall开始运行")
	time.Sleep(100 * time.Minute)
}
