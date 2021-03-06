# goulist

**Goulist** package is an **Urolled Linked List** written in pure Go.

[![GoDoc](https://godoc.org/github.com/hIMEI29A/goulist?status.svg)](http://godoc.org/github.com/hIMEI29A/goulist) [![Apache-2.0 License](https://img.shields.io/badge/license-Apache--2.0-red.svg)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/hIMEI29A/goulist)](https://goreportcard.com/report/github.com/hIMEI29A/goulist) [![Coverage Status](https://coveralls.io/repos/github/hIMEI29A/goulist/badge.svg?branch=master)](https://coveralls.io/github/hIMEI29A/goulist?branch=master) [![Build Status](https://travis-ci.org/hIMEI29A/goulist.svg?branch=master)](https://travis-ci.org/hIMEI29A/goulist)

## TOC
- [Version](#version)
- [Install](#install)
- [Content](#content)
- [Usage](#usage)

## Version

`1.0.3`

## Install

`go get github.com/hIMEI29A/goulist`

## Content

### Unrolled linked list

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

## Usage

Import package:

```
import (
	"github.com/hIMEI29A/goulist"
)
```
