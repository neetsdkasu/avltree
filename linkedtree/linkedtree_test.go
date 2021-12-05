package linkedtree

import (
	"sort"
	"testing"
	"testing/quick"

	"avltree"
)

type IntKey = avltree.IntKey

type keyAndValue struct {
	Key   int
	Value int
}

func omitDuplicates(list []keyAndValue) []*keyAndValue {
	set := make(map[int]bool)
	result := []*keyAndValue{}
	for i := range list {
		kv := &list[i]
		if set[kv.Key] {
			continue
		}
		set[kv.Key] = true
		result = append(result, kv)
	}
	return result
}

func omitAllDuplicates(allList [][]keyAndValue) [][]*keyAndValue {
	set := make(map[int]bool)
	result := [][]*keyAndValue{}
	for _, list := range allList {
		newList := []*keyAndValue{}
		for i := range list {
			kv := &list[i]
			if set[kv.Key] {
				continue
			}
			set[kv.Key] = true
			newList = append(newList, kv)
		}
		result = append(result, newList)
	}
	return result
}

func TestInsertOneEntry(t *testing.T) {

	f := func(k, v int) *LinkedTree {
		tree := NewLinkedTree(false)
		avltree.Insert(tree, false, IntKey(k), v)
		return tree
	}

	g := func(k, v int) *LinkedTree {
		root := &linkedTreeNode{nil, nil, 1, nil, 1, IntKey(k), v}
		tree := &LinkedTree{root, false}
		return tree
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInsertTwoEntries(t *testing.T) {

	f := func(k1, v1, k2, v2 int) *LinkedTree {
		if k1 == k2 {
			return nil
		}
		tree := NewLinkedTree(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		return tree
	}

	g := func(k1, v1, k2, v2 int) *LinkedTree {
		if k1 == k2 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v1}
		child := &linkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.leftChild = child
		} else {
			root.rightChild = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	var cfg quick.Config
	cfg.MaxCount = 1000

	if err := quick.CheckEqual(f, g, &cfg); err != nil {
		t.Fatal(err)
	}
}

func TestRejectDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := NewLinkedTree(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, false, IntKey(k1), v3)
		avltree.Insert(tree, false, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v1}
		child := &linkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.leftChild = child
		} else {
			root.rightChild = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	var cfg quick.Config
	cfg.MaxCount = 1000

	if err := quick.CheckEqual(f, g, &cfg); err != nil {
		t.Fatal(err)
	}
}

func TestReplaceDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := NewLinkedTree(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, true, IntKey(k1), v3)
		avltree.Insert(tree, true, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v3}
		child := &linkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v4}
		if k2 < k1 {
			root.leftChild = child
		} else {
			root.rightChild = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	var cfg quick.Config
	cfg.MaxCount = 1000

	if err := quick.CheckEqual(f, g, &cfg); err != nil {
		t.Fatal(err)
	}
}

func TestAllowDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := NewLinkedTree(true)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, false, IntKey(k1), v3)
		avltree.Insert(tree, false, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		var root *linkedTreeNode
		if k2 < k1 {
			root = &linkedTreeNode{nil, nil, 3, nil, 4, IntKey(k1), v1}
			lChild := &linkedTreeNode{nil, nil, 2, root, 2, IntKey(k2), v2}
			rChild := &linkedTreeNode{nil, nil, 1, root, 1, IntKey(k1), v3}
			lrChild := &linkedTreeNode{nil, nil, 1, lChild, 1, IntKey(k2), v4}
			lChild.rightChild = lrChild
			root.leftChild = lChild
			root.rightChild = rChild
		} else {
			root = &linkedTreeNode{nil, nil, 3, nil, 4, IntKey(k1), v3}
			lChild := &linkedTreeNode{nil, nil, 1, root, 1, IntKey(k1), v1}
			rChild := &linkedTreeNode{nil, nil, 2, root, 2, IntKey(k2), v2}
			rrChild := &linkedTreeNode{nil, nil, 1, rChild, 1, IntKey(k2), v4}
			rChild.rightChild = rrChild
			root.leftChild = lChild
			root.rightChild = rChild
		}
		tree := &LinkedTree{root, true}
		return tree
	}

	var cfg quick.Config
	cfg.MaxCount = 1000

	if err := quick.CheckEqual(f, g, &cfg); err != nil {
		t.Fatal(err)
	}
}

