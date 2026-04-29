package Queue

import (
	"fmt"
	"sync"
)

type Node struct {
	Val  string
	Next *Node
}

type Queue struct {
	Name       string
	head       *Node
	tail       *Node
	NotifyChan chan struct{}
	mu         sync.Mutex
}

func InitNewQ(QName string) *Queue {
	return &Queue{Name: QName, head: nil, tail: nil, NotifyChan: make(chan struct{}, 1)}
}

func (q *Queue) AddQ(val string) {
	q.mu.Lock()
	newNode := &Node{Val: val}

	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.Next = newNode
		q.tail = newNode
	}
	q.mu.Unlock()
	select {
	case q.NotifyChan <- struct{}{}:
	default:

	}
}

func (q *Queue) GetQ() *Node {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == nil {
		return nil
	}
	retQ := q.head

	q.head = q.head.Next

	if q.head == nil {
		q.tail = nil
	}
	return retQ
}

func (q *Queue) PrintQ() {
	curr := q.head

	for curr != nil {
		fmt.Printf("%s -> ", curr.Val)
		curr = curr.Next
	}
	fmt.Println("nil")
}
