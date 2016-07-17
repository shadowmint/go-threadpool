# go-threadpool

A simple threadpool implementation.

# Usage

    import "ntoolkit/threadpool"

    ...

    pool := threadpool.New()
    pool.MaxThreads = 10

    value := 0

    // Dispatch requests until all 50 have been serviced
    for hits < 50 {
      pool.WaitNext()
      if err := pool.Run(func() { value += 1 }); err == nil {
        hits++
      }
    }

    // Wait for any remaining requests to finish
    pool.Wait()
