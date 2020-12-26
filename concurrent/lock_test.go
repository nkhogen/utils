package concurrent

import (
	"fmt"
	"testing"
	"time"
)

// Rudimentary test
func TestNewCondLock(t *testing.T) {
	cond := NewCondLock()
	start := time.Now()
	fmt.Printf("Waiting ...\n")
	cond.Lock()
	cond.TimedWait(time.Second * 5)
	fmt.Printf("Done ... %+v\n", time.Since(start)/time.Second)
	go func() {
		t := time.Tick(time.Second * 5)
		t1 := time.Tick(time.Second * 10)
		for {
			select {
			case <-t:
				fmt.Printf("Notifying...\n")
				func() {
					cond.Lock()
					defer cond.Unlock()
					cond.Notify()
				}()

			case <-t1:
				fmt.Printf("Notifying all...\n")
				func() {
					cond.Lock()
					defer cond.Unlock()
					cond.NotifyAll()
				}()
			}
		}
	}()

	for i := 0; i < 5; i++ {
		fmt.Printf("Waiting again %d ...\n", i)
		go func() {
			cond.Lock()
			defer cond.Unlock()
			cond.Wait()
		}()
	}
	cond.Wait()
	fmt.Printf("Done ... %+v\n", time.Since(start)/time.Second)
	fmt.Printf("Waiting again ...\n")
	cond.TimedWait(time.Second * 5)
	fmt.Printf("Done ... %+v\n", time.Since(start)/time.Second)
	cond.Unlock()
}
