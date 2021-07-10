package goulist

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var (
	nodeSize = 4
)

func Test_newUlistNode(t *testing.T) {
	var (
		newNode = &ulistNode{nil, nil, 0, nodeSize, []interface{}{nil, nil, nil, nil}}
	)

	type args struct {
		c int
	}

	tests := []struct {
		name string
		args args
		want *ulistNode
	}{
		{"newNodeTest", args{nodeSize}, newNode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUlistNode(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUlistNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: refactoring
func Test_ulistNode_add(t *testing.T) {
	var (
		node = newUlistNode(nodeSize)
	)

	var (
		toAdd       = 555
		toAffIfFull = 333
	)

	var (
		nodeAfter = &ulistNode{
			nil,
			nil,
			1,
			nodeSize,
			[]interface{}{toAdd, nil, nil, nil},
		}
	)

	var (
		halfFullNode = &ulistNode{nil,
			nil,
			3,
			nodeSize,
			[]interface{}{toAdd, toAdd, toAffIfFull, nil},
		}
	)

	var (
		halfFullNodeAfter = &ulistNode{
			nil,
			nil,
			2,
			nodeSize,
			[]interface{}{toAdd, toAdd, toAffIfFull, nil},
		}
	)

	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type newFields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type args struct {
		val interface{}
	}

	tests := []struct {
		name      string
		fields    fields
		newFields newFields
		args      args
		want      *ulistNode // returned node
		self      *ulistNode // node itself after adding
	}{
		{
			// test case of adding to empty node
			"nodeIsNotFullTest",
			fields{nil, nil, 0, nodeSize, []interface{}{nil, nil, nil, nil}},
			newFields{nil, nil, 1, nodeSize, []interface{}{toAdd, nil, nil, nil}},
			args{toAdd},
			node,
			nodeAfter,
		},

		{
			// test case of adding to full node
			"nodeIsFullTest",
			fields{nil, nil, 4, nodeSize, []interface{}{toAdd, toAdd, toAdd, toAdd}},
			newFields{
				nil,
				nil,
				2,
				nodeSize,
				[]interface{}{toAdd, toAdd, toAffIfFull, nil},
			},
			args{toAffIfFull},
			halfFullNode,
			halfFullNodeAfter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// node itself
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			// returned node
			nn := &ulistNode{
				next:     tt.newFields.next,
				prev:     tt.newFields.prev,
				size:     tt.newFields.size,
				capacity: tt.newFields.capacity,
				elems:    tt.newFields.elems,
			}

			if got := un.add(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ulistNode.add() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(nn, tt.self) {
				t.Errorf("value not added")
			}
		})
	}
}

func Test_ulistNode_del(t *testing.T) {
	var (
		errn = errors.New("Element index is out of range")
	)

	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type args struct {
		index int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			// deletion with error
			"delWithErrorTest",
			fields{nil, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{777}, // index out of range
			0,
			true,
		},

		{
			// first element deletion
			"delFirstElemTest",
			fields{nil, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{0},
			0,
			false,
		},

		{
			// second element deletion
			"delFirstElemTest",
			fields{nil, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{1},
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			got, err := un.del(tt.args.index)

			if (err != nil && err.Error() == errn.Error()) != tt.wantErr {
				t.Errorf("ulistNode.del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ulistNode.del() = %v, want %v", got, tt.want)
			}

			if (err == nil && tt.args.index == 0) && (un.elems[0] != nil) {
				t.Errorf("ulistNode.del() error: elements mismatch")
			}

			if (err == nil && tt.args.index == 1) && (un.elems[1] != nil) {
				t.Errorf("ulistNode.del() error: elements mismatch")
			}
		})
	}
}

func Test_ulistNode_delAt(t *testing.T) {
	var (
		node2 = &ulistNode{nil, nil, 2, nodeSize, []interface{}{3, 4, nil, nil}}
		node3 = &ulistNode{nil, nil, 2, nodeSize, []interface{}{5, 6, nil, nil}}
	)

	node2.next = node3
	node3.prev = node2

	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type args struct {
		index int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			// delition with error
			"delWithErrorTest",
			fields{nil, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{777}, // index out of range
			0,
			true,
		},

		{
			"delFirstElementTest",
			fields{node2, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{0},
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			got, err := un.delAt(tt.args.index)

			if (err != nil) != tt.wantErr {
				t.Errorf("ulistNode.delAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ulistNode.delAt() = %v, want %v", got, tt.want)
			}

			if err != nil &&
				err.Error() == fmt.Sprintf(
					"Element with index %d deletion error", tt.args.index) {
				t.Errorf("Element deletion error")
			}

			if (err == nil &&
				got == tt.want) &&
				(tt.args.index == 0) &&
				(un.next != node3) {
				t.Errorf("ulistNode.next == nil after successful deletion")
			}

			if ((err == nil && got == tt.want) && (tt.args.index == 0)) &&
				((un.elems[0] != 2) || (un.elems[1] != 4) || (un.elems[2] != 3)) {
				t.Errorf("Order of ulistNode.elems is wrong")
			}
		})
	}
}

func Test_ulistNode_delOccurrences(t *testing.T) {
	var (
		node2 = &ulistNode{nil, nil, 2, nodeSize, []interface{}{3, 4, nil, nil}}
		node3 = &ulistNode{nil, nil, 2, nodeSize, []interface{}{5, 6, nil, nil}}
	)

	node2.next = node3
	node3.prev = node2

	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type args struct {
		val interface{}
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"delOneElementsTest",
			fields{node2, nil, 2, nodeSize, []interface{}{1, 2, nil, nil}},
			args{1},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			if got := un.delOccurrences(tt.args.val); got != tt.want {
				t.Errorf("ulistNode.delOccurrences() = %v, want %v", got, tt.want)
			}

			if un.delOccurrences(tt.args.val); un.elems[0] != 2 ||
				un.elems[1] != 4 || un.elems[2] != 3 {
				t.Errorf("Order of ulistNode.elems is wrong after deletion")
			}
		})
	}
}

func Test_ulistNode_shift(t *testing.T) {
	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			"shiftTest1",
			fields{nil, nil, 2, nodeSize, []interface{}{3, nil, 4, nil}},
		},

		{
			"shiftTest2",
			fields{nil, nil, 1, nodeSize, []interface{}{nil, nil, 4, nil}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			un.shift()

			if un.size == 2 && (un.elems[1] != 4 || un.elems[2] != nil) {
				t.Errorf("Wrong elements order")
			}

			if un.size == 1 && (un.elems[0] != 4 || un.elems[2] != nil) {
				t.Errorf("Wrong elements order")
			}
		})
	}
}

func Test_ulistNode_do(t *testing.T) {
	f := func(i *interface{}) {
		*i = nil
	}

	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	type args struct {
		fn func(*interface{})
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"doTest",
			fields{nil, nil, 4, nodeSize, []interface{}{0, 1, 2, 3}},
			args{f},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			un.do(tt.args.fn)

			for i := 0; i < 4; i++ {
				if un.elems[i] != nil {
					t.Errorf(
						"Function %T was not called at element %d",
						tt.args.fn, i,
					)
					break
				}
			}
		})
	}
}

