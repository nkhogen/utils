package concurrent

import (
	"sync"
	"time"
)

// CondLock ...
type CondLock struct {
	l          *sync.Mutex
	il         *sync.Mutex // Internal lock
	ch         chan struct{}
	isChanOpen bool
}

// NewCondLock returns a conditional lock
func NewCondLock() *CondLock {
	return &CondLock{
		l:  &sync.Mutex{},
		il: &sync.Mutex{},
	}
}

// Lock ...
func (lock *CondLock) Lock() {
	lock.l.Lock()
}

// Unlock ...
func (lock *CondLock) Unlock() {
	lock.l.Unlock()
}

// Wait ... already in lock
func (lock *CondLock) Wait() {
	if !lock.isChanOpen {
		lock.ch = make(chan struct{})
		lock.isChanOpen = true
	}
	// This is required because lock.l.Unlock can be called when it is already being locked
	lock.il.Lock()
	defer lock.il.Unlock()
	// Release the lock
	lock.l.Unlock()
	<-lock.ch
	// Lock again before returning
	lock.l.Lock()
}

// TimedWait ...
func (lock *CondLock) TimedWait(d time.Duration) {
	if !lock.isChanOpen {
		lock.ch = make(chan struct{})
		lock.isChanOpen = true
	}
	lock.il.Lock()
	defer lock.il.Unlock()
	lock.l.Unlock()
	select {
	case <-time.After(d):
		break
	case <-lock.ch:
	}
	lock.l.Lock()
}

// Notify ...already locked
func (lock *CondLock) Notify() {
	if lock.isChanOpen {
		lock.ch <- struct{}{}
	}
}

// NotifyAll ...already locked
func (lock *CondLock) NotifyAll() {
	if lock.isChanOpen {
		// Closing the channel makes all the go-routines waiting
		// on the channel to return
		close(lock.ch)
		lock.isChanOpen = false
	}
}
