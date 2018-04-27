package sorted_linklist

import (
	"fmt"
	"log"
)

const PI = 3.14159265359

type Centroid struct {
	Mean   float32
	Weight float32
	Next   *Centroid
}

func (c *Centroid) String() string {
	if c.Next != nil {
		return fmt.Sprintf("Centroid{Mean=%f Weight=%f Next}", c.Mean, c.Weight)
	}
	return fmt.Sprintf("Centroid{Mean=%f Weight=%f nil}", c.Mean, c.Weight)
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
	Count       int
	Head        *Centroid
}

func NewLinkedList(compression float32) *LinkedList {
	return &LinkedList{
		Compression: compression,
		Weights:     0,
		Max:         0,
		Count:       0,
	}
}

func (l *LinkedList) Append(c *Centroid, that *Centroid) {
	that.Next = c.Next
	c.Next = that
}

func (l *LinkedList) Size() int {
	return l.Count
}

func (l *LinkedList) Add(x float32, w float32, q float32) float32 {
	l.Weights += w

	index := q * l.Weights
	var p float32 = -1
	var node *Centroid

	c := &Centroid{Mean: x, Weight: w}
	if l.Count == 0 {
		l.Head = c
		l.Count++

		p = l.Head.Mean
	} else {
		inserted := false
		if x < l.Head.Mean {
			c.Next = l.Head
			l.Head = c
			inserted = true
		}

		var wSoFar float32 = 0

		wDw := l.Head.Weight / 2.0
		normalizer := l.Compression / (PI * l.Weights)

		l.Count = 0
		node = l.Head
		next := node.Next
		for next != nil {
			if !inserted && next.Mean > x {
				l.Append(node, c)
				l.Count++
				inserted = true
			}

			proposed := node.Weight + next.Weight

			// Percentile
			dw := proposed / 2.0
			if wDw+dw > index {
				z1 := index - wDw
				z2 := wDw + dw - index
				p = (node.Mean*z2 + next.Mean*z1) / (z1 + z2)
			}
			wDw += dw

			z := proposed * normalizer
			q0 := wSoFar / l.Weights
			q2 := (wSoFar + proposed) / l.Weights
			if z*z <= q0*(1-q0) && z*z <= q2*(1-q2) {
				node.Update(next)
				node.Next = next.Next

				next = node.Next
			} else {
				wSoFar += node.Weight
				l.Count++

				node = next
				next = next.Next
			}
		}
		if !inserted {
			node.Next = c
			node = node.Next
			l.Count++
		}
		l.Count++

		if float32(l.Count) > 2*l.Compression+20 {
			log.Fatal("Overflow ", l.Count)
		}
	}

	if x > l.Max {
		l.Max = x
	}

	if p == -1 && node != nil {
		z1 := index - l.Weights - node.Weight/2.0
		z2 := node.Weight/2.0 - z1
		p = (node.Mean*z1 + l.Max*z2) / (z1 + z2)
	}

	if p == -1 {
		p = 0
	}

	return p
}

func (l *LinkedList) Percentile(q float32) float32 {
	index := q * l.Weights

	var node *Centroid

	wSoFar := l.Head.Weight / 2.0
	for node = l.Head; node.Next != nil; node = node.Next {
		dw := (node.Weight + node.Next.Weight) / 2.0
		if wSoFar+dw > index {
			z1 := index - wSoFar
			z2 := wSoFar + dw - index
			return (node.Mean*z2 + node.Next.Mean*z1) / (z1 + z2)
		}
		wSoFar += dw
	}

	z1 := index - l.Weights - node.Weight/2.0
	z2 := node.Weight/2.0 - z1
	return (node.Mean*z1 + l.Max*z2) / (z1 + z2)
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
