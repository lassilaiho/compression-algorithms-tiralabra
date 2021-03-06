package huffman

// queueItem is an item in a priorityQueue.
type queueItem struct {
	node      *codeTreeNode
	frequency int64
}

// priorityQueue is a priority queue of codeTreeNodes.
type priorityQueue []*queueItem

// Init establishes the heap property for q.
func (q *priorityQueue) Init() {
	for i := (len(*q) - 1) / 2; i >= 0; i-- {
		q.siftDown(i)
	}
}

// Push adds item to q.
func (q *priorityQueue) Push(item *queueItem) {
	node := len(*q)
	q.Append(item)
	parent := (node - 1) / 2
	for parent >= 0 && q.less(node, parent) {
		q.swap(node, parent)
		node, parent = parent, (parent-1)/2
	}
}

// Pop returns the item with the smallest frequency and removes it from the q.
func (q *priorityQueue) Pop() *queueItem {
	item := (*q)[0]
	(*q)[0] = (*q)[len(*q)-1]
	(*q) = (*q)[:len(*q)-1]
	q.siftDown(0)
	return item
}

// Append appends item to the end of q. The result may violate the heap
// property. It must be restored with a call to Init.
func (q *priorityQueue) Append(item *queueItem) {
	if len(*q) == cap(*q) {
		newCap := 2 * cap(*q)
		if newCap == 0 {
			newCap = 1
		}
		newQueue := make(priorityQueue, len(*q), newCap)
		for i := 0; i < len(*q); i++ {
			newQueue[i] = (*q)[i]
		}
		*q = newQueue
	}
	*q = (*q)[:len(*q)+1]
	(*q)[len(*q)-1] = item
}

// Len returns the number of items in q.
func (q *priorityQueue) Len() int {
	return len(*q)
}

// siftDown moves node downward in the heap until it satisfies the heap
// property.
func (q *priorityQueue) siftDown(node int) {
	for {
		min := node
		left := 2*node + 1
		right := 2*node + 2
		if left < len(*q) && q.less(left, min) {
			min = left
		}
		if right < len(*q) && q.less(right, min) {
			min = right
		}
		if min == node {
			return
		}
		q.swap(node, min)
		node = min
	}
}

// swap swaps the elements at positions i and j.
func (q *priorityQueue) swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

// less returns true if the frequency of the i'th element is less than that of
// the j'th element and false otherwise.
func (q *priorityQueue) less(i, j int) bool {
	return (*q)[i].frequency < (*q)[j].frequency
}
