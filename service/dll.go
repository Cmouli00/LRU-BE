package service

import (
	"time"
)

type Node struct {
	key        string
	value      interface{}
	expiration time.Time
	prev       *Node
	next       *Node
}

type DoublyLinkedList interface {
	MoveToFront(node *Node)
	PushFront(node *Node)
	Remove(node *Node)
	RemoveLast() *Node
}

type Dll struct {
	head *Node
	tail *Node
}

func NewDll() *Dll {
	return &Dll{}
}

func (dll *Dll) MoveToFront(node *Node) {
	if dll.head == node {
		return
	}
	dll.Remove(node)
	dll.PushFront(node)
}

func (dll *Dll) PushFront(node *Node) {
	node.next = dll.head
	node.prev = nil
	if dll.head != nil {
		dll.head.prev = node
	}
	dll.head = node
	if dll.tail == nil {
		dll.tail = node
	}
}

func (dll *Dll) Remove(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
	}
}

func (dll *Dll) RemoveLast() *Node {
	if dll.tail == nil {
		return nil
	}
	node := dll.tail
	dll.Remove(node)
	return node
}
