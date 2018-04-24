package linkedlist

import (
	"testing"
	"math/rand"
	"fmt"
)

func TestLinkedList_ViewFromHead(t *testing.T) {
	N := 1000
	ll := NewLinkedList(100)
	for i := 0; i < N; i++ {
		ll.Add(rand.Float32()*1000, 1)
	}
	count := ll.ViewFromHead()
	fmt.Println(ll.Size(), count)
}

func TestLinkedList_ViewFromTail(t *testing.T) {
	N := 1000
	ll := NewLinkedList(100)
	for i := 0; i < N; i++ {
		ll.Add(rand.Float32()*1000, 1)
	}
	count := ll.ViewFromTail()
	fmt.Println(ll.Size(), count)
}
