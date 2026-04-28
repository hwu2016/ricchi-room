package models

import (
	"math/rand"
	"sync"
	"time"
)

type TableManager struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

var (
	manager  *TableManager
	onceMgr  sync.Once
)

func GetTableManager() *TableManager {
	onceMgr.Do(func() {
		manager = &TableManager{
			tables: make(map[string]*Table),
		}
	})
	return manager
}

func (m *TableManager) CreateTable() (string, *Table) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for attempts := 0; attempts < 100; attempts++ {
		id := generateTableID()
		if _, exists := m.tables[id]; !exists {
			table := NewTable(id)
			m.tables[id] = table
			return id, table
		}
	}
	return "", nil
}

func (m *TableManager) GetTable(id string) *Table {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tables[id]
}

func (m *TableManager) DeleteTable(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tables, id)
}

func (m *TableManager) TableCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tables)
}

func generateTableID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return string(rune('1'+r.Intn(9))) + string(rune('0'+r.Intn(10))) + string(rune('0'+r.Intn(10))) + string(rune('0'+r.Intn(10)))
}