func TestNodeCount(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			avltree.Delete(tree, IntKey(list[0].Key))
		}
		var invalidNode Node = nil
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
			count := node.(avltree.NodeCounter).NodeCount()
			var cLeft, cRight int
			if leftChild, ok := node.LeftChild().(avltree.NodeCounter); ok {
				cLeft = leftChild.NodeCount()
			}
			if rightChild, ok := node.RightChild().(avltree.NodeCounter); ok {
				cRight = rightChild.NodeCount()
			}
			if count != 1+cLeft+cRight {
				invalidNode = node
				breakIteration = true
				return
			}
			return
		})
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestHeight(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			avltree.Delete(tree, IntKey(list[0].Key))
		}
		var invalidNode Node = nil
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
			height := node.(avltree.RealNode).Height()
			var hLeft, hRight int
			if leftChild, ok := node.LeftChild().(avltree.RealNode); ok {
				hLeft = leftChild.Height()
			}
			if rightChild, ok := node.RightChild().(avltree.RealNode); ok {
				hRight = rightChild.Height()
			}
			hMin, hMax := hLeft, hRight
			if hMax < hMin {
				hMin, hMax = hMax, hMin
			}
			if 1 < hMax-hMin {
				invalidNode = node
				breakIteration = true
				return
			}
			if height-hMax != 1 {
				invalidNode = node
				breakIteration = true
				return
			}
			return
		})
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestParent(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			avltree.Delete(tree, IntKey(list[0].Key))
		}
		if root, ok := tree.Root().(avltree.ParentGetter); ok {
			if root.Parent() != nil {
				return root
			}
		}
		var invalidNode Node = nil
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
			if parent := node.(avltree.ParentGetter).Parent(); parent != nil {
				leftChild := parent.LeftChild()
				rightChild := parent.RightChild()
				if leftChild != node && rightChild != node {
					invalidNode = node
					breakIteration = true
					return
				}
			} else if tree.Root() != node {
				invalidNode = node
				breakIteration = true
				return
			}
			if leftChild, ok := node.LeftChild().(avltree.ParentGetter); ok {
				if node != leftChild.Parent() {
					invalidNode = leftChild
					breakIteration = true
					return
				}
			}
			if rightChild, ok := node.RightChild().(avltree.ParentGetter); ok {
				if node != rightChild.Parent() {
					invalidNode = rightChild
					breakIteration = true
					return
				}
			}
			return
		})
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if node, ok := avltree.Min(tree); ok {
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

func TestMax(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if node, ok := avltree.Max(tree); ok {
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

func TestCount(t *testing.T) {

	f := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		return avltree.Count(tree)
	}

	g := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		return len(list)
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {

	f := func(listBase []keyAndValue) *keyAndValue {
		list := omitDuplicates(listBase)
		if len(list) < 2 {
			return nil
		}
		tree := NewLinkedTree(false)
		for _, kv := range list[1:] {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list[1:] {
			node, ok := avltree.Find(tree, IntKey(kv.Key))
			if !ok {
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

func TestInsertAndDelete1(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := NewLinkedTree(false)
		for _, list := range lists[0:2] {
			for _, kv := range list {
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			avltree.Delete(tree, IntKey(kv.Key))
		}
		for _, kv := range lists[2] {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		list := append(lists[1], lists[2]...)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
		})
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}
func TestInsertAndDelete2(t *testing.T) {

	f := func(ins1, ins2del4, ins3 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		tree := NewLinkedTree(false)
		for _, list := range lists {
			for _, kv := range list {
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			avltree.Delete(tree, IntKey(kv.Key))
		}
		result := []int{}
		avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(ins1, ins2del4, ins3 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		list := append(lists[0], lists[2]...)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
		})
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscSorted(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
		})
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescSorted(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		avltree.Iterate(tree, true, func(node avltree.Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		})
		return result
	}

	g := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key > list[j].Key
		})
		result := []int{}
		for _, kv := range list {
			result = append(result, kv.Key)
			result = append(result, kv.Value)
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		appender := func(node avltree.Node) (breakIteration bool) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
			return
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.RangeIterate(tree, false, lower, upper, appender)
		avltree.RangeIterate(tree, false, nil, lower, appender)
		avltree.RangeIterate(tree, false, upper, nil, appender)
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = IntKey(k1)
			upper = IntKey(k2)
			avltree.RangeIterate(tree, false, lower, upper, appender)
			avltree.RangeIterate(tree, false, nil, lower, appender)
			avltree.RangeIterate(tree, false, upper, nil, appender)
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
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
		})
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

func TestDescRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.RangeIterate(tree, true, lower, upper, func(node avltree.Node) (breakIteration bool) {
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
		list := omitDuplicates(listBase)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key > list[j].Key
		})
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

func TestDuplicateKeyAscRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				avltree.RangeIterate(tree, false, IntKey(lower), IntKey(upper), func(node Node) (breakIteration bool) {
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

func TestDuplicateKeyDescRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				avltree.RangeIterate(tree, true, IntKey(lower), IntKey(upper), func(node Node) (breakIteration bool) {
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
			for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
				list[i], list[j] = list[j], list[i]
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		appender := func(nodes []avltree.Node) {
			for _, node := range nodes {
				result = append(result, int(node.Key().(IntKey)))
				result = append(result, node.Value().(int))
			}
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		appender(avltree.Range(tree, false, lower, upper))
		appender(avltree.Range(tree, false, nil, lower))
		appender(avltree.Range(tree, false, upper, nil))
		if len(list) > 1 {
			k1, k2 = list[0].Key, list[1].Key
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			lower = IntKey(k1)
			upper = IntKey(k2)
			appender(avltree.Range(tree, false, lower, upper))
			appender(avltree.Range(tree, false, nil, lower))
			appender(avltree.Range(tree, false, upper, nil))
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
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
		})
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

func TestDescRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		for _, node := range avltree.Range(tree, true, lower, upper) {
			result = append(result, int(node.Key().(IntKey)))
			result = append(result, node.Value().(int))
		}
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key > list[j].Key
		})
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

func TestDuplicateKeyAscRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range avltree.Range(tree, false, IntKey(lower), IntKey(upper)) {
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

func TestDuplicateKeyDescRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				values := []int(nil)
				for _, node := range avltree.Range(tree, true, IntKey(lower), IntKey(upper)) {
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
			for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
				list[i], list[j] = list[j], list[i]
			}
		}
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestCountRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		return avltree.CountRange(tree, lower, upper)
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

func TestDuplicateKeyCountRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := []int{}
		for lower := 0; lower < keymax; lower++ {
			for upper := lower; upper < keymax; upper++ {
				count := avltree.CountRange(tree, IntKey(lower), IntKey(upper))
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

func TestFindAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for key := 0; key < keymax; key++ {
			values := []int(nil)
			if nodes, ok := avltree.FindAll(tree, IntKey(key)); ok {
				for _, node := range nodes {
					values = append(values, node.Value().(int))
				}
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

func TestMinAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		if nodes, ok := avltree.MinAll(tree); ok {
			for _, node := range nodes {
				result = append(result, node.Value().(int))
			}
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

func TestMaxAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		if nodes, ok := avltree.MaxAll(tree); ok {
			for _, node := range nodes {
				result = append(result, node.Value().(int))
			}
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
