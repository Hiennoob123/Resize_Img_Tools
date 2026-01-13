package main
type Query int 
const (
	Enqueue Query = iota
	Dequeue
	None
)

// Request query for queue
type Request struct {
	query Query
	data interface{}
}

type TopStruct[T any] struct {
	value T 
	err error
}

// Implementing safe multi-threads job queue
type Queue[T any] struct {
	items []T
	request chan Request
	top chan TopStruct[T]
	stop bool 
	sz int
}

// Initialize newQueue with multi-threads Enqueue and Dequeue handler
func newQueue(items []T) *Queue {
	queue = new(queue[T])
	queue.items = items
	queue.Request = chan Request
	queue.top = chan TopStruct[T]
	queue.stop = true
	queue.sz = len(items)
	go func() {
		for !queue.stop {
			request := <- queue.request
			switch request.query {
			case None:
				break
			case Enqueue:
				queue = append(queue, request.data)
			case Dequeue:
				if q.IsEmpty() {
					queue.top <- TopStruct{
						value: new(T),
						err: nil,
					}
				}
				top_item := q.items[0]
				if len(q.items) == 1 {
					q.items = nil
				} else {
					q.items = q.items[1:]
				}
				queue.top <- TopStruct{
					value: top_item,
					err: nil,
				}
			}
		}
	}
}

func (q *Queue[T]) stop() {
	q.stop = true
	q.Request <- Request {
		query: None,
		data: 0,
	}
}


func (q *Queue[T]) enqueue(item T) {
	sz += 1
	q.request <- Request{
		query: Enqueue,
		data: item,
	}
}

func (q *Queue[T]) dequeue() T {
	sz -= 1
	q.request <- Request{
		query: Dequeue,
		data: 0,
	}
	
	top := <- q.top

	return top
}
