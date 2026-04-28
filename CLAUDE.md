# Ricchi-room 项目文档

## 项目概述
Ricchi-room是一款日本麻将（立直麻将）的多人牌桌记分微信小程序，支持4人同时在线记分、支付和结算。

## 技术栈
- 后端：Golang + GIN + WebSocket
- 前端：微信小程序 (WXSS, WXML, JS)
- 日志：zap
- 配置：viper

## 目录结构
```
ricchi-room/
├── CLAUDE.md              # 项目文档
├── DESIGN.md             # 架构设计
├── backend/              # 后端代码
│   ├── main.go          # 入口
│   ├── config/         # 配置
│   ├── models/         # 数据模型
│   ├── handlers/       # HTTP处理
│   ├── websocket/      # WebSocket处理
│   └── go.mod
└── frontend/           # 微信小程序
    ├── app.js
    ├── app.json
    ├── app.wxss
    └── pages/
        ├── index/      # 首页
        └── table/    # 牌桌页面
```

## 后端API

### REST API
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/table/create | 创建牌桌 |
| POST | /api/table/join | 加入牌桌 |
| POST | /api/table/leave | 离开牌桌 |
| POST | /api/table/pay | 支付点数 |
| POST | /api/table/settle | 结算 |
| POST | /api/table/dismiss | 解散房间 |
| POST | /api/table/reconnect | 重连恢复 |
| GET | /api/table/:id | 获取牌桌状态 |

### WebSocket
| 路径 | 说明 |
|------|------|
| /ws/table/:id | 实时状态同步 |

## 数据模型

### Table (牌桌)
- ID: 4位数字
- Players: 玩家列表 (最多4人)
- Status: 状态 (waiting/playing/settled)
- CreatedAt: 创建时间

### Player (玩家)
- OpenID: 微信ID
- Nickname: 昵称
- Points: 点数 (初始25000)
- IsOwner: 是否房主
- IsConnected: 是否在线

## 业务规则
1. 每人初始25000点
2. 支付必须是100的倍数
3. 点数<0触发击飞结算
4. 全局最多100桌
5. 所有人离开时牌桌解散

## 配色
- #c6ffdd (浅绿)
- #fbd786 (淡黄)
- #f7797d (珊瑚红)