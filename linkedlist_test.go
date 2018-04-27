package sorted_linklist

import (
	"testing"
	"math/rand"
	"fmt"
	"time"
)

func TestLinkedList_ViewFromHead(t *testing.T) {
	N := 1000
	ll := NewLinkedList(100)
	start := time.Now()

	var p float32 = 0
	for i := 0; i < N; i++ {
		p = ll.Add(rand.Float32()*1000, 1, 0.95)
	}
	fmt.Println("Cost", time.Now().Sub(start))
	fmt.Println("Percentile", ll.Percentile(0.95), p)
	fmt.Println("Size", ll.Size())
}
