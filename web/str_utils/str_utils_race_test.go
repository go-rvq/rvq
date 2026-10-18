package str_utils

import (
	"sync"
	"testing"
)

// HumanizeString is called by many requests at once.
func TestHumanizeStringConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 200 {
				if got := HumanizeString("OrderItemID"); got != "Order Item ID" {
					t.Errorf("got %q", got)
					return
				}
			}
		}()
	}
	wg.Wait()
}
