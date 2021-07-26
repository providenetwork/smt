package smt

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// MapStore is a key-value store.
type MapStore interface {
	Get(key []byte) ([]byte, error)     // Get gets the value for a key.
	Set(key []byte, value []byte) error // Set updates the value for a key.
	Delete(key []byte) error            // Delete deletes a key.
}

// InvalidKeyError is thrown when a key that does not exist is being accessed.
type InvalidKeyError struct {
	Key []byte
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("invalid key: %x", e.Key)
}

// SimpleMap is a simple in-memory map.
type SimpleMap struct {
	m map[string][]byte
}

// NewSimpleMap creates a new empty SimpleMap.
func NewSimpleMap() *SimpleMap {
	return &SimpleMap{
		m: make(map[string][]byte),
	}
}

// SimpleMapFactory initializes a SimpleMap with m.
func SimpleMapFactory(m map[string][]byte) *SimpleMap {
	return &SimpleMap{
		m: m,
	}
}

// Get gets the value for a key.
func (sm *SimpleMap) Get(key []byte) ([]byte, error) {
	if value, ok := sm.m[string(key)]; ok {
		return value, nil
	}
	return nil, &InvalidKeyError{Key: key}
}

// Set updates the value for a key.
func (sm *SimpleMap) Set(key []byte, value []byte) error {
	sm.m[string(key)] = value
	return nil
}

// Delete deletes a key.
func (sm *SimpleMap) Delete(key []byte) error {
	_, ok := sm.m[string(key)]
	if ok {
		delete(sm.m, string(key))
		return nil
	}
	return &InvalidKeyError{Key: key}
}

// MarshalJSON marshal to JSON
func (sm *SimpleMap) MarshalJSON() ([]byte, error) {
	val := map[string]interface{}{}
	for k := range sm.m {
		val[hex.EncodeToString([]byte(k))] = hex.EncodeToString(sm.m[k])
	}
	return json.Marshal(val)
}

// UnmarshalJSON unmarshal from JSON
func (sm *SimpleMap) UnmarshalJSON(raw []byte) error {
	var val map[string]string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return err
	}

	sm.m = map[string][]byte{}
	for k := range val {
		_key, err := hex.DecodeString(k)
		if err != nil {
			return err
		}

		_val, err := hex.DecodeString(val[k])
		if err != nil {
			return err
		}

		sm.m[string(_key)] = _val
	}

	return json.Unmarshal(raw, &sm.m)
}
