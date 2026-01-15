package main

import "errors"

type Query int 
const (
	Enqueue Query = iota
	Dequeue
	Length
)

// Request query for queue
type Request[T any] struct {
	query Query
	data T
	top chan TopStruct[T]
}

type TopStruct[T any] struct {
	value T 
	length int
	err error
}

// Implementing safe multi-threads job queue
type Queue[T any] struct {
	items []T
	request chan *Request[T]
}

// Initialize newQueue with multi-threads Enqueue and Dequeue handler
func newQueue[T any](items []T) *Queue[T] {
	queue := new(Queue[T])
	queue.items = items
	queue.request = make(chan *Request[T])
	go func() {
		for {
			request, more := <- queue.request
			if !more {
				break
			}
			switch request.query {
			case Enqueue:
				queue.items = append(queue.items, request.data)
			case Dequeue:
				if len(queue.items) == 0 {
					request.top <- TopStruct[T]{
						value: *new(T),
						length: 0,
						err: errors.New("Queue is empty"),
					}
					continue
				}
				top_item := queue.items[0]
				if len(queue.items) == 1 {
					queue.items = nil
				} else {
					queue.items = queue.items[1:]
				}
				request.top <- TopStruct[T]{
					value: top_item,
					length: len(queue.items),
					err: nil,
				}
			case Length:
				request.top <- TopStruct[T]{
					value: *new(T),
					length: len(queue.items),
					err: nil,
				}
			}
		}
	}()
	return queue
}

func (q *Queue[T]) stop_queue() {
	close(q.request)
}


func (q *Queue[T]) enqueue(item T) {
	rep := make(chan TopStruct[T])
	q.request <- &Request[T]{
		query: Enqueue,
		data: item,
		top: rep,
	}
}

func (q *Queue[T]) dequeue() (T, error) {
	rep := make(chan TopStruct[T])
	q.request <- &Request[T]{
		query: Dequeue,
		data: *new(T),
		top: rep, 
	}
	
	top := <- rep

	return top.value, top.err
}

func (q *Queue[T]) length() int {
	rep := make(chan TopStruct[T])
	q.request <- &Request[T]{
		query: Length,
		data: *new(T),
		top: rep,
	}
	sz := <- rep
	return sz.length
}

