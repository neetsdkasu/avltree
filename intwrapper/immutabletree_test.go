package intwrapper

import (
	"math"
	"sort"
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree/immutabletree"
)

type (
	ImmutableTree     = immutabletree.ImmutableTree
	ImmutableTreeNode = immutabletree.ImmutableTreeNode
)

func TestImmutableTreeInsertOneEntry(t *testing.T) {

	f := func(k, v int) []Tree {
		tree := New(immutabletree.New(false))
		tree0 := tree.Tree
		tree.Insert(k, v)
		tree1 := tree.Tree
		return []Tree{tree0, tree1}
	}

	g := func(k, v int) []Tree {
		root := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k), v}
		tree1 := &ImmutableTree{root, false}
		return []Tree{
			&ImmutableTree{nil, false},
			tree1,
		}
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeInsertTwoEntries(t *testing.T) {

	f := func(k1, v1, k2, v2 int) []Tree {
		if k1 == k2 {
			return nil
		}
		tree := New(immutabletree.New(false))
		tree0 := tree.Tree
		tree.Insert(k1, v1)
		tree1 := tree.Tree
		tree.Insert(k2, v2)
		tree2 := tree.Tree
		return []Tree{tree0, tree1, tree2}
	}

	g := func(k1, v1, k2, v2 int) []Tree {
		if k1 == k2 {
			return nil
		}
		root := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v1}
		child := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.LeftChildNode = child
		} else {
			root.RightChildNode = child
		}
		tree2 := &ImmutableTree{root, false}
		return []Tree{
			&ImmutableTree{nil, false},
			&ImmutableTree{&ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}, false},
			tree2,
		}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeRejectDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(immutabletree.New(false))
		tree0 := tree.Tree
		tree.Insert(k1, v1)
		tree1 := tree.Tree
		tree.Insert(k2, v2)
		tree2 := tree.Tree
		tree.Insert(k1, v3)
		tree3 := tree.Tree
		tree.Insert(k2, v4)
		tree4 := tree.Tree
		return []Tree{tree0, tree1, tree2, tree3, tree4}
	}

	g := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v1}
		child := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.LeftChildNode = child
		} else {
			root.RightChildNode = child
		}
		tree2 := &ImmutableTree{root, false}
		return []Tree{
			&ImmutableTree{nil, false},
			&ImmutableTree{&ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}, false},
			tree2,
			tree2,
			tree2,
		}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeReplaceDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(immutabletree.New(false))
		tree0 := tree.Tree
		tree.Insert(k1, v1)
		tree1 := tree.Tree
		tree.Insert(k2, v2)
		tree2 := tree.Tree
		tree.InsertOrReplace(k1, v3)
		tree3 := tree.Tree
		tree.InsertOrReplace(k2, v4)
		tree4 := tree.Tree
		return []Tree{tree0, tree1, tree2, tree3, tree4}
	}

	g := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root2 := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v1}
		child2 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
		if k2 < k1 {
			root2.LeftChildNode = child2
		} else {
			root2.RightChildNode = child2
		}
		tree2 := &ImmutableTree{root2, false}
		root3 := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v3}
		child3 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
		if k2 < k1 {
			root3.LeftChildNode = child3
		} else {
			root3.RightChildNode = child3
		}
		tree3 := &ImmutableTree{root3, false}
		root4 := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v3}
		child4 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v4}
		if k2 < k1 {
			root4.LeftChildNode = child4
		} else {
			root4.RightChildNode = child4
		}
		tree4 := &ImmutableTree{root4, false}
		return []Tree{
			&ImmutableTree{nil, false},
			&ImmutableTree{&ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}, false},
			tree2,
			tree3,
			tree4,
		}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeAllowDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(immutabletree.New(true))
		tree0 := tree.Tree
		tree.Insert(k1, v1)
		tree1 := tree.Tree
		tree.Insert(k2, v2)
		tree2 := tree.Tree
		tree.Insert(k1, v3)
		tree3 := tree.Tree
		tree.Insert(k2, v4)
		tree4 := tree.Tree
		return []Tree{tree0, tree1, tree2, tree3, tree4}
	}

	g := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		var root2, root3, root4 *ImmutableTreeNode
		if k2 < k1 {
			root2 = &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v1}
			lChild2 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
			root2.LeftChildNode = lChild2
			root3 = &ImmutableTreeNode{nil, nil, 2, 3, IntKey(k1), v1}
			lChild3 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
			rChild3 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v3}
			root3.LeftChildNode = lChild3
			root3.RightChildNode = rChild3
			root4 = &ImmutableTreeNode{nil, nil, 3, 4, IntKey(k1), v1}
			lChild4 := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k2), v2}
			rChild4 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v3}
			lrChild4 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v4}
			lChild4.RightChildNode = lrChild4
			root4.LeftChildNode = lChild4
			root4.RightChildNode = rChild4
		} else {
			root2 = &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k1), v1}
			rChild2 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
			root2.RightChildNode = rChild2
			root3 = &ImmutableTreeNode{nil, nil, 2, 3, IntKey(k1), v3}
			lChild3 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}
			rChild3 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v2}
			root3.LeftChildNode = lChild3
			root3.RightChildNode = rChild3
			root4 = &ImmutableTreeNode{nil, nil, 3, 4, IntKey(k1), v3}
			lChild4 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}
			rChild4 := &ImmutableTreeNode{nil, nil, 2, 2, IntKey(k2), v2}
			rrChild4 := &ImmutableTreeNode{nil, nil, 1, 1, IntKey(k2), v4}
			rChild4.RightChildNode = rrChild4
			root4.LeftChildNode = lChild4
			root4.RightChildNode = rChild4
		}
		return []Tree{
			&ImmutableTree{nil, true},
			&ImmutableTree{&ImmutableTreeNode{nil, nil, 1, 1, IntKey(k1), v1}, true},
			&ImmutableTree{root2, true},
			&ImmutableTree{root3, true},
			&ImmutableTree{root4, true},
		}
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		if node := tree.Min(); node != nil {
			result := []int{
				node.Key(),
				node.Value(),
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

func TestImmutableTreeMax(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		if node := tree.Max(); node != nil {
			result := []int{
				node.Key(),
				node.Value(),
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

func TestImmutableTreeCount(t *testing.T) {

	f := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
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

func TestImmutableTreeFind(t *testing.T) {

	f := func(listBase []keyAndValue) *keyAndValue {
		list := omitDuplicates(listBase)
		if len(list) < 2 {
			return nil
		}
		tree := New(immutabletree.New(false))
		for _, kv := range list[1:] {
			tree.Insert(kv.Key, kv.Value)
		}
		for _, kv := range list[1:] {
			node := tree.Find(kv.Key)
			if node == nil {
				return kv
			}
			if kv.Key != node.Key() {
				return kv
			}
			if kv.Value != node.Value() {
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

func TestImmutableTreeInsertAndDelete1(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := New(immutabletree.New(false))
		for _, list := range lists[0:2] {
			for _, kv := range list {
				tree.Insert(kv.Key, kv.Value)
			}
		}
		for _, kv := range lists[0] {
			var dv KeyAndValue
			dv = tree.Delete(kv.Key)
			if dv == nil || dv.Value() != kv.Value {
				panic("wrong")
			}
		}
		for _, kv := range lists[2] {
			tree.Insert(kv.Key, kv.Value)
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

func TestImmutableTreeInsertAndDelete2(t *testing.T) {

	f := func(ins1, ins2del4, ins3 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		tree := New(immutabletree.New(false))
		for _, list := range lists {
			for _, kv := range list {
				tree.Insert(kv.Key, kv.Value)
			}
		}
		for _, kv := range lists[1] {
			tree.Delete(kv.Key)
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

func TestImmutableTreeAscIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int(nil)
		tree.Iterate(func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeDescIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int(nil)
		tree.IterateRev(func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeAscHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		tree.Iterate(func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeDescHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		tree.IterateRev(func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeAscRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int{}
		appender := func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
			return
		}
		lower := k1
		upper := k2
		tree.RangeIterate(lower, upper, appender)
		tree.RangeIterate(math.MinInt, lower, appender)
		tree.RangeIterate(upper, math.MaxInt, appender)
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = k1
			upper = k2
			tree.RangeIterate(lower, upper, appender)
			tree.RangeIterate(math.MinInt, lower, appender)
			tree.RangeIterate(upper, math.MaxInt, appender)
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

func TestImmutableTreeDescRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int{}
		lower := k1
		upper := k2
		tree.RangeIterateRev(lower, upper, func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeDescHalfRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int{}
		lower := k1
		upper := k2
		stopKey := (k2 + k1) / 2
		tree.RangeIterateRev(lower, upper, func(node Node) (breakIteration bool) {
			result = append(result, node.Key())
			result = append(result, node.Value())
			breakIteration = stopKey < node.Key()
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

func TestImmutableTreeDuplicateKeyAscRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				tree.RangeIterate(lower, upper, func(node Node) (breakIteration bool) {
					values = append(values, node.Value())
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

func TestImmutableTreeDuplicateKeyDescRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				tree.RangeIterateRev(lower, upper, func(node Node) (breakIteration bool) {
					values = append(values, node.Value())
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

func TestImmutableTreeAscRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int{}
		appender := func(nodes []Node) {
			for _, node := range nodes {
				result = append(result, node.Key())
				result = append(result, node.Value())
			}
		}
		lower := k1
		upper := k2
		appender(tree.Range(lower, upper))
		appender(tree.Range(math.MinInt, lower))
		appender(tree.Range(upper, math.MaxInt))
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = k1
			upper = k2
			appender(tree.Range(lower, upper))
			appender(tree.Range(math.MinInt, lower))
			appender(tree.Range(upper, math.MaxInt))
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

func TestImmutableTreeDescRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		result := []int{}
		lower := k1
		upper := k2
		for _, node := range tree.RangeRev(lower, upper) {
			result = append(result, node.Key())
			result = append(result, node.Value())
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

func TestImmutableTreeDuplicateKeyAscRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range tree.Range(lower, upper) {
					values = append(values, node.Value())
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

func TestImmutableTreeDuplicateKeyDescRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range tree.RangeRev(lower, upper) {
					values = append(values, node.Value())
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

func TestImmutableTreeCountRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
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

func TestImmutableTreeDuplicateKeyCountRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := []int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				count := tree.CountRange(lower, upper)
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

func TestImmutableTreeDeleteAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		beforeTree := New(tree.Tree)
		delValues := tree.DeleteAll(delkey % keymax)
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := node.Key()
			result[key] = append(result[key], node.Value())
			return
		})
		result = append(result, toKeyValueInts(delValues))
		beforeValues := []int(nil)
		beforeTree.Iterate(func(node Node) (breakIteration bool) {
			key := node.Key()
			beforeValues = append(beforeValues, key, node.Value())
			return
		})
		result = append(result, beforeValues)
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
		delValues := []int(nil)
		for _, v := range result[delkey%keymax] {
			delValues = append(delValues, delkey%keymax, v)
		}
		beforeValues := []int(nil)
		for key, row := range result {
			for _, v := range row {
				beforeValues = append(beforeValues, key, v)
			}
		}
		result[delkey%keymax] = []int(nil)
		result = append(result, delValues, beforeValues)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestImmutableTreeFindAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := [][]int{}
		for key := 0; key < keymax; key++ {
			values := []int(nil)
			nodes := tree.FindAll(key)
			for _, node := range nodes {
				values = append(values, node.Value())
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

func TestImmutableTreeMinAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := []int(nil)
		nodes := tree.MinAll()
		for _, node := range nodes {
			result = append(result, node.Value())
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

func TestImmutableTreeMaxAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		result := []int(nil)
		nodes := tree.MaxAll()
		for _, node := range nodes {
			result = append(result, node.Value())
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

func TestImmutableTreeUpdateValueByFind(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		// Immutableなのでこの方法では更新できない
		for _, kv := range list {
			node := tree.Find(kv.Key)
			value := node.Value()
			newValue := value >> 1
			node.SetValue(newValue)
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

func TestImmutableTreeUpdateValueByIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		// Immutableなのでこの方法では更新できない
		tree.IterateRev(func(node Node) (breakIteration bool) {
			value := node.Value()
			newValue := value >> 1
			node.SetValue(newValue)
			return
		})
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

func TestImmutableTreeUpdateValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		for _, kv := range list {
			tree.Update(kv.Key, func(key, oldValue int) (newValue int, keepOldValue bool) {
				value := oldValue
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

func TestImmutableTreeKeepOldValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		for _, kv := range list {
			tree.Update(kv.Key, func(key, oldValue int) (newValue int, keepOldValue bool) {
				value := oldValue
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

func TestImmutableTreeReplaceValue(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		for _, kv := range list {
			tree.Replace(kv.Key, kv.Value^value)
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

func TestImmutableTreeAscUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		tree.UpdateIterate(func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeAscHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		tree.UpdateIterate(func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeDescUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		tree.UpdateIterateRev(func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeDescHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.InsertOrReplace(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		tree.UpdateIterateRev(func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeAscUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		tree.UpdateRangeIterate(lower, upper, func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeAscHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		tree.UpdateRangeIterate(lower, upper, func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeDescUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		tree.UpdateRangeIterateRev(lower, upper, func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeDescHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		tree.UpdateRangeIterateRev(lower, upper, func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool) {
			value := oldValue
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

func TestImmutableTreeAscUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		tree.UpdateRange(lower, upper, func(key, oldValue int) (newValue int, keepOldValue bool) {
			value := oldValue
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

func TestImmutableTreeDescUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		tree.UpdateRangeRev(lower, upper, func(key, oldValue int) (newValue int, keepOldValue bool) {
			value := oldValue
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

func TestImmutableTreeReplaceRange(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
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

func TestImmutableTreeUpdateAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		count := 0
		tree.UpdateAll(updkey%keymax, func(key, oldValue int) (newValue int, keepOldValue bool) {
			value := oldValue
			newValue = value ^ updkey
			count++
			keepOldValue = count%2 == 0
			return
		})
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := node.Key()
			result[key] = append(result[key], node.Value())
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

func TestImmutableTreeReplaceAll(t *testing.T) {

	const keymax = 4
	const value int = 123456

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		tree.ReplaceAll(updkey%keymax, value)
		result := make([][]int, keymax)
		tree.Iterate(func(node Node) (breakIteration bool) {
			key := node.Key()
			result[key] = append(result[key], node.Value())
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

func TestImmutableTreeAscDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		values := tree.DeleteIterate(func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeAscHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		values := tree.DeleteIterate(func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeDescDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		values := tree.DeleteIterateRev(func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeDescHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		values := tree.DeleteIterateRev(func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(tree)
		for _, v := range values {
			result = append(result, v.Key())
			result = append(result, v.Value())
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

func TestImmutableTreeAscDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values := tree.DeleteRangeIterate(lower, upper, func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeAscHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		values := tree.DeleteRangeIterate(lower, upper, func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeDescDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values := tree.DeleteRangeIterateRev(lower, upper, func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeDescHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		values := tree.DeleteRangeIterateRev(lower, upper, func(key, value int) (deleteNode, breakIteration bool) {
			deleteNode = value < 0
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

func TestImmutableTreeAscDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
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

func TestImmutableTreeDescDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
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

func TestImmutableTreeAlter(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			var delValue KeyAndValue
			delValue, _ = tree.Alter(kv.Key, func(node AlterNode) (request AlterRequest) {
				value := node.Value()
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
				values = append(values, delValue.Value())
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

func TestImmutableTreeAscAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		delValues, _ := tree.AlterIterate(func(node AlterNode) (request AlterRequest, breakIteration bool) {
			value := node.Value()
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

func TestImmutableTreeAscHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		delValues, _ := tree.AlterIterate(func(node AlterNode) (request AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value()
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

func TestImmutableTreeDescAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		delValues, _ := tree.AlterIterateRev(func(node AlterNode) (request AlterRequest, breakIteration bool) {
			value := node.Value()
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

func TestImmutableTreeDescHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		delValues, _ := tree.AlterIterateRev(func(node AlterNode) (request AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value()
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

func TestImmutableTreeAscAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values, _ := tree.AlterRangeIterate(lower, upper, func(node AlterNode) (request AlterRequest, breakIteration bool) {
			value := node.Value()
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

func TestImmutableTreeAscHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		values, _ := tree.AlterRangeIterate(lower, upper, func(node AlterNode) (request AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value()
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

func TestImmutableTreeDescAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values, _ := tree.AlterRangeIterateRev(lower, upper, func(node AlterNode) (request AlterRequest, breakIteration bool) {
			value := node.Value()
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

func TestImmutableTreeDescHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := k1
		upper := k2
		values, _ := tree.AlterRangeIterateRev(lower, upper, func(node AlterNode) (request AlterRequest, breakIteration bool) {
			count--
			if count < 0 {
				breakIteration = true
				return
			}
			value := node.Value()
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

func TestImmutableTreeAscAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values, _ := tree.AlterRange(lower, upper, func(node AlterNode) (request AlterRequest) {
			value := node.Value()
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

func TestImmutableTreeDescAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(immutabletree.New(false))
		for _, kv := range list {
			tree.Insert(kv.Key, kv.Value)
		}
		lower := k1
		upper := k2
		values, _ := tree.AlterRangeRev(lower, upper, func(node AlterNode) (request AlterRequest) {
			value := node.Value()
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

func TestImmutableTreeAlterAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(immutabletree.New(true))
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree.Insert(key%keymax, kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		values, _ := tree.AlterAll(updkey%keymax, func(node AlterNode) (request AlterRequest) {
			value := node.Value()
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
			key := node.Key()
			result[key] = append(result[key], node.Value())
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
