package sort

import (
	"container/heap"
	"fmt"
)

// 普通的堆排序
type myHeap []int

func (h *myHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *myHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myHeap) Len() int {
	return len(*h)
}

// 把最后一个弹出，因为最小的值已经被换到了最后。
func (h *myHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myHeap) Push(v any) {
	*h = append(*h, v.(int))
}

// 优先队列
type Item struct {
	value    string
	priority int
	index    int
}
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func PriorityQueueTest() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}

type TimeSortedQueueItem struct {
	Time  int64
	Value interface{}
}

type TimeSortedQueue []*TimeSortedQueueItem

func (q TimeSortedQueue) Len() int           { return len(q) }
func (q TimeSortedQueue) Less(i, j int) bool { return q[i].Time < q[j].Time }
func (q TimeSortedQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *TimeSortedQueue) Push(v interface{}) {
	*q = append(*q, v.(*TimeSortedQueueItem))
}

func (q *TimeSortedQueue) Pop() interface{} {
	n := len(*q)
	item := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return item
}

func NewTimeSortedQueue(items ...*TimeSortedQueueItem) *TimeSortedQueue {
	q := make(TimeSortedQueue, len(items))
	for i, item := range items {
		q[i] = item
	}
	heap.Init(&q)
	return &q
}

func (q *TimeSortedQueue) PushItem(time int64, value interface{}) {
	heap.Push(q, &TimeSortedQueueItem{
		Time:  time,
		Value: value,
	})
}

func (q *TimeSortedQueue) PopItem() interface{} {
	if q.Len() == 0 {
		return nil
	}

	return heap.Pop(q).(*TimeSortedQueueItem).Value
}
