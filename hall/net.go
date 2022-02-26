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

func (h *Hall) SendHallPlayersByPlayerId(pId []int32, hToC *client.HToC) {
	for _, playerId := range pId {
		player := h.GetHallPlayer(playerId)
		if player != nil && player.Conn != nil {
			err := player.Conn.Send(hToC)
			if err != nil {
				logger.Errorf("玩家id:%d 连接异常", pId)
			}
		}
	}
}
