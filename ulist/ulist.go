// Copyright 2020  himei@tuta.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ulist implements unrolled linked list.
//
// 	Each node holds up to a certain maximum number of elements, typically just
// 	large enough so that the node fills a single cache line or a small multiple
// 	thereof.
//
// 	To insert a new element, we simply find the node the element should
// 	be in and insert the element into the elements array, incrementing
// 	numElements. If the array is already full, we first insert a new node
// 	either preceding or following the current one and move half of the
// 	elements in the current node into it.
//
// 	To remove an element, we simply find the node it is in and delete it
// 	from the elements array, decrementing numElements. If this reduces
// 	the node to less than half-full, then we move elements from the next node
// 	to fill it back up above half. If this leaves the next node less
// 	than half full, then we move all its remaining elements into the
// 	current node, then bypass and delete it.
//
// See http://en.wikipedia.org/wiki/Unrolled_linked_list for details
package ulist

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/cpu"
)

// CacheLineSize represents the cache line size of the current CPU.
// See https://en.wikipedia.org/wiki/CPU_cache for details.
const CacheLineSize = int(unsafe.Sizeof(cpu.CacheLinePad{}))

// ulistNode is a single node of the unrolled linked list.
// It contains links to previous and next node, number of stored elements and
// slice of elements.
type ulistNode struct {
	next     *ulistNode
	prev     *ulistNode
	size     int // number of elements
	capacity int // max number of elements
	elems    []interface{}
}

// newUlistNode creates empty instance of list's node.
// All elements in emtpy node is set to nil.
func newUlistNode(c int) *ulistNode {
	elems := make([]interface{}, 0)

	// fill elems field with nil elemets
	for i := 0; i < c; i++ {
		elems = append(elems, nil)
	}

	return &ulistNode{
		next:     nil,
		prev:     nil,
		size:     0,
		elems:    elems,
		capacity: c,
	}
}

// add sets the first non-nil element equal to the given value
// and increments size of node. If the node is full, this function creates
// a new node and moves to it a number of elements equal to half the
// length of the cuttent node. in this case, the new element is
// added to the end of the new node. The function returns a new node,
// empty if no elements were moved.
func (un *ulistNode) add(val interface{}) *ulistNode {
	newNode := newUlistNode()

	if !un.isFull() {
		for i := range un.elems {
			if un.elems[i] == nil {
				un.elems[i] = val
				break
			}
		}

		un.size++
	} else {
		// elements to move
		tmv := un.capacity / 2
		// element to start moving
		start := un.capacity - tmv

		for i := 0; i < tmv; i++ {
			newNode.elems[i] = un.elems[start+i]
			newNode.size++
			un.elems[start+i] = nil
			un.size--
		}

		newNode.elems[newNode.size] = val
		newNode.size++
	}

	return newNode
}

// del removes the element with the given index from the node.
// If this reduces the node to less than half-full, then it moves
// elements from the next node (if that not nil) to fill node back up
// above half. If this leaves the next node less than half full, then it move all
// next node's remaining elements into the current node, then delete it.
// It returns zero if next node was not deleted and 1 in other case. If given
// index is greater than node's capacity, it returns error.
func (un *ulistNode) del(index int) (int, error) {
	var (
		err error
		n   = 0
	)

	if index > un.capacity {
		err = errors.New("Element index is out of range")
		return n, err
	}

	if un.size > un.capacity/2 {
		un.elems[index] = nil
		un.size--
	} else {
		if un.next != nil {
			un.elems[index] = nil
			un.size--

			un.elems[index] = un.next.elems[un.next.size-1]
			un.size++
			un.next.elems[un.next.size-1] = nil
			un.next.size--

			if un.next.size < un.capacity/2 {
				for i := range un.next.elems {
					if un.next.elems[i] != nil {
						un.elems[un.size+i] = un.next.elems[i]
						un.size++
					} else {
						break
					}
				}

				// if next node exists
				if un.next.next != nil {
					un.next = un.next.next
					un.next.next.prev = un
				}

				// indicate that the next node has been removed
				n++
			}
		}
	}

	return n, err
}

func (un *ulistNode) do(fn func(interface{})) {
	for i := range un.elems {
		if un.elems[i] != nil {
			fn(un.elems[i])
		} else {
			break
		}
	}
}

// isFull checks if node is full
func (un *ulistNode) isFull() bool {
	return un.size == un.capacity
}

// Ulist is an unrolled linked list itself.
// It contains links to first and last nodes and number of nodes.
type Ulist struct {
	first *ulistNode
	last  *ulistNode
	size  int // number of nodes
}

