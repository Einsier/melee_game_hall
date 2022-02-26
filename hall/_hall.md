# 大厅服务器功能及逻辑

## A 客户端

### 1.登录功能
前端使用`Serve`这个grpc,发送`LoginRequest`给大厅服务器,,其中包含`phone`和`password`.
```
message LoginRequest{
  string phoneNum = 1;
  string password = 2;
}
```

后端通过数据库连接模块,
查询数据库,并返回`LoginResponse`
```
message LoginResponse{
  bool  ok = 1;     //登录是否成功
  LoginErrorType error = 2; //如果不成功返回错误信息
  PlayerInfo pInfo = 3; //如果成功返回玩家信息
}
```

如果校验通过,在`PlayerInfo`字段中包含玩家信息(下称PlayerInfo),同时设置ok字段为true,返回给前端后,通过玩家信息
生成一个`HallPlayer`,作为当前加入大厅的玩家实体.并且将该玩家的grpc stream保存到`HallPlayer`结构中.
```
message PlayerInfo{
  int32 PlayerId = 1;
  string NickName = 2;
  int32 GameCount = 3;		  //参与游戏数
  int32 KillNum = 4;				//总击杀数
  int32 MaxKill = 5;			  //最高单局击杀数
}
```

如果校验未通过,在LoginErrorType返回错误信息,同时ok字段为false.如果grpc连接有问题,通过grpc函数的error返回值返回(前端可以直接展示error的text)
```
enum LoginErrorType{
  PhoneNotExist = 0;    //手机号未注册
  WrongPassword = 1;    //密码不正确
  AccountNotExist = 2;  //账户不存在
  DBInnerError   = 3;   //数据库内部错误
  HallInnerError = 4;
}
```

### 2.排队功能
前端使用使用`Serve`这个grpc,发出`StartQueueRequest`给大厅服务器:
```
message StartQueueRequest {
  int32     playerId = 1;   //待排队的玩家id
  GameType  gameType = 2;   //游戏类型
}
```
后端不做响应.前端可以开始显示排队计时画面.

等玩家要玩的游戏类型的人凑够之后,成功开启房间/排队超时/`game_server`通信异常/`game_room`无法正常开启时,
后端通过`Serve`这个grpc向前端发送`StartQueueResponse`作为排队的响应:
```
message StartQueueResponse{
  bool                  ok = 1;
  GameServerConnection  gameServerConnection = 2; //如果排队成功,返回gameServer的联系方式
  QueueErrorType        error = 3;    //如果排队失败,返回错误信息
}
```

如果排队失败,前端可以向玩家展示失败信息,如果成功,则`gameServerConnection`字段存放了`game_server`的联系方式,玩家可以向
`game_server`发送kcp报文(可以见`melee_game_server`中的文档),以此来加入对局房间.


## B game_server

### 1.开启
