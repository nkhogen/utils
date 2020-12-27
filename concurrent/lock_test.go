package concurrent

import (
	"fmt"
	"testing"
	"time"
	"sync"
)

// Rudimentary test
func TestNewCondLock(t *testing.T) {
	cond := NewCondLock()
	start := time.Now()
	fmt.Printf("Waiting for 5 secs in main goroutine...\n")
	cond.Lock()
	cond.TimedWait(time.Second * 5)
	fmt.Printf("Done waiting 5 secs in main goroutine... %+v\n", time.Since(start)/time.Second)
	go func() {
		t := time.Tick(time.Second * 3)
		t1 := time.Tick(time.Second * 10)
		for {
			select {
			case <-t:
				fmt.Printf("Notifying...\n")
				func() {
					cond.Lock()
					defer cond.Unlock()
					cond.Notify()
					fmt.Printf("Notified...\n")
				}()

			case <-t1:
				fmt.Printf("Notifying all...\n")
				func() {
					cond.Lock()
					defer cond.Unlock()
					cond.NotifyAll()
					fmt.Printf("Notified all...\n")
				}()
			}
		}
	}()
	fmt.Printf("Unlocking in main goroutine so that other goroutines can lock to wait")
	cond.Unlock()
	wg := &sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		fmt.Printf("Waiting again %d ...\n", i)
		go func(n int) {
			cond.Lock()
			defer cond.Unlock()
			cond.Wait()
			fmt.Printf("Done %d... %+v\n", n, time.Since(start)/time.Second)
			wg.Done()
		}(i)
	}
	fmt.Printf("Done ... %+v\n", time.Since(start)/time.Second)
	wg.Wait()
	fmt.Printf("Done all ... %+v\n", time.Since(start)/time.Second)
}
