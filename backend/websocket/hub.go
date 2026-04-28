package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"ricchi-room/backend/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	tables map[string]map[string]*websocket.Conn
	mu     sync.RWMutex
}

var (
	h  *Hub
	hOnce sync.Once
)

func GetHub() *Hub {
	hOnce.Do(func() {
		h = &Hub{
			tables: make(map[string]map[string]*websocket.Conn),
		}
	})
	return h
}

func (h *Hub) JoinTable(tableID, openID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.tables[tableID] == nil {
		h.tables[tableID] = make(map[string]*websocket.Conn)
	}
	h.tables[tableID][openID] = conn
}

func (h *Hub) LeaveTable(tableID, openID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.tables[tableID] != nil {
		delete(h.tables[tableID], openID)
	}
}

func (h *Hub) Broadcast(tableID string, msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if connections, ok := h.tables[tableID]; ok {
		data, _ := json.Marshal(msg)
		for _, conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("broadcast error: %v", err)
			}
		}
	}
}

func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request, tableID, openID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer func() {
		conn.Close()
		h.LeaveTable(tableID, openID)
	}()

	h.JoinTable(tableID, openID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *Hub) NotifyStateChange(table *models.Table) {
	msg := map[string]interface{}{
		"type":    "state_update",
		"table":   table,
	}
	h.Broadcast(table.ID, msg)
}