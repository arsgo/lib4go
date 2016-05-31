package pool

import (
	"sync/atomic"
	"testing"
)

type benchClient struct {
}

func (n *benchClient) Close() {
}

type bentchFactory struct {
}

func (n *bentchFactory) Create() (Object, error) {
	return &benchClient{}, nil
}
func (n *bentchFactory) Close() {

}
func BenchmarkAll(t *testing.B) {
	max := 1000
	min := 100
	var index int32
	var concurrent int32 = 1000000
	groupName := "test"
	benchPool := New()
	benchPool.Register(groupName, &bentchFactory{}, min, max)
	ch := make(chan int, max)
	close := make(chan int, 1)

	for i := 0; i < min*10; i++ {
		ch <- i
		go func() {
			for {
				if atomic.LoadInt32(&index) >= concurrent {
					close <- 1
					break
				}
				<-ch
				obj, err := benchPool.Get(groupName)
				if err != nil {
					t.Error(err)
				} else {
					benchPool.Recycle(groupName, obj)
				}
				atomic.AddInt32(&index, 1)
				ch <- 1
			}

		}()
	}
	<-close

}
