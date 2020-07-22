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
// slice of elements. Length of this slice is equial to CacheLineSize.
type ulistNode struct {
	next  *ulistNode
	prev  *ulistNode
	size  int // number of elements
	elems []interface{}
}

// newUlistNode creates empty instance of list's node.
// Elems field has a length equal to CacheLineSize.
// All elements in emtpy node is set to nil.
func newUlistNode() *ulistNode {
	elems := make([]interface{}, 0)

	// fill elems field with nil elemets
	for i := 0; i < CacheLineSize; i++ {
		elems = append(elems, nil)
	}

	return &ulistNode{
		next:  nil,
		prev:  nil,
		size:  0,
		elems: elems,
	}
}

// add sets the first non-nil element equal to the given value
// and increments size of node
func (un *ulistNode) add(val interface{}) *ulistNode {
	for i := range un.elems {
		if un.elems[i] == nil {
			un.elems[i] = val
			break
		}
	}

	un.size++

	return un
}

func (un *ulistNode) addIfFull(nn *ulistNode, val interface{}) *ulistNode {
	var (
		// elements to move
		tmv   = CacheLineSize / 2
		start = CacheLineSize - tmv
	)

	for i := 0; i < tmv; i++ {
		nn = nn.add(un.elems[start+i])
		err := un.del(start + i)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	}

	// add val to the end of new node
	nn = nn.add(val)

	return nn
}

func (un *ulistNode) del(index int) error {
	var err error

	if index > CacheLineSize {
		err = errors.New("Element index is out of range")
		return err
	}

	un.elems[index] = nil
	un.size--

	return err
}

func (un *ulistNode) isFull() bool {
	return un.size == CacheLineSize
}

// Ulist is an unrolled linked list itself.
// It contains links to first and last nodes and number of nodes.
type Ulist struct {
	first *ulistNode
	last  *ulistNode
	size  int // number of nodes
}

// NewUlist creates new empty unrolled linked list. It has only one (empty)
// node wich is first and last same time.
func NewUlist() *Ulist {
	var (
		ul   = &Ulist{}
		node = newUlistNode()
	)

	ul.first = node
	ul.last = ul.first

	ul.first.next = ul.last
	ul.last.prev = ul.first

	ul.size = 1

	return ul
}

func (ul *Ulist) GetSize() int {
	return ul.size
}

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

func (ul *Ulist) Add(val interface{}) {
	newNode := newUlistNode()

	if !ul.last.isFull() {
		ul.last = ul.last.add(val) // append
	} else {
		newNode = ul.last.addIfFull(newNode, val)

		// link new node and list's last node
		ul.last.next = newNode
		newNode.prev = ul.last

		// set new node to list's last node
		ul.last = newNode

		// increment list's size
		ul.size++
	}
}

func (ul *Ulist) AddTo(val interface{}, num int) error {
	var (
		targetNode = &ulistNode{}
		err        error
	)

	targetNode, err = ul.findNode(num)

	if err != nil {
		return err
	}

	if !targetNode.isFull() {
		targetNode = targetNode.add(val)
	} else {
		newNode := newUlistNode()
		newNode = targetNode.addIfFull(newNode, val)

		targetNode.next.prev = newNode
		newNode.next = targetNode.next
		newNode.prev = targetNode

		ul.size++
	}

	return err
}

func (ul *Ulist) Do(fn func(interface{})) {
	var (
		newNode = newUlistNode()
		count   = 0
	)

	newNode = ul.first

	for count < ul.GetSize() {
		for i := range newNode.elems {
			if newNode.elems[i] != nil {
				fn(newNode.elems[i])
			}
		}

		newNode = newNode.next
		count++
	}
}

func (ul *Ulist) Print() {
	fn := func(i interface{}) {
		fmt.Printf("%v\n", i)
	}

	ul.Do(fn)
}

func (ul *Ulist) Clear() int {
	fn := func(i interface{}) {
		i = nil
	}

	ul.Do(fn)

	return ul.GetSize()
}
