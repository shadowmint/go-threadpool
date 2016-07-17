package threadpool_test

import (
	"ntoolkit/assert"
	"ntoolkit/errors"
	"ntoolkit/threadpool"
	"testing"
)

func TestRun(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		pool := threadpool.New()
		value := 0

		T.Assert(pool.Run(func() { value += 1 }) == nil)

		pool.Wait()
		T.Assert(value == 1)
	})
}

func TestBusy(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		pool := threadpool.New()
		pool.MaxThreads = 2

		value := 0

		T.Assert(pool.Run(func() { value += 1 }) == nil)
		T.Assert(pool.Run(func() { value += 1 }) == nil)
		err := pool.Run(func() { value += 1 })

		T.Assert(err != nil)
		T.Assert(errors.Is(err, threadpool.ErrBusy{}))

		pool.Wait()
		T.Assert(value == 2)
	})
}

func TestWaitNext(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		pool := threadpool.New()
		pool.MaxThreads = 10

		value := 0
		procs := 0
		hits := 0

		for hits < 50 {
			pool.WaitNext()
			if err := pool.Run(func() { value += 1 }); err == nil {
				hits++
			}
			active := pool.Active()
			if active > procs {
				procs = active
			}
		}

		pool.Wait()
		T.Assert(procs == 10)
		T.Assert(value == 50)
	})
}
