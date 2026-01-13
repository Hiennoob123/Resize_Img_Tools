package main

import "fmt"

type Query int 
const (
	Enqueue Query = iota
	Dequeue
	None
)

// Request query for queue
type Request[T any] struct {
	query Query
	data T
}

type TopStruct[T any] struct {
	value T 
	err error
}

// Implementing safe multi-threads job queue
type Queue[T any] struct {
	items []T
	request chan Request[T]
	top chan TopStruct[T]
	stop bool 
	sz int
}

// Initialize newQueue with multi-threads Enqueue and Dequeue handler
func newQueue[T any](items []T) *Queue {
	queue := new(Queue[T])
	queue.items = items
	queue.request = make(chan Request)
	queue.top = make(chan TopStruct[T])
	queue.stop = true
	queue.sz = len(items)
	go func() {
		for !queue.stop {
			request := <- queue.request
			switch request.query {
			case None:
				break
			case Enqueue:
				queue.items = append(queue.items, request.data)
			case Dequeue:
				if len(queue.items) == 0 {
					queue.top <- TopStruct{
						value: new(T),
						err: nil,
					}
				}
				top_item := queue.items[0]
				if len(queue.items) == 1 {
					queue.items = nil
				} else {
					queue.items = queue.items[1:]
				}
				queue.top <- TopStruct{
					value: top_item,
					err: nil,
				}
			}
		}
	}()
	return queue
}

func (q *Queue[T]) stop_queue() {
	q.stop = true
	q.request <- Request {
		query: None,
		data: new(T),
	}
}


func (q *Queue[T]) enqueue(item T) {
	q.sz += 1
	q.request <- Request{
		query: Enqueue,
		data: item,
	}
}

func (q *Queue[T]) dequeue() T {
	q.sz -= 1
	q.request <- Request{
		query: Dequeue,
		data: new(T),
	}
	
	top := <- q.top

	if top.err != nil {
		fmt.Println("Queue Access Error:", top.err)
	}

	return top.value
}
