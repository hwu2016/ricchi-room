package handlers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"ricchi-room/backend/models"
	ricchiWebsocket "ricchi-room/backend/websocket"
)

type PayRequest struct {
	TableID    string `json:"table_id"`
	FromOpenID string `json:"from_open_id"`
	ToOpenID   string `json:"to_open_id"`
	Amount     int    `json:"amount"`
}

func Pay(c *gin.Context) {
	var req PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Amount <= 0 || req.Amount%100 != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive multiple of 100"})
		return
	}

	mgr := models.GetTableManager()
	table := mgr.GetTable(req.TableID)
	if table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	fromPlayer := table.GetPlayer(req.FromOpenID)
	toPlayer := table.GetPlayer(req.ToOpenID)
	if fromPlayer == nil || toPlayer == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player not found"})
		return
	}

	fromPlayer.Points -= req.Amount
	toPlayer.Points += req.Amount

	ricchiWebsocket.GetHub().NotifyStateChange(table)

	if fromPlayer.Points < 0 {
		SettleTable(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "paid"})
}

type SettleRequest struct {
	TableID string `json:"table_id"`
}

type RankResult struct {
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	Points   int    `json:"points"`
	Rank     int    `json:"rank"`
}

func SettleTable(c *gin.Context) {
	var req SettleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	mgr := models.GetTableManager()
	table := mgr.GetTable(req.TableID)
	if table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	players := table.GetPlayers()
	type scored struct {
		player *models.Player
		points int
	}
	var ranked []scored
	for _, p := range players {
		ranked = append(ranked, scored{player: p, points: p.Points})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].points > ranked[j].points
	})

	var results []RankResult
	for i, s := range ranked {
		results = append(results, RankResult{
			OpenID:   s.player.OpenID,
			Nickname: s.player.Nickname,
			Points:   s.player.Points,
			Rank:     i + 1,
		})
	}

	table.Status = models.TableStatusSettled
	c.JSON(http.StatusOK, gin.H{"results": results})
}

type DismissRequest struct {
	TableID string `json:"table_id"`
}

func DismissTable(c *gin.Context) {
	var req DismissRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	mgr := models.GetTableManager()
	mgr.DeleteTable(req.TableID)
	c.JSON(http.StatusOK, gin.H{"status": "dismissed"})
}

type ReconnectRequest struct {
	TableID   string `json:"table_id"`
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func Reconnect(c *gin.Context) {
	var req ReconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	mgr := models.GetTableManager()
	table := mgr.GetTable(req.TableID)
	if table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	if player := table.GetPlayer(req.OpenID); player != nil {
		player.IsConnected = true
	} else {
		table.AddPlayer(req.OpenID, req.Nickname, req.AvatarURL, false)
	}

	ricchiWebsocket.GetHub().NotifyStateChange(table)
	c.JSON(http.StatusOK, gin.H{"status": "reconnected"})
}

func GetTable(c *gin.Context) {
	id := c.Param("id")
	mgr := models.GetTableManager()
	table := mgr.GetTable(id)
	if table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	players := table.GetPlayers()
	c.JSON(http.StatusOK, gin.H{
		"id":      table.ID,
		"status":  table.Status,
		"players": players,
	})
}