func Test_ulistNode_isFull(t *testing.T) {
	type fields struct {
		next     *ulistNode
		prev     *ulistNode
		size     int
		capacity int
		elems    []interface{}
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"isFullTestFalse",
			fields{nil, nil, 2, nodeSize, []interface{}{3, 45, nil, nil}},
			false,
		},

		{
			"isFullTestTrue",
			fields{nil, nil, 4, nodeSize, []interface{}{3, 45, 32, 87}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &ulistNode{
				next:     tt.fields.next,
				prev:     tt.fields.prev,
				size:     tt.fields.size,
				capacity: tt.fields.capacity,
				elems:    tt.fields.elems,
			}

			if got := un.isFull(); got != tt.want {
				t.Errorf("ulistNode.isFull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUlist(t *testing.T) {
	type args struct {
		c int
	}

	tests := []struct {
		name      string
		args      args
		wantSize  int
		elemsLen  int
		elemsSize int
		elemsVal  interface{}
	}{
		{
			"newListTest",
			args{nodeSize},
			1,
			nodeSize,
			0,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newUlist(tt.args.c)

			if got.size != tt.wantSize {
				t.Errorf("Ulist size = %d but %d needed", got.size, tt.wantSize)
			}

			if len(got.first.elems) != tt.elemsLen {
				t.Errorf(
					"Ulist node length = %d, but %d needed",
					len(got.first.elems), tt.elemsLen)
			}

			if got.first.size != tt.elemsSize {
				t.Errorf(
					"Ulist node size = %d, but %d needed",
					got.first.size, tt.elemsSize)
			}

			for i := range got.first.elems {
				if got.first.elems[i] != tt.elemsVal {
					t.Errorf("All new list's elements must be nil")
					break
				}
			}

		})
	}
}

func TestNewUlist(t *testing.T) {
	type args struct {
		c int
	}

	tests := []struct {
		name      string
		args      args
		wantSize  int
		elemsLen  int
		elemsSize int
		elemsVal  interface{}
	}{
		{
			"newListWithCacheLineSizeTest",
			args{CacheLineSize},
			1,
			CacheLineSize,
			0,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newUlist(tt.args.c)

			if got.size != tt.wantSize {
				t.Errorf("Ulist size = %d but %d needed", got.size, tt.wantSize)
			}

			if len(got.first.elems) != tt.elemsLen {
				t.Errorf(
					"Ulist node length = %d, but %d needed",
					len(got.first.elems), CacheLineSize)
			}

			if got.first.size != tt.elemsSize {
				t.Errorf(
					"Ulist node size = %d, but %d needed",
					got.first.size, tt.elemsSize)
			}

			for i := range got.first.elems {
				if got.first.elems[i] != tt.elemsVal {
					t.Errorf("All new list's elements must be nil")
					break
				}
			}
		})
	}
}

func TestNewUlistCustomCap(t *testing.T) {
	type args struct {
		c int
	}

	tests := []struct {
		name      string
		args      args
		wantSize  int
		elemsLen  int
		elemsSize int
		elemsVal  interface{}
	}{
		{
			"newListWithCustomSizeTest",
			args{nodeSize},
			1,
			nodeSize,
			0,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newUlist(tt.args.c)

			if got.size != tt.wantSize {
				t.Errorf("Ulist size = %d but %d needed", got.size, tt.wantSize)
			}

			if len(got.first.elems) != tt.elemsLen {
				t.Errorf(
					"Ulist node length = %d, but %d needed",
					len(got.first.elems), tt.elemsLen)
			}

			if got.first.size != tt.elemsSize {
				t.Errorf(
					"Ulist node size = %d, but %d needed",
					got.first.size, tt.elemsSize)
			}

			for i := range got.first.elems {
				if got.first.elems[i] != tt.elemsVal {
					t.Errorf("All new list's elements must be nil")
					break
				}
			}
		})
	}
}

func TestUlist_GetSize(t *testing.T) {
	type fields struct {
		first *ulistNode
		last  *ulistNode
		size  int
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"ulistGetSizeTest",
			fields{&ulistNode{}, &ulistNode{}, 2},
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul := &Ulist{
				first: tt.fields.first,
				last:  tt.fields.last,
				size:  tt.fields.size,
			}

			if got := ul.GetSize(); got != tt.want {
				t.Errorf("Ulist.GetSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_GetFirst(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	tests := []struct {
		name string
		want []interface{}
	}{
		{
			"getFirstTest1",
			[]interface{}{0, 1, 2, 3},
		},

		{
			"getFirstTest2",
			[]interface{}{0, 1},
		},
	}

	for _, tt := range tests {
		if tt.name == "getFirstTest1" {
			for i := 0; i < 4; i++ {
				ul.Push(i)
			}
		}

		if tt.name == "getFirstTest2" {
			for i := 0; i < 5; i++ {
				ul.Push(i)
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			if got := ul.GetFirst(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ulist.GetFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_GetLast(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	tests := []struct {
		name string
		want []interface{}
	}{
		{
			"getLastTest1",
			[]interface{}{0, 1, 2, 3},
		},

		{
			"getLastTest2",
			[]interface{}{4, 5, 6},
		},
	}

	for _, tt := range tests {
		if tt.name == "getLastTest1" {
			for i := 0; i < 4; i++ {
				ul.Push(i)
			}
		}

		if tt.name == "getLastTest2" {
			for i := 0; i < 7; i++ {
				ul.Push(i)
			}
		}

		t.Run(tt.name, func(t *testing.T) {

			if got := ul.GetLast(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ulist.GetLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_findNode(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 23; i++ {
		ul.Push(i)
	}

	type args struct {
		num int
	}

	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			"findNodeWithErrorTest",
			args{777},
			[]interface{}{},
			true,
		},

		{
			// start from begin of the list
			"findNodeTest1",
			args{2},
			[]interface{}{4, 5, nil, nil},
			false,
		},

		{
			// start from end of the list
			"findNodeTest2",
			args{9},
			[]interface{}{18, 19, nil, nil},
			false,
		},

		{
			// start from end of the list, index = 0
			"findNodeTest3WithZero",
			args{0},
			[]interface{}{0, 1, nil, nil}, // first node returned
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ul.findNode(tt.args.num)

			if tt.args.num == 777 {
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"Ulist.findNode() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
			}

			if tt.args.num == 2 {
				if !reflect.DeepEqual(got.elems, tt.want) {
					t.Errorf("Find node test1 error")
				}
			}

			if tt.args.num == 9 {
				if !reflect.DeepEqual(got.elems, tt.want) {
					t.Errorf("Find node test2 error")
				}
			}

			if tt.args.num == 0 {
				if !reflect.DeepEqual(got.elems, tt.want) {
					t.Errorf("Find node test2 error")
				}
			}
		})
	}
}

func TestUlist_Push(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	type args struct {
		val interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"pushTest1",
			args{1},
			false,
		},

		{
			"pushTest2",
			args{144},
			false,
		},

		{
			"pushTest3",
			args{6},
			false,
		},

		{
			"pushTest4",
			args{890},
			false,
		},

		{
			"pushTest5",
			args{123},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ul.Push(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Ulist.Push() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name == "pushTest1" {
				if ul.last.size != 1 {
					t.Errorf("Wrong size of list's last node")
				}
			}

			if tt.name == "pushTest4" {
				if ul.last.size != 4 {
					t.Errorf("Wrong size of list's last node")
				}
			}

			if tt.name == "pushTest5" {
				if ul.first.size != 2 || ul.last.size != 3 {
					t.Errorf("Wrong size of list's last node")
				}
			}
		})
	}
}

func TestUlist_Insert(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 4; i++ { // after it, ul.size == 1
		ul.Push(i)
	}

	type args struct {
		val interface{}
		num int
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantSize int
	}{
		{
			"insertTest1",
			args{13, 0},
			false,
			2,
		},

		{
			"insertTest3ErrorReturned",
			args{65, 2134},
			true,
			2,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := ul.Insert(tt.args.val, tt.args.num); (err != nil) !=
				tt.wantErr {
				t.Errorf("Ulist.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestUlist_ExportElems(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 8; i++ {
		ul.Push(i)
	}

	tests := []struct {
		name string
		want []interface{}
	}{
		{
			"exportElemsTest",
			[]interface{}{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ul.ExportElems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ulist.ExportElems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_IsContains(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	ul.Push(2)
	ul.Push(1)
	ul.Push(2)
	ul.Push(67)

	type args struct {
		val interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"isContainsTestTrue",
			args{2},
			true,
		},

		{
			"isContainsTestFalse",
			args{888},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ul.IsContains(tt.args.val); got != tt.want {
				t.Errorf("Ulist.IsContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_IsContainsAll(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 10; i++ {
		ul.Push(i)
	}

	s1 := []interface{}{2, 4, 6, 1, 2}
	s2 := []interface{}{2, 4, 6, 1, 999}

	type args struct {
		vals []interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"isContainsAllTrueTest",
			args{s1},
			true,
		},

		{
			"isContainsAllFalseTest",
			args{s2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ul.IsContainsAll(tt.args.vals); got != tt.want {
				t.Errorf("Ulist.IsContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_PushAll(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 10; i++ {
		ul.Push(i)
	}

	n := []interface{}{2, 4, 6, 8}
	m := []interface{}{2, 4, 6, 8}

	type args struct {
		vals []interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"pushAllTest",
			args{n},
			m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul.PushAll(tt.args.vals)

			if !reflect.DeepEqual(ul.last.elems, tt.want) {
				t.Errorf("Some elements was not added")
			}
		})
	}
}

func TestUlist_RemoveFromNode(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 10; i++ {
		ul.Push(i)
	}

	n := []interface{}{3, 5, 4, nil}

	type args struct {
		nodeNum int
		elemNum int
	}

	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			// node number 1
			"removeFromNodeOneTest",
			args{1, 0},
			n,
			false,
		},

		{
			// node number 777 - error
			"removeFromNodeOneTest",
			args{777, 0},
			n,
			true,
		},

		{
			// node number 1, element number 777 - error
			"removeFromNodeOneTest",
			args{1, 777},
			n,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ul.RemoveFromNode(
				tt.args.nodeNum,
				tt.args.elemNum,
			); ((err != nil) != tt.wantErr) ||
				!reflect.DeepEqual(
					ul.first.next.elems,
					tt.want,
				) {
				t.Errorf(
					"Ulist.RemoveFromNode() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestUlist_RemoveAllOccurrences(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(0)
	ul.Push(2)
	ul.Push(1)
	ul.Push(2)

	n := []interface{}{0, 1, nil, nil}

	type args struct {
		val interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"removeAllOccurrencesTest",
			args{2},
			n,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul.RemoveAllOccurrences(tt.args.val)

			if !reflect.DeepEqual(ul.first.elems, tt.want) {
				t.Errorf("Error of element %d deletion", tt.args.val)
			}
		})
	}
}

func TestUlist_RemoveAllOfSlice(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(0)
	ul.Push(2)
	ul.Push(1)
	ul.Push(2)
	ul.Push(3)

	s := []interface{}{2, 3}
	n := []interface{}{0, 1, nil, nil}

	type args struct {
		vals []interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"removeAllOfSliceTest",
			args{s},
			n,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul.RemoveAllOfSlice(tt.args.vals)

			if !reflect.DeepEqual(ul.first.elems, tt.want) {
				t.Errorf("Error of element %d deletion", tt.args.vals)
			}
		})
	}
}

/*
func TestUlist_Set(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)
	ul.Push(33)

	type args struct {
		nodeNum int
		elemNum int
		val     interface{}
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"setTest",
			args{0, 1, 55},
			55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ul.Set(
				tt.args.nodeNum,
				tt.args.elemNum,
				tt.args.val,
			); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Ulist.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func TestUlist_Set(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)
	ul.Push(33)

	type args struct {
		nodeNum int
		elemNum int
		val     interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"setTest",
			args{0, 1, 55},
			55,
			false,
		},

		{
			"setWithErrorTest",
			args{567, 1, 55}, // wrong node index
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ul.Set(tt.args.nodeNum, tt.args.elemNum, tt.args.val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Ulist.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ulist.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_Len(t *testing.T) {
	ul := NewUlistCustomCap(nodeSize)

	for i := 0; i < 10; i++ {
		ul.Push(i)
	}

	tests := []struct {
		name string
		want int
	}{
		{
			"lenTest",
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ul.Len(); got != tt.want {
				t.Errorf("Ulist.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_Get(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)
	ul.Push(33)

	type args struct {
		nodeNum int
		elemNum int
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"getTest",
			args{0, 1},
			33,
			false,
		},

		{
			"getTestWithError",
			args{989, 1}, // wrong node index
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ul.Get(tt.args.nodeNum, tt.args.elemNum)

			if (err != nil) != tt.wantErr {
				t.Errorf("Ulist.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ulist.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUlist_Do(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)
	ul.Push(33)

	f := func(i *interface{}) {
		*i = nil
	}

	type args struct {
		fn func(*interface{})
	}

	tests := []struct {
		name string
		args args
	}{
		{
			"uListDoTest",
			args{f},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul.Do(tt.args.fn)

			for i := 0; i < nodeSize; i++ {
				if ul.first.elems[i] != nil {
					t.Errorf(
						"Function %T was not called at elements the list",
						f,
					)
				}
			}
		})
	}
}

func TestUlist_Clear(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)
	ul.Push(55)
	ul.Push(33)
	ul.Push(999)

	tests := []struct {
		name string
	}{
		{
			"clearTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul.Clear()

			for i := 0; i < nodeSize; i++ {
				if ul.first.elems[i] != nil {
					t.Errorf("Clear() test failed")
				}
			}
		})
	}
}

func TestUlist_Printc(t *testing.T) {
	ul := NewUlistCustomCap(4)

	ul.Push(22)

	tests := []struct {
		name  string
		wantW string
	}{
		{
			"printcTest",
			fmt.Sprintf("%d\n", 22),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ul.Printc(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Log(reflect.DeepEqual(gotW, tt.wantW))
				t.Errorf("Ulist.Printc() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
