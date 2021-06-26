package sdb

import "errors"

type MemoryDb struct {
	data map[string]interface{}
}

func (m *MemoryDb) Get(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}

	return nil, errors.New("not found")
}


func (m *MemoryDb) Set(key string, value interface{}) error {
	m.data[key] = value
	return nil
}

func NewMemoryDb() *MemoryDb {
	return &MemoryDb{data: map[string]interface{}{}}
}