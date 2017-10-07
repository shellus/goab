package goab

import (
	"sync"
	"time"
	"sort"
)

type process struct {
	arr []time.Duration
	mux sync.Mutex
}

func NewProcess() *process {
	return &process{
		arr: make([]time.Duration, 0),
		mux: sync.Mutex{},
	}
}
func (t *process) Len() int {
	return len(t.arr)
}
func (t *process) Less(i, j int) bool {
	return int64(t.arr[i]) < int64(t.arr[j])
}
func (t *process) Swap(i, j int) {
	t.arr[i], t.arr[j] = t.arr[j], t.arr[i]
}

func (t *process) Add(processTime time.Duration) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.arr = append(t.arr, processTime)
}

func (t *process) Dump() (resultSlice []time.Duration) {
	t.mux.Lock()
	defer t.mux.Unlock()

	sort.Sort(t)
	const slice = 10
	var sliceLen int = int(len(t.arr) / slice)

	for i := 1; i <= slice; i++ {
		var sliceTotal time.Duration
		var sliceMaxIndex int
		if i == slice {
			// 最后一次就统计所有的
			sliceMaxIndex = len(t.arr)
		} else {
			sliceMaxIndex = i * sliceLen
		}
		for j := 0; j < sliceMaxIndex; j++ {
			sliceTotal = sliceTotal + t.arr[j]
		}
		resultSlice = append(resultSlice, sliceTotal/time.Duration(slice*sliceMaxIndex))
	}
	return
}

type SafeCounter struct {
	incs map[string]int
	mux  sync.Mutex
}

func NewCounter() *SafeCounter {
	return &SafeCounter{
		incs: make(map[string]int),
		mux:  sync.Mutex{},
	}
}

func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.incs[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.incs[key]
}

func (c *SafeCounter) Dump() map[string]int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.incs
}
