package linkedtree

import (
	"avltree"
	"sort"
	"testing"
	"testing/quick"
)

type intKey int

func (key intKey) CompareTo(other avltree.Key) avltree.KeyOrdering {
	v1 := int(key)
	v2 := int(other.(intKey))
	switch {
	case v1 < v2:
		return avltree.LessThanOtherKey
	case v1 > v2:
		return avltree.GreaterThanOtherKey
	default:
		return avltree.EqualToOtherKey
	}
	// return v1 - v2 は 算術オーバーフローがこわい
}

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

func TestInsertOneEntry(t *testing.T) {

	f := func(k, v int) *LinkedTree {
		tree := NewLinkedTree(false)
		avltree.Insert(tree, false, intKey(k), v)
		return tree
	}

	g := func(k, v int) *LinkedTree {
		root := &linkedTreeNode{nil, nil, 1, 1, intKey(k), v}
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
		avltree.Insert(tree, false, intKey(k1), v1)
		avltree.Insert(tree, false, intKey(k2), v2)
		return tree
	}

	g := func(k1, v1, k2, v2 int) *LinkedTree {
		if k1 == k2 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, 2, intKey(k1), v1}
		child := &linkedTreeNode{nil, nil, 1, 1, intKey(k2), v2}
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
		avltree.Insert(tree, false, intKey(k1), v1)
		avltree.Insert(tree, false, intKey(k2), v2)
		avltree.Insert(tree, false, intKey(k1), v3)
		avltree.Insert(tree, false, intKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, 2, intKey(k1), v1}
		child := &linkedTreeNode{nil, nil, 1, 1, intKey(k2), v2}
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
		avltree.Insert(tree, false, intKey(k1), v1)
		avltree.Insert(tree, false, intKey(k2), v2)
		avltree.Insert(tree, true, intKey(k1), v3)
		avltree.Insert(tree, true, intKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &linkedTreeNode{nil, nil, 2, 2, intKey(k1), v3}
		child := &linkedTreeNode{nil, nil, 1, 1, intKey(k2), v4}
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
		avltree.Insert(tree, false, intKey(k1), v1)
		avltree.Insert(tree, false, intKey(k2), v2)
		avltree.Insert(tree, false, intKey(k1), v3)
		avltree.Insert(tree, false, intKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) *LinkedTree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		var root *linkedTreeNode
		if k2 < k1 {
			root = &linkedTreeNode{nil, nil, 3, 4, intKey(k1), v1}
			lChild := &linkedTreeNode{nil, nil, 2, 2, intKey(k2), v2}
			rChild := &linkedTreeNode{nil, nil, 1, 1, intKey(k1), v3}
			lrChild := &linkedTreeNode{nil, nil, 1, 1, intKey(k2), v4}
			lChild.rightChild = lrChild
			root.leftChild = lChild
			root.rightChild = rChild
		} else {
			root = &linkedTreeNode{nil, nil, 3, 4, intKey(k1), v3}
			lChild := &linkedTreeNode{nil, nil, 1, 1, intKey(k1), v1}
			rChild := &linkedTreeNode{nil, nil, 2, 2, intKey(k2), v2}
			rrChild := &linkedTreeNode{nil, nil, 1, 1, intKey(k2), v4}
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

func TestMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		if node, ok := avltree.Min(tree); ok {
			result := []int{
				int(node.Key().(intKey)),
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
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		if node, ok := avltree.Max(tree); ok {
			result := []int{
				int(node.Key().(intKey)),
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

func TestNodeCount(t *testing.T) {

	f := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		for _, kv := range list[1:] {
			node, ok := avltree.Find(tree, intKey(kv.Key))
			if !ok {
				return kv
			}
			if kv.Key != int(node.Key().(intKey)) {
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

func TestAscSorted(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		result := []int{}
		avltree.Iterate(tree, false, func(node avltree.Node) bool {
			result = append(result, int(node.Key().(intKey)))
			result = append(result, node.Value().(int))
			return true
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
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		result := []int{}
		avltree.Iterate(tree, true, func(node avltree.Node) bool {
			result = append(result, int(node.Key().(intKey)))
			result = append(result, node.Value().(int))
			return true
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

func TestAscRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := intKey(k1)
		upper := intKey(k2)
		avltree.Range(tree, false, lower, upper, func(node avltree.Node) bool {
			result = append(result, int(node.Key().(intKey)))
			result = append(result, node.Value().(int))
			return true
		})
		return result
	}

	g := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Key < list[j].Key
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

func TestDescRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := NewLinkedTree(false)
		for _, kv := range list {
			avltree.Insert(tree, false, intKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := intKey(k1)
		upper := intKey(k2)
		avltree.Range(tree, true, lower, upper, func(node avltree.Node) bool {
			result = append(result, int(node.Key().(intKey)))
			result = append(result, node.Value().(int))
			return true
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

func TestDuplicateKeyRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := NewLinkedTree(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, intKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for key := 0; key < keymax; key++ {
			values := []int(nil)
			avltree.Range(tree, false, intKey(key), intKey(key), func(node Node) bool {
				values = append(values, node.Value().(int))
				return true
			})
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
