package hall

import (
	"fmt"
	"melee_game_hall/api/gs"
	"melee_game_hall/hall/entity"
	"strconv"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/2/26
*@Version 1.0
*@Description:
 */

const MaxTestPlayerNum = 500

var testHall = NewHall("localhost:8000")

//testHallPlayerSlice 存储0~MaxTestPlayerNum,下标为x的存储playerId为x的Player
var testHallPlayerSlice = make([]*entity.HallPlayer, MaxTestPlayerNum)

func InitTest() {
	for i := 0; i < MaxTestPlayerNum; i++ {
		testHallPlayerSlice[i] = entity.NewHallPlayer(&entity.PlayerInfo{
			PlayerId:  int32(i),
			NickName:  "player" + strconv.Itoa(i),
			GameCount: int32(i),
			KillNum:   int32(i),
			MaxKill:   int32(i),
		}, nil)
		go testHall.AddHallPlayer(testHallPlayerSlice[i])
	}
}

func TestHall_BasicOption(t *testing.T) {
	InitTest()
	for i := 0; i < MaxTestPlayerNum; i++ {
		go func(id int32) {
			hallPlayer := testHall.GetHallPlayer(id)
			if hallPlayer == nil || hallPlayer.PlayerId != id {
				panic(fmt.Sprintf("id为%d的玩家信息错误!", id))
				return
			}
		}(int32(i))
	}
	for i := 0; i < MaxTestPlayerNum; i++ {
		go func(id int32, quit bool) {
			err := testHall.JoinQueue(entity.GameType(gs.NormalGameType), id)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			time.Sleep(1 * time.Millisecond)
			if quit {
				testHall.StopQueuing(id, entity.GameType(gs.NormalGameType))
			}
		}(int32(i), i%3 == 0)
	}
	time.Sleep(1 * time.Second)
}

func TestGetGameId(t *testing.T) {
	fmt.Printf("%v\n", time.Now().UnixNano())
	fmt.Printf("%v\n", GetGameId())
}
