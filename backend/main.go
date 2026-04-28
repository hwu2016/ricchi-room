package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"ricchi-room/backend/config"
	"ricchi-room/backend/handlers"
	ricchiWebsocket "ricchi-room/backend/websocket"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	zapLog, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger = zapLog.Sugar()
	defer zapLog.Sync()

	cfg := config.Load()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/api/table/create", handlers.CreateTable)
	r.POST("/api/table/join", handlers.JoinTable)
	r.POST("/api/table/leave", handlers.LeaveTable)
	r.POST("/api/table/pay", handlers.Pay)
	r.POST("/api/table/settle", handlers.SettleTable)
	r.POST("/api/table/dismiss", handlers.DismissTable)
	r.POST("/api/table/reconnect", handlers.Reconnect)
	r.GET("/api/table/:id", handlers.GetTable)

	r.GET("/ws/table/:id/:openid", func(c *gin.Context) {
		tableID := c.Param("id")
		openID := c.Param("openid")
		ricchiWebsocket.GetHub().HandleConnection(c.Writer, c.Request, tableID, openID)
	})

	logger.Infow("server starting", "port", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		logger.Fatalw("server failed", "error", err)
	}
}