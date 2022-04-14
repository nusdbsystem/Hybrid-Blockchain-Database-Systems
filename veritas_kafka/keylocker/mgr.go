package keylocker

import (
	"fmt"
	"sync"
)

type KLocker interface {
	Lock(key string)
	Unlock(key string)
}

type KMutex struct {
	mu     sync.Mutex
	locks  map[string]*sync.Mutex
	counts map[string]int
}

func (km *KMutex) Lock(key string) {
	km.mu.Lock()
	if km.locks == nil {
		km.locks = make(map[string]*sync.Mutex)
		km.counts = make(map[string]int)
	}
	lock, ok := km.locks[key]
	if !ok {
		lock = &sync.Mutex{}
		km.locks[key] = lock
	}
	km.counts[key]++
	km.mu.Unlock()
	lock.Lock()
}

func (km *KMutex) Unlock(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()
	lock, ok := km.locks[key]
	if !ok || km.counts[key] == 0 {
		fmt.Printf("klocker: unlock unlocked kmutex of %#v\n", key)
		return
	}
	lock.Unlock()
	km.counts[key]--
	if km.counts[key] == 0 {
		delete(km.locks, key)
		delete(km.counts, key)
	}
}

type locker struct {
	kl  KLocker
	key string
}

func (l *locker) Lock()   { l.kl.Lock(l.key) }
func (l *locker) Unlock() { l.kl.Unlock(l.key) }

func (km *KMutex) Locker(key string) sync.Locker {
	return &locker{kl: km, key: key}
}
