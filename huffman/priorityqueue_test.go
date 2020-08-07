package huffman

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func validateHeapProperty(t *testing.T, q *priorityQueue, i int) {
	left := 2*i + 1
	right := 2*i + 2
	if left < q.Len() {
		if q.less(left, i) {
			t.Fatalf("heap property invalidated at index %d", i)
		}
		validateHeapProperty(t, q, left)
	}
	if right < q.Len() {
		if q.less(right, i) {
			t.Fatalf("heap property invalidated at index %d", i)
		}
		validateHeapProperty(t, q, right)
	}
}

func logPriorityQueue(t *testing.T, q *priorityQueue) {
	b := strings.Builder{}
	for _, item := range *q {
		b.WriteString(strconv.FormatInt(item.frequency, 10))
		b.WriteByte(' ')
	}
	t.Log(b.String())
}

func TestPriorityQueue(t *testing.T) {
	rand.Seed(9873214)
	newItem := func() *queueItem {
		return &queueItem{frequency: rand.Int63n(20)}
	}
	queue := priorityQueue{}
	for i := 0; i < 10; i++ {
		queue.Append(newItem())
	}
	t.Run("Init", func(t *testing.T) {
		logPriorityQueue(t, &queue)
		queue.Init()
		logPriorityQueue(t, &queue)
		validateHeapProperty(t, &queue, 0)
	})
	t.Run("Usage", func(t *testing.T) {
		queue.Push(newItem())
		validateHeapProperty(t, &queue, 0)
		queue.Push(newItem())
		validateHeapProperty(t, &queue, 0)
		queue.Push(newItem())
		validateHeapProperty(t, &queue, 0)

		check(t, int64(0), queue.Pop().frequency)
		validateHeapProperty(t, &queue, 0)
		queue.Push(newItem())
		validateHeapProperty(t, &queue, 0)
		check(t, int64(2), queue.Pop().frequency)
		validateHeapProperty(t, &queue, 0)
		check(t, int64(4), queue.Pop().frequency)
		validateHeapProperty(t, &queue, 0)
	})
}
