# 棋牌游戏中心 服务端目录

所有棋牌共用游戏服务器项目

## 目录说明

internal 游戏逻辑主目录
--config 游戏可变参数配置文件
--const  游戏常量，枚举定义
--game   游戏核心玩法
  ---card 麻将牌,算法管理
  ---protoapi 游戏业务逻辑
  ---room 房间管理
  ---seat 座位管理
  ---table 游戏牌桌
  ---timer 计时器相关
  --- utils 常用类
--logs 游戏日志生成文件夹
--player 玩家数据管理
  ---api 接口调用的 token 的合法性验证
vendor 是依赖库三方导入，不能修改

