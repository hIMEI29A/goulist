# golists

**Golists** package is a collection of different kinds of linked lists written in pure Go.

[![GoDoc](https://godoc.org/github.com/hIMEI29A/golists?status.svg)](http://godoc.org/github.com/hIMEI29A/golists) [![Apache-2.0 License](https://img.shields.io/badge/license-Apache--2.0-red.svg)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/hIMEI29A/golists)](https://goreportcard.com/report/github.com/hIMEI29A/golists) [![Coverage Status](https://coveralls.io/repos/github/hIMEI29A/golists/badge.svg?branch=master)](https://coveralls.io/github/hIMEI29A/golists?branch=master)

## TOC
- [Version](#version)
- [Install](#install)
- [Content](#content)
- [Usage](#usage)
- [TODO](#todo)

## Version

`0.2.0-alpha`

**Status**

Work in progress. Only Unrolled linked list implemented.

## Install

`go get github.com/hIMEI29A/golists`

## Content

### Unrolled linked list

[github.com/hIMEI29A/golists/ulist](https://github.com/hIMEI29A/golists/tree/master/ulist)

**Unrolled linked list** (_ULL_) is an doubly linked list each node of which holds up to a certain maximum number of elements, typically just large enough so that the node fills a single cache line or a small multiple thereof.

It has a special algorithm for insertion and deletion of elements. 

As Wikipedia says, 
>To insert a new element, we simply find the node the element should
>be in and insert the element into the elements array, incrementing
>numElements. If the array is already full, we first insert a new node
>either preceding or following the current one and move half of the
>elements in the current node into it.
>
>To remove an element, we simply find the node it is in and delete it
>from the elements array, decrementing numElements. If this reduces
>the node to less than half-full, then we move elements from the next node
>to fill it back up above half. If this leaves the next node less
>than half full, then we move all its remaining elements into the
>current node, then bypass and delete it.

After deletion, all **nil** elements in the node are shifted to the right to avoid empty spaces.

Default constructor of _ULL_ creates list the length of the arrays (slices in fact) at the nodes of which is equal to cache line size. Creation with custom array (slice) length is also possible.

See [Wikipedia article](http://en.wikipedia.org/wiki/Unrolled_linked_list) for details.

**Related**

[Java implementation](https://github.com/l-tamas/Unrolled-linked-list)

[Another Golang implementation](https://github.com/ryszard/unrolledlist), but it seems does not work as expected.

### XOR linked list

Not implemented yet

### Skip list

Not implemented yet

### Vlist

Not implemented yet

## Usage

Import required subpackage:

```
import (
	"github.com/hIMEI29A/golists/ulist"
)
```