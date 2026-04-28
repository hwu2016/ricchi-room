package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type TableStatus string

const (
	TableStatusWaiting  TableStatus = "waiting"
	TableStatusPlaying  TableStatus = "playing"
	TableStatusSettled TableStatus = "settled"
)

type Table struct {
	ID         string            `json:"id"`
	Players   map[string]*Player `json:"players"`
	Status    TableStatus       `json:"status"`
	mu        sync.RWMutex
	Conns     map[string]*websocket.Conn
	CreatedAt time.Time `json:"created_at"`
}

type Player struct {
	OpenID      string `json:"open_id"`
	Nickname    string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	Points     int    `json:"points"`
	IsOwner    bool   `json:"is_owner"`
	IsConnected bool   `json:"is_connected"`
	SeatIndex  int    `json:"seat_index"`
}

func NewTable(id string) *Table {
	return &Table{
		ID:        id,
		Players:  make(map[string]*Player),
		Status:   TableStatusWaiting,
		Conns:    make(map[string]*websocket.Conn),
		CreatedAt: time.Now(),
	}
}

func (t *Table) AddPlayer(openID, nickname, avatarURL string, isOwner bool) (*Player, int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seat := len(t.Players)
	player := &Player{
		OpenID:       openID,
		Nickname:     nickname,
		AvatarURL:    avatarURL,
		Points:      25000,
		IsOwner:      isOwner,
		IsConnected: true,
		SeatIndex:   seat,
	}
	t.Players[openID] = player
	return player, seat
}

func (t *Table) RemovePlayer(openID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Players, openID)
}

func (t *Table) GetPlayer(openID string) *Player {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Players[openID]
}

func (t *Table) SetPoints(openID string, points int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if p, ok := t.Players[openID]; ok {
		p.Points = points
	}
}

func (t *Table) GetPlayers() []*Player {
	t.mu.RLock()
	defer t.mu.RUnlock()
	players := make([]*Player, 0, len(t.Players))
	for _, p := range t.Players {
		players = append(players, p)
	}
	return players
}

func (t *Table) IsEmpty() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.Players) == 0
}

func (t *Table) ConnectedCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	count := 0
	for _, p := range t.Players {
		if p.IsConnected {
			count++
		}
	}
	return count
}