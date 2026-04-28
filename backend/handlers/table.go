package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ricchi-room/backend/config"
	"ricchi-room/backend/models"
	ricchiWebsocket "ricchi-room/backend/websocket"
)

type CreateTableRequest struct {
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
}

type CreateTableResponse struct {
	TableID string `json:"table_id"`
}

func CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	mgr := models.GetTableManager()
	if mgr.TableCount() >= config.Get().Server.MaxTables {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "max tables reached"})
		return
	}

	tableID, table := mgr.CreateTable()
	if tableID == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to create table"})
		return
	}

	table.AddPlayer(req.OpenID, req.Nickname, true)
	c.JSON(http.StatusOK, CreateTableResponse{TableID: tableID})
}

type JoinTableRequest struct {
	TableID  string `json:"table_id"`
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
}

func JoinTable(c *gin.Context) {
	var req JoinTableRequest
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

	if len(table.Players) >= config.Get().Game.MaxPlayers {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table is full"})
		return
	}

	table.AddPlayer(req.OpenID, req.Nickname, false)
	ricchiWebsocket.GetHub().NotifyStateChange(table)
	c.JSON(http.StatusOK, gin.H{"status": "joined"})
}

type LeaveTableRequest struct {
	TableID string `json:"table_id"`
	OpenID  string `json:"open_id"`
}

func LeaveTable(c *gin.Context) {
	var req LeaveTableRequest
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

	table.RemovePlayer(req.OpenID)
	ricchiWebsocket.GetHub().NotifyStateChange(table)

	if table.IsEmpty() {
		mgr.DeleteTable(req.TableID)
	}

	c.JSON(http.StatusOK, gin.H{"status": "left"})
}