// NewUlist creates new empty unrolled linked list. It has only one (empty)
// node wich is first and last same time. Returns pointer to empty list.
func newUlist(c int) *Ulist {
	var (
		ul   = &Ulist{}
		node = newUlistNode(c)
	)

	ul.first = node
	ul.last = ul.first

	ul.first.next = ul.last
	ul.last.prev = ul.first

	ul.size = 1

	return ul
}

// NewUlist creates new empty unrolled linked list. It has only one (empty)
// node wich is first and last same time. Elem fields of list's nodes
// have length equal to CacheLineSize. Returns pointer to empty list.
func NewUlist() *Ulist {
	return newUlist(CacheLineSize)
}

// NewUlistCustomCap creates new empty unrolled linked list.
// Elem fields of list's nodes have length equal c. Returns pointer to empty list.
func NewUlistCustomCap(c int) *Ulist {
	return newUlist(c)
}

// GetSize returns number of list's nodes
func (ul *Ulist) GetSize() int {
	return ul.size
}

// findNode finds node with given index num. If num is greater than half-size of
// list, search starts from first node. Else search starts from last node.
// If num is greater then node size, it returns error.
func (ul *Ulist) findNode(num int) (*ulistNode, error) {
	var (
		err     error
		newNode = &ulistNode{}
	)

	if num > ul.GetSize() {
		err = errors.New("Node index is out of range")
		return newNode, err
	}

	count := 0

	// start from front
	if (ul.size - num) >= (ul.size / 2) {
		newNode.next = ul.first.next

		for count != num {
			newNode = newNode.next
			count++
		}
	} else { // start from back
		count = ul.size
		newNode.prev = ul.last.prev

		for count != num {
			newNode = newNode.prev
			count--
		}
	}

	return newNode, err
}

// Push appends new element val to the end of list.
func (ul *Ulist) Push(val interface{}) {
	newNode := ul.last.add(val)

	if newNode.size != 0 {
		// link new node and list's last node
		ul.last.next = newNode
		newNode.prev = ul.last

		// set new node to the list's last node
		ul.last = newNode

		// increment list's size
		ul.size++
	}
}

// Insert inserts a new element val at the target node with index num.
// If target node is full, it creates a new node and moves there the number
// of elements of the target node equal to half the length of the node.
// New element val will be added to the end of new node. New node
// will be inserted to list after target node. Function returns error if given index
// is greater than node.capacity.
func (ul *Ulist) Insert(val interface{}, num int) error {
	var (
		targetNode = &ulistNode{}
		err        error
	)

	targetNode, err = ul.findNode(num)

	if err != nil {
		return err
	}

	newNode := targetNode.add(val)

	if newNode.size != 0 {
		targetNode.next.prev = newNode
		newNode.next = targetNode.next
		newNode.prev = targetNode
		ul.size++
	}

	return err
}

// RemoveInNode removes element at index elemNum from node at index nodeNum.
func (ul *Ulist) RemoveFromNode(nodeNum, elemNum int) {
	var (
		err  error
		n    int
		node = &ulistNode{}
	)

	node, err = ul.findNode(nodeNum)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	n, err = node.del(elemNum)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if n != 0 {
		ul.size -= n
	}
}

// Do calls function fn on each list's element.
func (ul *Ulist) Do(fn func(interface{})) {
	var (
		newNode = newUlistNode()
		count   = 0
	)

	newNode = ul.first

	for count < ul.GetSize() {
		newNode.do(fn)
		newNode = newNode.next
		count++
	}
}

// Print prints each list's element.
func (ul *Ulist) Print() {
	fn := func(i interface{}) {
		fmt.Printf("%v\n", i)
	}

	ul.Do(fn)
}

// Clear removes all elements from list.
func (ul *Ulist) Clear() int {
	fn := func(i interface{}) {
		i = nil
	}

	ul.Do(fn)

	return ul.GetSize()
}

// ExportElems returns slice filled with all list's elements.
func (ul *Ulist) ExportElems() []interface{} {
	var target = []interface{}{}

	fn := func(i interface{}) {
		target = append(target, i)
	}

	ul.Do(fn)

	return target
}

// IsContains returns returns true if list contains
// at least one element val.
func (ul *Ulist) IsContains(val interface{}) bool {
	var check = false

	fn := func(i interface{}) {
		if val == i {
			check = true
		}
	}

	ul.Do(fn)

	return check
}

// PushAll appends all of the elements in the given slice vals to the end of
// the list, in the original order.
func (ul *Ulist) PushAll(vals []interface{}) {
	for i := range vals {
		ul.Push(v[i])
	}
}

/*
func (ul *Ulist) RemoveAll(vals []interface{}) {
	var n interface{}

	fn := func(i interface{}) {
		if i == n {

		}
	}
}*/
