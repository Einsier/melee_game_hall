
# 大厅服务器功能及逻辑

## 客户端(Client)

#### 1.登录功能

前端使用`Serve`这个grpc,发送`LoginRequest`给大厅服务器,,其中包含`phone`和`password`.
```
message LoginRequest{
  string phoneNum = 1;
  string password = 2;
}
```

后端通过数据库连接模块,查询数据库,并返回`LoginResponse`
```
message LoginResponse{
  bool  ok = 1;     //登录是否成功
  LoginErrorType error = 2; //如果不成功返回错误信息
  PlayerInfo pInfo = 3; //如果成功返回玩家信息
}
```

如果校验通过,在`PlayerInfo`字段中包含玩家信息(下称PlayerInfo),同时设置ok字段为true,返回给前端后,通过玩家信息生成一个`HallPlayer`,作为当前加入大厅的玩家实体.并且将该玩家的grpc stream保存到`HallPlayer`结构中.
```
message PlayerInfo{
  int32 PlayerId = 1;
  string NickName = 2;
  int32 GameCount = 3;      //参与游戏数
  int32 KillNum = 4;            //总击杀数
  int32 MaxKill = 5;           //最高单局击杀数
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

#### 2.排队功能

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

如果排队失败,前端可以向玩家展示失败信息,如果成功,则`gameServerConnection`字段存放了`game_server`的联系方式,玩家可以向`game_server`发送kcp报文(可以见`melee_game_server`中的文档),以此来加入对局房间.

#### 异常处理

1.玩家在大厅的时候直接退出了客户端:
`client`和`hall`之间的`grpc`双向`stream`会检测到`error`,并且加锁,删除`HallPlayer`信息.

2.玩家在排队的时候直接退出了客户端:
假设玩家a正在匹配10个人玩的`normal game`,没有匹配到的时候退出了游戏.同1中相同,玩家a退出后,删除玩家a的`HallPlayer`实体之后,会同时从a正在排队的队列中删除a(会加锁,使其同其它入队go程互斥)

3.在玩家a退出的时候没来得及删除a,而匹配成功,10个人进入游戏后a没连入game_server or 玩家可以正常连接hall,但是无法正常连接game_server:
hall正常给所有排队成功的玩家(包括a在内)发送排队成功响应,玩家a的流此时异常,通过它给a的客户端发送排队成功响应会产生`error`,但是这个`error`不会被上层处理.其他九个玩家从hall拿到game_server地址,成功进入游戏服务器(`game_server`中的某个`game_room`)后,game_server会因为没有按时正常开始(需要10个人都连上game_server才可以开始游戏)游戏而直接结束对局,给`hall`返回一个同正常对局结束相同的信息(通过rpc的方式)
其中包括了对局结束异常信息,这场对局会被入库时记录异常.这里的处理方式参考了英雄联盟.缺点是其他九个玩家就要因为a没进去等待一段时间.但是这种情况出现概率极小.

4.玩家在游戏过程中强行退出游戏:
由于当前`game_server`暂时还没有容灾备份等处理,所以玩家在游戏过程中退出,再重新打开客户端后,只能重新连接`hall`,而无法连接`game_server`.由于游戏启动时间和退出时间的原因,当前`hall`中该玩家存根`HallPlayer`大概率被删掉(即使没被删掉,玩家重新登录处理逻辑也一样).所以登录后会在大厅服务器中重新进入`PlayerIdle`状态,可以重新加入匹配.



## 游戏对局服务器(game_server)

#### 1.开启房间

以开启`normal_game`为例,在`normal_game`的匹配玩家数量达到开启游戏的人数(下面以`10`为例)之后,`hall`会通过`CreateGameRoom`(定义于`hall/gs.go`)这个函数申请建立一个房间,而如果创建失败,则返回给客户端错误,并可以让玩家重新排队.

如果开启房间没出现错误,那么给客户端返回`game_server`的联系方式和roomId后,再向服务器通过`StartGame`(定义于`hall/gs.go`)函数发送一个开始游戏请求,`game_server`中的对应`game_room`即开始进行游戏.

#### 2.删除房间

用于强行删掉gs上的房间,切断`game_server`通过`kcp`/`tcp`和客户端进行的连接.

#### 3.对局结算

用于对局信息的结算.(todo)

