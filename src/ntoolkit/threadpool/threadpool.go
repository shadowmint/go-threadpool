package threadpool

import (
	"ntoolkit/errors"
	"sync"
)

// ThreadPool is a common high level interface for a single producer
// multi-consumer pattern.
type ThreadPool struct {
	MaxThreads int
	active     int
	lock       *sync.Mutex
	busy       *sync.Mutex
	any        *sync.Mutex
	queue      chan bool
}

// New returns a new empty ThreadPool
func New() *ThreadPool {
	return &ThreadPool{
		MaxThreads: -1,
		active:     0,
		lock:       &sync.Mutex{},
		busy:       &sync.Mutex{},
		any:        &sync.Mutex{},
		queue:      make(chan bool)}
}

// Run starts a new task, or raises an error.
// If the task fails, the error is silently consumed.
func (pool *ThreadPool) Run(task func()) error {
	var err error = nil
	pool.lock.Lock()
	if pool.activeUp() {
		go func() {
			defer func() {
				pool.lock.Lock()
				if r := recover(); r != nil {
				}
				pool.activeDown()
				pool.lock.Unlock()
			}()
			task()
		}()
	} else {
		err = errors.Fail(ErrBusy{}, nil, "No available threads for task")
	}
	pool.lock.Unlock()
	return err
}

// Active returns a count of active threads.
func (pool *ThreadPool) Active() int {
	return pool.active
}

// WaitNext blocks until a task slot is free and then returns.
// Notice that due to the async nature of the thread pool, immediately
// calling Run() after wait next is not guaranteed not to return an error.
// It just means that a task finished, and a slot opened; the next request
// to run will be serviced, regardless of where it comes from.
func (pool *ThreadPool) WaitNext() {
	pool.busy.Lock()
	pool.busy.Unlock()
}

// Wait blocks until all tasks are completed.
func (pool *ThreadPool) Wait() {
	pool.any.Lock()
	pool.any.Unlock()
}

// Update the active count and lock state
func (pool *ThreadPool) activeUp() bool {
	if pool.MaxThreads < 0 || pool.active < pool.MaxThreads {
		pool.active++
		if pool.active == 1 {
			pool.any.Lock()
		}
		if pool.active == pool.MaxThreads {
			pool.busy.Lock()
		}
		return true
	}
	return false
}

// Update the active count and lock state
func (pool *ThreadPool) activeDown() {
	if pool.active == pool.MaxThreads {
		pool.busy.Unlock()
	}
	pool.active--
	if pool.active == 0 {
		pool.any.Unlock()
	}
}
