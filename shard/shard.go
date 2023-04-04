package shard

import (
	"runtime"
	"sync"

	"github.com/cespare/xxhash"
	"github.com/tidwall/hashmap"
)

// Map is a hashmap. Like map[string]any, but sharded and thread-safe.
type Map[V any] struct {
	once    sync.Once
	capcity int
	shards  int
	//seed    uint32
	mus  []sync.RWMutex
	maps []*hashmap.Map[string, V]
}

// New returns a new hashmap with the specified capacity.
func New[V any](capcity int) (m *Map[V]) {
	m = &Map[V]{capcity: capcity}
	m.initDo()
	return
}

// Clear out all values from map
func (m *Map[V]) Clear() {
	for i := 0; i < m.shards; i++ {
		m.mus[i].Lock()
		m.maps[i] = hashmap.New[string, V](m.capcity / m.shards)
		m.mus[i].Unlock()
	}
}

// Set assigns a value to a key.
// Returns the previous value, or false when no value was assigned.
func (m *Map[V]) Set(key string, value V) (prev any, replaced bool) {
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
func (m *Map[V]) SetAccept(
	key string, value V,
	accept func(prev V, replaced bool) bool,
) (prev V, replaced bool) {
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
			//prev = nil
			replaced = false
		}
	}
	return
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map[V]) Get(key string) (value V, ok bool) {
	shard := m.choose(key)
	m.mus[shard].RLock()
	value, ok = m.maps[shard].Get(key)
	m.mus[shard].RUnlock()
	return
}

// Looks up an item under specified key
func (m *Map[V]) Has(key string) (ok bool) {
	shard := m.choose(key)
	m.mus[shard].RLock()
	_, ok = m.maps[shard].Get(key)
	m.mus[shard].RUnlock()
	return
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map[V]) Delete(key string) (prev V, deleted bool) {
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
func (m *Map[V]) DeleteAccept(
	key string,
	accept func(prev V, replaced bool) bool,
) (prev V, deleted bool) {
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
			//prev, deleted = nil, false
			deleted = false
		}
	}

	return
}

// Len returns the number of values in map.
func (m *Map[V]) Len() (length int) {
	for i := 0; i < m.shards; i++ {
		m.mus[i].RLock()
		length += m.maps[i].Len()
		m.mus[i].RUnlock()
	}
	return
}

// Range iterates overall all key/values.
// It's not safe to call or Set or Delete while ranging.
func (m *Map[V]) Range(iter func(key string, value V) bool) {
	var done bool
	for i := 0; i < m.shards; i++ {
		func() {
			m.mus[i].RLock()
			defer m.mus[i].RUnlock()
			m.maps[i].Scan(func(key string, value V) bool {
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

func (m *Map[V]) choose(key string) int {
	return int(xxhash.Sum64String(key) & uint64(m.shards-1))
}

func (m *Map[V]) initDo() {
	m.once.Do(func() {
		m.shards = 1
		for m.shards < runtime.NumCPU()*16 {
			m.shards *= 2
		}
		scap := m.capcity / m.shards
		m.mus = make([]sync.RWMutex, m.shards)
		m.maps = make([]*hashmap.Map[string, V], m.shards)
		for i := 0; i < len(m.maps); i++ {
			m.maps[i] = hashmap.New[string, V](scap)
		}
	})
}
