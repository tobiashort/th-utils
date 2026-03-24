package main

import (
	"math/rand"
	"time"

	"github.com/tobiashort/th-utils/lib/worker"
)

func main() {
	pool := worker.NewPool(10)

	for i := range 20 {
		worker := pool.GetWorker()
		worker.Go(
			func() {
				worker.Printf("%d: started", i)
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				worker.Printf("%d: processing", i)
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				worker.Logf("%d: done", i)
			})
	}

	pool.Wait()
}
