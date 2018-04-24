package linkedlist

import (
	"fmt"
	"math/rand"
)

const PI = 1.5 * 3.14159265359

type Centroid struct {
	Mean   float32
	Weight float32
	Next   *Centroid
	Prev   *Centroid
}

func (c *Centroid) String() string {
	return fmt.Sprintf("Centroid{Mean=%f Weight=%f}", c.Mean, c.Weight)
}

func (c *Centroid) Update(that *Centroid) {
	c.Weight += that.Weight
	c.Mean += that.Weight * (that.Mean - c.Mean) / c.Weight
}

func (c *Centroid) Set(that *Centroid) {
	c.Mean = that.Mean
	c.Weight = that.Weight
}

type LinkedList struct {
	Compression float32
	Weights     float32
	Max         float32
	Fractor     float32
	Count       int
	Head        *Centroid
	Tail        *Centroid
}

func NewLinkedList(compression float32) *LinkedList {
	return &LinkedList{
		Compression: compression,
		Weights:     0,
		Max:         0,
		Fractor:     PI,
		Count:       0,
	}
}

func (l *LinkedList) Append(c *Centroid, that *Centroid) {
	that.Next = c.Next
	that.Prev = c

	if c != l.Tail {
		c.Next.Prev = that
	}
	c.Next = that

	if c == l.Tail {
		l.Tail = that
	}
}

func (l *LinkedList) Size() int {
	return l.Count
}

func (l *LinkedList) Add(x float32, w float32) {
	l.Weights += w

	c := &Centroid{Mean: x, Weight: w}
	if l.Count == 0 {
		l.Head = c
		l.Tail = c
		l.Count++
	} else {
		inserted := false
		if x < l.Head.Mean {
			c.Next = l.Head
			l.Head.Prev = c
			l.Head = c
			inserted = true
		} else if x > l.Tail.Mean {
			l.Tail.Next = c
			c.Prev = l.Tail
			l.Tail = c
			inserted = true
		}

		var wSoFar float32 = 0
		normalizer := l.Compression / (l.Fractor * l.Weights)

		count := 0
		node := l.Head
		next := node.Next
		for next != nil {
			if !inserted && next.Mean > x {
				l.Append(node, c)
				count++
				inserted = true
			}

			proposed := node.Weight + next.Weight
			z := proposed * normalizer
			q0 := wSoFar / l.Weights
			q2 := (wSoFar + proposed) / l.Weights
			if z*z <= q0*(1-q0) && z*z <= q2*(1-q2) {
				node.Update(next)
				node.Next = next.Next
				if node.Next != nil {
					node.Next.Prev = node
				}

				next = node.Next
			} else {
				wSoFar += node.Weight
				count++

				node = next
				next = next.Next
			}
		}
		l.Tail = node
		count++

		l.Count = count

		factor := l.Fractor * 1.5
		for float32(l.Count) > l.Compression+20 {
			l.compress(factor)
			factor = factor * 1.5
		}
	}

	if x > l.Max {
		l.Max = x
	}
}

func (l *LinkedList) compress(factor float32) {
	cs := make([]*Centroid, 0, l.Count)

	node := l.Head
	for node != nil {
		cs = append(cs, node)

		node = node.Next
	}

	for i := l.Count - 1; i > 1; i-- {
		j := rand.Intn(i - 1)
		cs[i], cs[j] = cs[j], cs[i]
	}

	l.Clear()
	l.Fractor = factor
	for _, c := range cs {
		l.Add(c.Mean, c.Weight)
	}
	l.Fractor = PI
}

func (l *LinkedList) Clear() {
	l.Weights = 0
	l.Max = 0
	l.Count = 0
	l.Head = nil
	l.Tail = nil
}

func (l *LinkedList) ViewFromHead() int {
	count := 0
	node := l.Head
	for node != nil {
		fmt.Println(node)

		node = node.Next
		count++
	}
	return count
}

func (l *LinkedList) ViewFromTail() int {
	count := 0
	node := l.Tail
	for node != nil {
		fmt.Println(node)

		node = node.Prev
		count++
	}
	return count
}
