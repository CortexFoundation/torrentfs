package shard

import (
	"runtime"
	"sync"

	"github.com/cespare/xxhash"
	"github.com/tidwall/hashmap"
)

// Map is a hashmap. Like map[string]any, but sharded and thread-safe.
type Map struct {
	once    sync.Once
	capcity int
	shards  int
	//seed    uint32
	mus  []sync.RWMutex
	maps []*hashmap.Map[string, any]
}

// New returns a new hashmap with the specified capacity. This function is only
// needed when you must define a minimum capacity, otherwise just use:
//
//	var m shard.Map
func New(capcity int) *Map {
	return &Map{capcity: capcity}
}

// Clear out all values from map
func (m *Map) Clear() {
	m.initDo()
	for i := 0; i < m.shards; i++ {
		m.mus[i].Lock()
		m.maps[i] = hashmap.New[string, any](m.capcity / m.shards)
		m.mus[i].Unlock()
	}
}

// Set assigns a value to a key.
// Returns the previous value, or false when no value was assigned.
func (m *Map) Set(key string, value any) (prev any, replaced bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	prev, replaced = m.maps[shard].Set(key, value)
	m.mus[shard].Unlock()
	return
}

// SetAccept assigns a value to a key. The "accept" function can be used to
// inspect the previous value, if any, and accept or reject the change.
// It's also provides a safe way to block other others from writing to the
// same shard while inspecting.
// Returns the previous value, or false when no value was assigned.
func (m *Map) SetAccept(
	key string, value any,
	accept func(prev any, replaced bool) bool,
) (prev any, replaced bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	defer m.mus[shard].Unlock()
	prev, replaced = m.maps[shard].Set(key, value)
	if accept != nil {
		if !accept(prev, replaced) {
			// revert unaccepted change
			if !replaced {
				// delete the newly set data
				m.maps[shard].Delete(key)
			} else {
				// reset updated data
				m.maps[shard].Set(key, prev)
			}
			prev, replaced = nil, false
		}
	}
	return
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map) Get(key string) (value any, ok bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].RLock()
	value, ok = m.maps[shard].Get(key)
	m.mus[shard].RUnlock()
	return
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) Delete(key string) (prev any, deleted bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	prev, deleted = m.maps[shard].Delete(key)
	m.mus[shard].Unlock()
	return
}

// DeleteAccept deletes a value for a key. The "accept" function can be used to
// inspect the previous value, if any, and accept or reject the change.
// It's also provides a safe way to block other others from writing to the
// same shard while inspecting.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) DeleteAccept(
	key string,
	accept func(prev any, replaced bool) bool,
) (prev any, deleted bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	defer m.mus[shard].Unlock()
	prev, deleted = m.maps[shard].Delete(key)
	if accept != nil {
		if !accept(prev, deleted) {
			// revert unaccepted change
			if deleted {
				// reset updated data
				m.maps[shard].Set(key, prev)
			}
			prev, deleted = nil, false
		}
	}

	return
}

// Len returns the number of values in map.
func (m *Map) Len() (length int) {
	m.initDo()
	for i := 0; i < m.shards; i++ {
		m.mus[i].RLock()
		length += m.maps[i].Len()
		m.mus[i].RUnlock()
	}
	return
}

// Range iterates overall all key/values.
// It's not safe to call or Set or Delete while ranging.
func (m *Map) Range(iter func(key string, value any) bool) {
	m.initDo()
	var done bool
	for i := 0; i < m.shards; i++ {
		func() {
			m.mus[i].RLock()
			defer m.mus[i].RUnlock()
			m.maps[i].Scan(func(key string, value any) bool {
				if !iter(key, value) {
					done = true
					return false
				}
				return true
			})
		}()
		if done {
			break
		}
	}
}

func (m *Map) choose(key string) int {
	return int(xxhash.Sum64String(key) & uint64(m.shards-1))
}

func (m *Map) initDo() {
	m.once.Do(func() {
		m.shards = 1
		for m.shards < runtime.NumCPU()*16 {
			m.shards *= 2
		}
		scap := m.capcity / m.shards
		m.mus = make([]sync.RWMutex, m.shards)
		m.maps = make([]*hashmap.Map[string, any], m.shards)
		for i := 0; i < len(m.maps); i++ {
			m.maps[i] = hashmap.New[string, any](scap)
		}
	})
}
