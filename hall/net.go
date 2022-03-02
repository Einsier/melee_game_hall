package hall

import (
	"melee_game_hall/api/client"
	"melee_game_hall/plugins/logger"
)

/**
*@Author Sly
*@Date 2022/2/26
*@Version 1.0
*@Description:
 */

//SendHallPlayersByPlayerId 向切片中的玩家发送消息,先检查有没有连接问题,如果有连接问题返回error,不发送任何信息
func (h *Hall) SendHallPlayersByPlayerId(pId []int32, hToC *client.HToC) (ok bool) {
	//处理逻辑可以看最外层readme文件
	for _, playerId := range pId {
		player := h.GetHallPlayer(playerId)
		if player == nil || player.Conn == nil {
			return false
		}
	}

	for _, playerId := range pId {
		player := h.GetHallPlayer(playerId)
		if player != nil && player.Conn != nil {
			err := player.Conn.Send(hToC)
			if err != nil {
				logger.Errorf("玩家id:%d 连接异常", pId)
			}
		}
	}
	return true
}
