package simplewrapper

import (
	"sort"
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/intarraytree"
)

type (
	IntArrayTree     = intarraytree.IntArrayTree
	IntArrayTreeNode = intarraytree.IntArrayTreeNode
)

const (
	HeaderSize                    = intarraytree.HeaderSize
	NodeSize                      = intarraytree.NodeSize
	PositionRootPosition          = intarraytree.PositionRootPosition
	PositionDuplicateKeysBehavior = intarraytree.PositionDuplicateKeysBehavior
	PositionIdleNodePosition      = intarraytree.PositionIdleNodePosition
	AllowDuplicateKeys            = intarraytree.AllowDuplicateKeys
	DisallowDuplicateKeys         = intarraytree.DisallowDuplicateKeys
	NodeIsNothing                 = intarraytree.NodeIsNothing
	OffsetLeftChildPosition       = intarraytree.OffsetLeftChildPosition
	OffsetRightChildPosition      = intarraytree.OffsetRightChildPosition
	OffsetHeight                  = intarraytree.OffsetHeight
	OffsetParentPosition          = intarraytree.OffsetParentPosition
	OffsetNodeCount               = intarraytree.OffsetNodeCount
	OffsetKey                     = intarraytree.OffsetKey
	OffsetValue                   = intarraytree.OffsetValue
)

