package database

/**
*@Author Sly
*@Date 2022/2/19
*@Version 1.0
*@Description:
 */

type DB interface {
	IsAccountLegal(req *IsAccountLegalRequest) *IsAccountLegalResponse
	SearchPlayerInfo(req *SearchPlayerInfoRequest) *SearchPlayerInfoResponse
	UpdatePlayerInfo(req *UpdatePlayerInfoRequest) *UpdatePlayerInfoResponse
	AddSingleGameInfo(req *AddSingleGameInfoRequest) *AddSingleGameInfoResponse
}