func TestIntArrayTreeInsertOneEntry(t *testing.T) {

	f := func(k, v int) Tree {
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k), v)
		return tree.Tree
	}

	g := func(k, v int) Tree {
		array := make([]int, HeaderSize+NodeSize)
		node := HeaderSize
		array[PositionRootPosition] = node
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		array[node+OffsetLeftChildPosition] = NodeIsNothing
		array[node+OffsetRightChildPosition] = NodeIsNothing
		array[node+OffsetHeight] = 1
		array[node+OffsetParentPosition] = NodeIsNothing
		array[node+OffsetNodeCount] = 1
		array[node+OffsetKey] = k
		array[node+OffsetValue] = v
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertTwoEntries(t *testing.T) {

	f := func(k1, v1, k2, v2 int) Tree {
		if k1 == k2 {
			return nil
		}
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2 int) Tree {
		if k1 == k2 {
			return nil
		}
		array := make([]int, HeaderSize+NodeSize*2)
		node1 := HeaderSize
		node2 := HeaderSize + NodeSize
		array[PositionRootPosition] = node1
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		if k2 < k1 {
			array[node1+OffsetLeftChildPosition] = node2
			array[node1+OffsetRightChildPosition] = NodeIsNothing
		} else {
			array[node1+OffsetLeftChildPosition] = NodeIsNothing
			array[node1+OffsetRightChildPosition] = node2
		}
		array[node1+OffsetHeight] = 2
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 2
		array[node1+OffsetKey] = k1
		array[node1+OffsetValue] = v1
		array[node2+OffsetLeftChildPosition] = NodeIsNothing
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = node1
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k2
		array[node2+OffsetValue] = v2
		return &IntArrayTree{array}

	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeRejectDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.Insert(IntKey(k1), v3)
		tree.Insert(IntKey(k2), v4)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		array := make([]int, HeaderSize+NodeSize*2)
		node1 := HeaderSize
		node2 := HeaderSize + NodeSize
		array[PositionRootPosition] = node1
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		if k2 < k1 {
			array[node1+OffsetLeftChildPosition] = node2
			array[node1+OffsetRightChildPosition] = NodeIsNothing
		} else {
			array[node1+OffsetLeftChildPosition] = NodeIsNothing
			array[node1+OffsetRightChildPosition] = node2
		}
		array[node1+OffsetHeight] = 2
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 2
		array[node1+OffsetKey] = k1
		array[node1+OffsetValue] = v1
		array[node2+OffsetLeftChildPosition] = NodeIsNothing
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = node1
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k2
		array[node2+OffsetValue] = v2
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeReplaceDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.InsertOrReplace(IntKey(k1), v3)
		tree.InsertOrReplace(IntKey(k2), v4)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		array := make([]int, HeaderSize+NodeSize*2)
		node1 := HeaderSize
		node2 := HeaderSize + NodeSize
		array[PositionRootPosition] = node1
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		if k2 < k1 {
			array[node1+OffsetLeftChildPosition] = node2
			array[node1+OffsetRightChildPosition] = NodeIsNothing
		} else {
			array[node1+OffsetLeftChildPosition] = NodeIsNothing
			array[node1+OffsetRightChildPosition] = node2
		}
		array[node1+OffsetHeight] = 2
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 2
		array[node1+OffsetKey] = k1
		array[node1+OffsetValue] = v3
		array[node2+OffsetLeftChildPosition] = NodeIsNothing
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = node1
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k2
		array[node2+OffsetValue] = v4
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAllowDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(intarraytree.New(true))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.Insert(IntKey(k1), v3)
		tree.Insert(IntKey(k2), v4)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		array := make([]int, HeaderSize+NodeSize*4)
		node1 := HeaderSize
		node2 := node1 + NodeSize
		node3 := node2 + NodeSize
		node4 := node3 + NodeSize
		array[PositionDuplicateKeysBehavior] = AllowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		if k2 < k1 {
			array[PositionRootPosition] = node1
			array[node1+OffsetLeftChildPosition] = node2
			array[node1+OffsetRightChildPosition] = node3
			array[node1+OffsetHeight] = 3
			array[node1+OffsetParentPosition] = NodeIsNothing
			array[node1+OffsetNodeCount] = 4
			array[node1+OffsetKey] = k1
			array[node1+OffsetValue] = v1
			array[node2+OffsetLeftChildPosition] = NodeIsNothing
			array[node2+OffsetRightChildPosition] = node4
			array[node2+OffsetHeight] = 2
			array[node2+OffsetParentPosition] = node1
			array[node2+OffsetNodeCount] = 2
			array[node2+OffsetKey] = k2
			array[node2+OffsetValue] = v2
			array[node3+OffsetLeftChildPosition] = NodeIsNothing
			array[node3+OffsetRightChildPosition] = NodeIsNothing
			array[node3+OffsetHeight] = 1
			array[node3+OffsetParentPosition] = node1
			array[node3+OffsetNodeCount] = 1
			array[node3+OffsetKey] = k1
			array[node3+OffsetValue] = v3
			array[node4+OffsetLeftChildPosition] = NodeIsNothing
			array[node4+OffsetRightChildPosition] = NodeIsNothing
			array[node4+OffsetHeight] = 1
			array[node4+OffsetParentPosition] = node2
			array[node4+OffsetNodeCount] = 1
			array[node4+OffsetKey] = k2
			array[node4+OffsetValue] = v4
		} else {
			array[PositionRootPosition] = node3
			array[node3+OffsetLeftChildPosition] = node1
			array[node3+OffsetRightChildPosition] = node2
			array[node3+OffsetHeight] = 3
			array[node3+OffsetParentPosition] = NodeIsNothing
			array[node3+OffsetNodeCount] = 4
			array[node3+OffsetKey] = k1
			array[node3+OffsetValue] = v3
			array[node1+OffsetLeftChildPosition] = NodeIsNothing
			array[node1+OffsetRightChildPosition] = NodeIsNothing
			array[node1+OffsetHeight] = 1
			array[node1+OffsetParentPosition] = node3
			array[node1+OffsetNodeCount] = 1
			array[node1+OffsetKey] = k1
			array[node1+OffsetValue] = v1
			array[node2+OffsetLeftChildPosition] = NodeIsNothing
			array[node2+OffsetRightChildPosition] = node4
			array[node2+OffsetHeight] = 2
			array[node2+OffsetParentPosition] = node3
			array[node2+OffsetNodeCount] = 2
			array[node2+OffsetKey] = k2
			array[node2+OffsetValue] = v2
			array[node4+OffsetLeftChildPosition] = NodeIsNothing
			array[node4+OffsetRightChildPosition] = NodeIsNothing
			array[node4+OffsetHeight] = 1
			array[node4+OffsetParentPosition] = node2
			array[node4+OffsetNodeCount] = 1
			array[node4+OffsetKey] = k2
			array[node4+OffsetValue] = v4
		}
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		if node := tree.Min(); node != nil {
			result := []int{
				int(node.Key().(IntKey)),
				node.Value().(int),
			}
			return result
		} else {
			return nil
		}
	}

	g := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		if len(list) == 0 {
			return nil
		}
		minimum := list[0]
		for _, kv := range list[1:] {
			if kv.Key < minimum.Key {
				minimum = kv
			}
		}
		result := []int{minimum.Key, minimum.Value}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeMax(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		if node := tree.Max(); node != nil {
			result := []int{
				int(node.Key().(IntKey)),
				node.Value().(int),
			}
			return result
		} else {
			return nil
		}
	}

	g := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		if len(list) == 0 {
			return nil
		}
		maximum := list[0]
		for _, kv := range list[1:] {
			if kv.Key > maximum.Key {
				maximum = kv
			}
		}
		result := []int{maximum.Key, maximum.Value}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeCount(t *testing.T) {

	f := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		return tree.Count()
	}

	g := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		return len(list)
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeFind(t *testing.T) {

	f := func(listBase []keyAndValue) *keyAndValue {
		list := omitDuplicates(listBase)
		if len(list) < 2 {
			return nil
		}
		tree := New(intarraytree.New(false))
		for _, kv := range list[1:] {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list[1:] {
			node := tree.Find(IntKey(kv.Key))
			if node == nil {
				return kv
			}
			if kv.Key != int(node.Key().(IntKey)) {
				return kv
			}
			if kv.Value != node.Value().(int) {
				return kv
			}
		}
		return list[0]
	}

	g := func(listBase []keyAndValue) *keyAndValue {
		list := omitDuplicates(listBase)
		if len(list) < 2 {
			return nil
		} else {
			return list[0]
		}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete1(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := New(intarraytree.New(false))
		for _, list := range lists[0:2] {
			for _, kv := range list {
				tree.Insert(IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			dv := tree.Delete(IntKey(kv.Key))
			if dv == nil || dv.Value().(int) != kv.Value {
				panic("wrong")
			}
		}
		for _, kv := range lists[2] {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		list := toAscSorted(append(lists[1], lists[2]...))
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete2(t *testing.T) {

	f := func(ins1, ins2del4, ins3 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		tree := New(intarraytree.New(false))
		for _, list := range lists {
			for _, kv := range list {
				tree.Insert(IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			tree.Delete(IntKey(kv.Key))
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(ins1, ins2del4, ins3 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		list := toAscSorted(append(lists[0], lists[2]...))
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int(nil)
		tree.Iterate(func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int(nil)
		tree.IterateRev(func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toDescSorted(omitDuplicates(listBase))
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		tree.Iterate(func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			count--
			breakIteration = count <= 0
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
			count--
			if count <= 0 {
				break
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		tree.IterateRev(func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			count--
			breakIteration = count <= 0
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toDescSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
			count--
			if count <= 0 {
				break
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		appender := func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.RangeIterate(lower, upper, appender)
		tree.RangeIterate(nil, lower, appender)
		tree.RangeIterate(upper, nil, appender)
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = IntKey(k1)
			upper = IntKey(k2)
			tree.RangeIterate(lower, upper, appender)
			tree.RangeIterate(nil, lower, appender)
			tree.RangeIterate(upper, nil, appender)
		}
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		var k11, k22 int
		if len(list) > 1 {
			k11, k22 = list[0].Key, list[1].Key
		}
		toAscSorted(list)
		result := []int{}
		for _, kv := range list {
			if k1 <= kv.Key && kv.Key <= k2 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		for _, kv := range list {
			if kv.Key <= k1 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		for _, kv := range list {
			if k2 <= kv.Key {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		if len(list) > 1 {
			k1, k2 = k11, k22
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			for _, kv := range list {
				if k1 <= kv.Key && kv.Key <= k2 {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
			for _, kv := range list {
				if kv.Key <= k1 {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
			for _, kv := range list {
				if k2 <= kv.Key {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.RangeIterateRev(lower, upper, func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toDescSorted(omitDuplicates(listBase))
		result := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				continue
			}
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		stopKey := IntKey((k2 + k1) / 2)
		tree.RangeIterateRev(lower, upper, func(node Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			breakIteration = stopKey.CompareTo(node.Key()).LessThan()
			return
		})
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toDescSorted(omitDuplicates(listBase))
		result := []int{}
		stopKey := (k2 + k1) / 2
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				continue
			}
			result = append(result, kv.Key)
			result = append(result, kv.Value)
			if stopKey < kv.Key {
				break
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDuplicateKeyAscRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				tree.RangeIterate(IntKey(lower), IntKey(upper), func(node Node) (breakIteration bool) {
					values = append(values, node.Value().(int))
					return
				})
				result = append(result, values)
			}
		}
		return result
	}

	g := func(list []keyAndValue) [][]int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for key := lower; key <= upper; key++ {
					values = append(values, table[key]...)
				}
				result = append(result, values)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDuplicateKeyDescRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				tree.RangeIterateRev(IntKey(lower), IntKey(upper), func(node Node) (breakIteration bool) {
					values = append(values, node.Value().(int))
					return
				})
				result = append(result, values)
			}
		}
		return result
	}

	g := func(list []keyAndValue) [][]int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for key := lower; key <= upper; key++ {
					values = append(values, table[key]...)
				}
				result = append(result, values)
			}
		}
		for _, list := range result {
			reversed(list)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		appender := func(nodes []Node) {
			for _, node := range nodes {
				result = append(result, int(node.Key().(IntKey)))
				result = append(result, node.Value().(int))
			}
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		appender(tree.Range(lower, upper))
		appender(tree.Range(nil, lower))
		appender(tree.Range(upper, nil))
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = IntKey(k1)
			upper = IntKey(k2)
			appender(tree.Range(lower, upper))
			appender(tree.Range(nil, lower))
			appender(tree.Range(upper, nil))
		}
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		var k11, k22 int
		if len(list) > 1 {
			k11, k22 = list[0].Key, list[1].Key
		}
		toAscSorted(list)
		result := []int{}
		for _, kv := range list {
			if k1 <= kv.Key && kv.Key <= k2 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		for _, kv := range list {
			if kv.Key <= k1 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		for _, kv := range list {
			if k2 <= kv.Key {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			}
		}
		if len(list) > 1 {
			k1, k2 = k11, k22
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			for _, kv := range list {
				if k1 <= kv.Key && kv.Key <= k2 {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
			for _, kv := range list {
				if kv.Key <= k1 {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
			for _, kv := range list {
				if k2 <= kv.Key {
					result = append(result, kv.Key)
					result = append(result, kv.Value)
				}
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		for _, node := range tree.RangeRev(lower, upper) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
		}
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toDescSorted(omitDuplicates(listBase))
		result := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				continue
			}
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDuplicateKeyAscRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range tree.Range(IntKey(lower), IntKey(upper)) {
					values = append(values, node.Value().(int))
				}
				result = append(result, values)
			}
		}
		return result
	}

	g := func(list []keyAndValue) [][]int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for key := lower; key <= upper; key++ {
					values = append(values, table[key]...)
				}
				result = append(result, values)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDuplicateKeyDescRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range tree.RangeRev(IntKey(lower), IntKey(upper)) {
					values = append(values, node.Value().(int))
				}
				result = append(result, values)
			}
		}
		return result
	}

	g := func(list []keyAndValue) [][]int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for key := lower; key <= upper; key++ {
					values = append(values, table[key]...)
				}
				result = append(result, values)
			}
		}
		for _, list := range result {
			reversed(list)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeCountRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		return tree.CountRange(lower, upper)
	}

	g := func(listBase []keyAndValue, k1, k2 int) int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		result := 0
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				continue
			}
			result++
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDuplicateKeyCountRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := []int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				count := tree.CountRange(IntKey(lower), IntKey(upper))
				result = append(result, count)
			}
		}
		return result
	}

	g := func(list []keyAndValue) []int {
		table := make([]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key]++
		}
		result := []int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				count := 0
				for key := lower; key <= upper; key++ {
					count += table[key]
				}
				result = append(result, count)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDeleteAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		tree.DeleteAll(IntKey(delkey % keymax))
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			result[key] = append(result[key], node.Value().(int))
			return
		})
		return result
	}

	g := func(list []keyAndValue, delkey int) [][]int {
		result := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			result[key] = append(result[key], kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		result[delkey%keymax] = []int(nil)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeFindAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for key := 0; key < keymax; key++ {
			values := []int(nil)
			nodes := tree.FindAll(IntKey(key))
			for _, node := range nodes {
				values = append(values, node.Value().(int))
			}
			result = append(result, values)
		}
		return result
	}

	g := func(list []keyAndValue) [][]int {
		result := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			result[key] = append(result[key], kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeMinAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		nodes := tree.MinAll()
		for _, node := range nodes {
			result = append(result, node.Value().(int))
		}
		return result
	}

	g := func(list []keyAndValue) []int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := []int(nil)
		for _, values := range table {
			if len(values) > 0 {
				result = values
				break
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeMaxAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		nodes := tree.MaxAll()
		for _, node := range nodes {
			result = append(result, node.Value().(int))
		}
		return result
	}

	g := func(list []keyAndValue) []int {
		table := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			table[key] = append(table[key], kv.Value)
		}
		result := []int(nil)
		for _, values := range table {
			if len(values) > 0 {
				result = values
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeUpdateValueByFind(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			node := tree.Find(IntKey(kv.Key))
			value := node.Value().(int)
			newValue := value >> 1
			node.SetValue(newValue)
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value>>1)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeUpdateValueByIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		tree.Iterate(func(node Node) (breakIteration bool) {
			value := node.Value().(int)
			newValue := value >> 1
			node.SetValue(newValue)
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value>>1)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeUpdateValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree.Update(IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
				value := oldValue.(int)
				newValue = value >> 1
				return
			})
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value>>1)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeKeepOldValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree.Update(IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
				value := oldValue.(int)
				newValue = value >> 1
				keepOldValue = true
				return
			})
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeReplaceValue(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree.Replace(IntKey(kv.Key), kv.Value^value)
		}
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value^value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		tree.UpdateIterate(func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		tree.UpdateIterate(func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if count <= 0 || kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
			count--
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		tree.UpdateIterateRev(func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.InsertOrReplace(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		tree.UpdateIterateRev(func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		count := len(list) - (len(list)+1)/2
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if count <= 0 && kv.Value >= 0 {
				result = append(result, kv.Value>>1)
			} else {
				result = append(result, kv.Value)
			}
			count--
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRangeIterate(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRangeIterate(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 || count <= 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
			if k1 <= kv.Key && kv.Key <= k2 {
				count--
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRangeIterateRev(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRangeIterateRev(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := -(len(list) + 1) / 2
		for _, kv := range list {
			if k1 <= kv.Key && kv.Key <= k2 {
				count++
			}
		}
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 || count > 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
			if k1 <= kv.Key && kv.Key <= k2 {
				count--
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRange(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.UpdateRangeRev(lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
			value := oldValue.(int)
			newValue = value >> 1
			keepOldValue = value < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key || kv.Value < 0 {
				result = append(result, kv.Value)
			} else {
				result = append(result, kv.Value>>1)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeReplaceRange(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree.ReplaceRange(lower, upper, value)
		result := getAllAscKeyAndValues(tree)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			result = append(result, kv.Key)
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Value)
			} else {
				result = append(result, value)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeUpdateAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		count := 0
		tree.UpdateAll(IntKey(updkey%keymax), func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
			value := oldValue.(int)
			newValue = value ^ updkey
			count++
			keepOldValue = count%2 == 0
			return
		})
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			result[key] = append(result[key], node.Value().(int))
			return
		})
		return result
	}

	g := func(list []keyAndValue, updkey int) [][]int {
		result := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			result[key] = append(result[key], kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		updList := result[updkey%keymax]
		for i := range updList {
			// count++
			// keepOldValue = count%2 == 0
			// なので(i+1)%2==0がキープされる、つまりi%2==1
			if i%2 == 0 {
				updList[i] ^= updkey
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeReplaceAll(t *testing.T) {

	const keymax = 4
	const value int = 123456

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		tree.ReplaceAll(IntKey(updkey%keymax), value)
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			result[key] = append(result[key], node.Value().(int))
			return
		})
		return result
	}

	g := func(list []keyAndValue, updkey int) [][]int {
		result := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			result[key] = append(result[key], kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		updList := result[updkey%keymax]
		for i := range updList {
			updList[i] = value
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		values := tree.DeleteIterate(func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		for _, kv := range list {
			if kv.Value >= 0 {
				result = append(result, kv.Key, kv.Value)
			}
		}
		for _, kv := range list {
			if kv.Value < 0 {
				result = append(result, kv.Key, kv.Value)
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		values := tree.DeleteIterate(func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if count <= 0 || kv.Value >= 0 {
				result = append(result, kv.Key, kv.Value)
			} else {
				deleted = append(deleted, kv.Key, kv.Value)
			}
			count--
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		values := tree.DeleteIterateRev(func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Value >= 0 {
				result = append(result, kv.Key, kv.Value)
			} else {
				deleted = append(deleted, kv.Value, kv.Key)
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		values := tree.DeleteIterateRev(func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		for _, v := range values {
			result = append(result, int(v.Key().(IntKey)))
			result = append(result, v.Value().(int))
		}
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		count := len(list) - (len(list)+1)/2
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if count > 0 || kv.Value >= 0 {
				result = append(result, kv.Key, kv.Value)
			} else {
				deleted = append(deleted, kv.Value, kv.Key)
			}
			count--
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRangeIterate(lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key || kv.Value >= 0 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Key)
				deleted = append(deleted, kv.Value)
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRangeIterate(lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key || kv.Value >= 0 || count <= 0 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Key)
				deleted = append(deleted, kv.Value)
			}
			if k1 <= kv.Key && kv.Key <= k2 {
				count--
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRangeIterateRev(lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key || kv.Value >= 0 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Value)
				deleted = append(deleted, kv.Key)
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRangeIterateRev(lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := -(len(list) + 1) / 2
		for _, kv := range list {
			if k1 <= kv.Key && kv.Key <= k2 {
				count++
			}
		}
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key || kv.Value >= 0 || count > 0 {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Value)
				deleted = append(deleted, kv.Key)
			}
			if k1 <= kv.Key && kv.Key <= k2 {
				count--
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRange(lower, upper)
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Key)
				deleted = append(deleted, kv.Value)
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values := tree.DeleteRangeRev(lower, upper)
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key)
				result = append(result, kv.Value)
			} else {
				deleted = append(deleted, kv.Value)
				deleted = append(deleted, kv.Key)
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAlter(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			delValue, _ := tree.Alter(IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
				value := node.Value().(int)
				switch value & 3 {
				case 0, 3:
					if value < 0 {
						return node.Keep()
					} else {
						request.Keep()
					}
				case 1:
					if value < 0 {
						return node.Replace(value >> 1)
					} else {
						request.Replace(value >> 1)
					}
				case 2:
					if value < 0 {
						return node.Delete()
					} else {
						request.Delete()
					}
				}
				return
			})
			if delValue != nil {
				values = append(values, delValue.Value().(int))
			}
		}
		result := getAllAscKeyAndValues(tree)
		sort.Ints(values)
		result = append(result, values...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			switch kv.Value & 3 {
			case 0, 3:
				result = append(result, kv.Key, kv.Value)
			case 1:
				result = append(result, kv.Key, kv.Value>>1)
			case 2:
				deleted = append(deleted, kv.Value)
			}
		}
		sort.Ints(deleted)
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		delValues, _ := tree.AlterIterate(func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		values := toKeyValueInts(delValues)
		result = append(result, values...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			switch kv.Value & 3 {
			case 0, 3:
				result = append(result, kv.Key, kv.Value)
			case 1:
				result = append(result, kv.Key, kv.Value>>1)
			case 2:
				deleted = append(deleted, kv.Key, kv.Value)
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		delValues, _ := tree.AlterIterate(func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		values := toKeyValueInts(delValues)
		result = append(result, values...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		count := (len(list) + 1) / 2
		for i, kv := range list {
			count--
			if count < 0 {
				result = append(result, toKeyValueInts(list[i:])...)
				break
			}
			switch kv.Value & 3 {
			case 0, 3:
				result = append(result, kv.Key, kv.Value)
			case 1:
				result = append(result, kv.Key, kv.Value>>1)
			case 2:
				deleted = append(deleted, kv.Key, kv.Value)
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		delValues, _ := tree.AlterIterateRev(func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		values := toKeyValueInts(delValues)
		result = append(result, values...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			switch kv.Value & 3 {
			case 0, 3:
				result = append(result, kv.Key, kv.Value)
			case 1:
				result = append(result, kv.Key, kv.Value>>1)
			case 2:
				deleted = append(deleted, kv.Value, kv.Key)
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		delValues, _ := tree.AlterIterateRev(func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		values := toKeyValueInts(delValues)
		result = append(result, values...)
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		count := len(list) - (len(list)+1)/2
		for _, kv := range list {
			count--
			if count >= 0 {
				result = append(result, kv.Key, kv.Value)
				continue
			}
			switch kv.Value & 3 {
			case 0, 3:
				result = append(result, kv.Key, kv.Value)
			case 1:
				result = append(result, kv.Key, kv.Value>>1)
			case 2:
				deleted = append(deleted, kv.Value, kv.Key)
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRangeIterate(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Key, kv.Value)
				}
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRangeIterate(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := (len(list) + 1) / 2
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				count--
				if count < 0 {
					result = append(result, kv.Key, kv.Value)
					continue
				}
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Key, kv.Value)
				}
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRangeIterateRev(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Value, kv.Key)
				}
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRangeIterateRev(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep(), false
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1), false
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete(), false
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		count := -(len(list) + 1) / 2
		for _, kv := range list {
			if k1 <= kv.Key && kv.Key <= k2 {
				count++
			}
		}
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				count--
				if count >= 0 {
					result = append(result, kv.Key, kv.Value)
					continue
				}
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Value, kv.Key)
				}
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAscAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRange(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep()
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1)
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete()
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Key, kv.Value)
				}
			}
		}
		result = append(result, deleted...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeDescAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(intarraytree.New(false))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		values, _ := tree.AlterRangeRev(lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep()
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1)
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete()
				} else {
					request.Delete()
				}
			}
			return
		})
		result := getAllAscKeyAndValues(tree)
		result = append(result, toKeyValueInts(values)...)
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := toAscSorted(omitDuplicates(listBase))
		result := []int(nil)
		deleted := []int{}
		for _, kv := range list {
			if kv.Key < k1 || k2 < kv.Key {
				result = append(result, kv.Key, kv.Value)
			} else {
				switch kv.Value & 3 {
				case 0, 3:
					result = append(result, kv.Key, kv.Value)
				case 1:
					result = append(result, kv.Key, kv.Value>>1)
				case 2:
					deleted = append(deleted, kv.Value, kv.Key)
				}
			}
		}
		result = append(result, reversed(deleted)...)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeAlterAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(intarraytree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		values, _ := tree.AlterAll(IntKey(updkey%keymax), func(node avltree.AlterNode) (request avltree.AlterRequest) {
			value := node.Value().(int)
			switch value & 3 {
			case 0, 3:
				if value < 0 {
					return node.Keep()
				} else {
					request.Keep()
				}
			case 1:
				if value < 0 {
					return node.Replace(value >> 1)
				} else {
					request.Replace(value >> 1)
				}
			case 2:
				if value < 0 {
					return node.Delete()
				} else {
					request.Delete()
				}
			}
			return
		})
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			result[key] = append(result[key], node.Value().(int))
			return
		})
		result = append(result, toKeyValueInts(values))
		return result
	}

	g := func(list []keyAndValue, updkey int) [][]int {
		result := make([][]int, keymax)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			key %= keymax
			result[key] = append(result[key], kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		remains := []int(nil)
		deleted := []int(nil)
		for _, v := range result[updkey%keymax] {
			switch v & 3 {
			case 0, 3:
				remains = append(remains, v)
			case 1:
				remains = append(remains, v>>1)
			case 2:
				deleted = append(deleted, updkey%keymax, v)
			}
		}
		result[updkey%keymax] = remains
		result = append(result, deleted)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete3(t *testing.T) {

	f := func(k, v int) Tree {
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k), v)
		tree.Delete(IntKey(k))
		return tree.Tree
	}

	g := func(k, v int) Tree {
		array := make([]int, HeaderSize+NodeSize)
		node1 := HeaderSize
		array[PositionRootPosition] = NodeIsNothing
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = node1
		array[node1+OffsetLeftChildPosition] = NodeIsNothing
		array[node1+OffsetRightChildPosition] = NodeIsNothing
		array[node1+OffsetHeight] = 1
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 1
		array[node1+OffsetKey] = k
		array[node1+OffsetValue] = v
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete4(t *testing.T) {

	f := func(k1, v1, k2, v2 int) Tree {
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Delete(IntKey(k1))
		tree.Insert(IntKey(k2), v2)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2 int) Tree {
		array := make([]int, HeaderSize+NodeSize)
		node1 := HeaderSize
		array[PositionRootPosition] = node1
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = NodeIsNothing
		array[node1+OffsetLeftChildPosition] = NodeIsNothing
		array[node1+OffsetRightChildPosition] = NodeIsNothing
		array[node1+OffsetHeight] = 1
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 1
		array[node1+OffsetKey] = k2
		array[node1+OffsetValue] = v2
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete5(t *testing.T) {

	f := func(k1, v1, k2, v2, k3, v3 int) Tree {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.Delete(IntKey(k1))
		tree.Delete(IntKey(k2))
		tree.Insert(IntKey(k3), v3)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, k3, v3 int) Tree {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		array := make([]int, HeaderSize+NodeSize*2)
		node1 := HeaderSize
		node2 := node1 + NodeSize
		array[PositionRootPosition] = node2
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = node1
		array[node1+OffsetLeftChildPosition] = NodeIsNothing
		array[node1+OffsetRightChildPosition] = node2
		array[node1+OffsetHeight] = 2
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 2
		array[node1+OffsetKey] = k1
		array[node1+OffsetValue] = v1
		array[node2+OffsetLeftChildPosition] = NodeIsNothing
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = NodeIsNothing
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k3
		array[node2+OffsetValue] = v3
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete6(t *testing.T) {

	f := func(k1, v1, k2, v2, k3, v3 int) Tree {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.Delete(IntKey(k2))
		tree.Delete(IntKey(k1))
		tree.Insert(IntKey(k3), v3)
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, k3, v3 int) Tree {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		array := make([]int, HeaderSize+NodeSize*2)
		node1 := HeaderSize
		node2 := node1 + NodeSize
		array[PositionRootPosition] = node1
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = node2
		array[node1+OffsetLeftChildPosition] = NodeIsNothing
		array[node1+OffsetRightChildPosition] = NodeIsNothing
		array[node1+OffsetHeight] = 1
		array[node1+OffsetParentPosition] = NodeIsNothing
		array[node1+OffsetNodeCount] = 1
		array[node1+OffsetKey] = k3
		array[node1+OffsetValue] = v3
		array[node2+OffsetLeftChildPosition] = NodeIsNothing
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = node1
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k2
		array[node2+OffsetValue] = v2
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeInsertAndDelete7(t *testing.T) {

	f := func(k1, v1, k2, v2, k3, v3 int) Tree {
		ks := []int{k1, k2, k3}
		sort.Ints(ks)
		k1, k2, k3 = ks[0], ks[1], ks[2]
		tree := New(intarraytree.New(false))
		tree.Insert(IntKey(k1), v1)
		tree.Insert(IntKey(k2), v2)
		tree.Insert(IntKey(k3), v3)
		tree.Delete(IntKey(k1))
		tree.Delete(IntKey(k3))
		tree.Delete(IntKey(k2))
		return tree.Tree
	}

	g := func(k1, v1, k2, v2, k3, v3 int) Tree {
		ks := []int{k1, k2, k3}
		sort.Ints(ks)
		k1, k2, k3 = ks[0], ks[1], ks[2]
		array := make([]int, HeaderSize+NodeSize*3)
		node1 := HeaderSize
		node2 := node1 + NodeSize
		node3 := node2 + NodeSize
		array[PositionRootPosition] = NodeIsNothing
		array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		array[PositionIdleNodePosition] = node2
		array[node1+OffsetLeftChildPosition] = NodeIsNothing
		array[node1+OffsetRightChildPosition] = NodeIsNothing
		array[node1+OffsetHeight] = 1
		array[node1+OffsetParentPosition] = node2
		array[node1+OffsetNodeCount] = 1
		array[node1+OffsetKey] = k1
		array[node1+OffsetValue] = v1
		array[node2+OffsetLeftChildPosition] = node3
		array[node2+OffsetRightChildPosition] = NodeIsNothing
		array[node2+OffsetHeight] = 1
		array[node2+OffsetParentPosition] = NodeIsNothing
		array[node2+OffsetNodeCount] = 1
		array[node2+OffsetKey] = k2
		array[node2+OffsetValue] = v2
		array[node3+OffsetLeftChildPosition] = node1
		array[node3+OffsetRightChildPosition] = NodeIsNothing
		array[node3+OffsetHeight] = 1
		array[node3+OffsetParentPosition] = node2
		array[node3+OffsetNodeCount] = 1
		array[node3+OffsetKey] = k3
		array[node3+OffsetValue] = v3
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestIntArrayTreeClear(t *testing.T) {

	f := func(list []keyAndValue, allowDuplicateKeys bool) Tree {
		tree := New(intarraytree.New(allowDuplicateKeys))
		for _, kv := range list {
			tree.Insert(IntKey(kv.Key), kv.Value)
		}
		tree.Clear()
		return tree.Tree
	}

	g := func(list []keyAndValue, allowDuplicateKeys bool) Tree {
		array := make([]int, HeaderSize)
		array[PositionRootPosition] = NodeIsNothing
		if allowDuplicateKeys {
			array[PositionDuplicateKeysBehavior] = AllowDuplicateKeys
		} else {
			array[PositionDuplicateKeysBehavior] = DisallowDuplicateKeys
		}
		array[PositionIdleNodePosition] = NodeIsNothing
		return &IntArrayTree{array}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}
