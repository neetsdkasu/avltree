package immutabletree

import (
	"sort"
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

type (
	Key      = avltree.Key
	Node     = avltree.Node
	RealNode = avltree.RealNode
	Tree     = avltree.Tree
)

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

func reversed(list []int) []int {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func toAscSorted(list []*keyAndValue) []*keyAndValue {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Key < list[j].Key
	})
	return list
}

func toDescSorted(list []*keyAndValue) []*keyAndValue {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Key > list[j].Key
	})
	return list
}

func toKeyValueInts(list interface{}) (result []int) {
	switch list := list.(type) {
	case []keyAndValue:
		for _, kv := range list {
			result = append(result, kv.Key, kv.Value)
		}
	case []*keyAndValue:
		for _, kv := range list {
			result = append(result, kv.Key, kv.Value)
		}
	case []Node:
		for _, kv := range list {
			result = append(result, int(kv.Key().(IntKey)))
			result = append(result, kv.Value().(int))
		}
	case []avltree.KeyAndValue:
		for _, kv := range list {
			result = append(result, int(kv.Key().(IntKey)))
			result = append(result, kv.Value().(int))
		}
	default:
		panic("unsupported type")
	}
	return
}

func getAllAscKeyAndValues(tree Tree) (result []int) {
	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
		result = append(result, int(node.Key().(IntKey)))
		result = append(result, node.Value().(int))
		return
	})
	return
}

func checkHeight(node Node) (ok bool) {
	if node == nil {
		return true
	}
	height := node.(RealNode).Height()
	var hLeft, hRight int
	if lChild, ok := node.LeftChild().(RealNode); ok {
		hLeft = lChild.Height()
	}
	if rChild, ok := node.RightChild().(RealNode); ok {
		hRight = rChild.Height()
	}
	hMin, hMax := hLeft, hRight
	if hMax < hMin {
		hMin, hMax = hMax, hMin
	}
	return hMax-hMin <= 1 && height-hMax == 1
}

func takeInvalidNode(tree Tree, check func(node Node) bool) (invalidNode Node) {
	stack := []Node{tree.Root()}
	for len(stack) > 0 {
		newsize := len(stack) - 1
		node := stack[newsize]
		stack = stack[:newsize]
		if node == nil {
			continue
		}
		if !check(node) {
			invalidNode = node
			return
		}
		stack = append(stack, node.RightChild(), node.LeftChild())
	}
	return
}

func takeInvalidHeightNode(tree Tree) (invalidNode Node) {
	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
		if !checkHeight(node) {
			invalidNode = node
			breakIteration = true
		}
		return
	})
	if invalidNode != nil {
		return
	}
	return takeInvalidNode(tree, checkHeight)
}

func checkBalance(node Node) bool {
	if node == nil {
		return true
	}
	key := int(node.Key().(IntKey))
	if lChild := node.LeftChild(); lChild != nil {
		lKey := int(lChild.Key().(IntKey))
		if key <= lKey {
			return false
		}
	}
	if rChild := node.RightChild(); rChild != nil {
		rKey := int(rChild.Key().(IntKey))
		if rKey <= key {
			return false
		}
	}
	return true
}

func takeInvalidBalanceNode(tree Tree) (invalidNode Node) {
	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
		if !checkBalance(node) {
			invalidNode = node
			breakIteration = true
		}
		return
	})
	if invalidNode != nil {
		return
	}
	return takeInvalidNode(tree, checkBalance)
}

func TestInsertOneEntry(t *testing.T) {

	f := func(k, v int) []Tree {
		tree0 := New(false)
		tree1, _ := avltree.Insert(tree0, false, IntKey(k), v)
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

func TestInsertTwoEntries(t *testing.T) {

	f := func(k1, v1, k2, v2 int) []Tree {
		if k1 == k2 {
			return nil
		}
		tree0 := New(false)
		tree1, _ := avltree.Insert(tree0, false, IntKey(k1), v1)
		tree2, _ := avltree.Insert(tree1, false, IntKey(k2), v2)
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

func TestRejectDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree0 := New(false)
		tree1, _ := avltree.Insert(tree0, false, IntKey(k1), v1)
		tree2, _ := avltree.Insert(tree1, false, IntKey(k2), v2)
		tree3, _ := avltree.Insert(tree2, false, IntKey(k1), v3)
		tree4, _ := avltree.Insert(tree3, false, IntKey(k2), v4)
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

func TestReplaceDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree0 := New(false)
		tree1, _ := avltree.Insert(tree0, false, IntKey(k1), v1)
		tree2, _ := avltree.Insert(tree1, false, IntKey(k2), v2)
		tree3, _ := avltree.Insert(tree2, true, IntKey(k1), v3)
		tree4, _ := avltree.Insert(tree3, true, IntKey(k2), v4)
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

func TestAllowDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) []Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree0 := New(true)
		tree1, _ := avltree.Insert(tree0, false, IntKey(k1), v1)
		tree2, _ := avltree.Insert(tree1, false, IntKey(k2), v2)
		tree3, _ := avltree.Insert(tree2, false, IntKey(k1), v3)
		tree4, _ := avltree.Insert(tree3, false, IntKey(k2), v4)
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

func TestNodeCount(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := New(true)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			tree, _ = avltree.Delete(tree, IntKey(list[0].Key))
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
		tree := New(true)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			tree, _ = avltree.Delete(tree, IntKey(list[0].Key))
		}
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if node := avltree.Min(tree); node != nil {
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
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if node := avltree.Max(tree); node != nil {
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
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		return avltree.Count(tree)
	}

	g := func(listBase []keyAndValue) int {
		list := omitDuplicates(listBase)
		return len(list)
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {

	f := func(listBase []keyAndValue) *keyAndValue {
		list := omitDuplicates(listBase)
		if len(list) < 2 {
			return nil
		}
		tree := New(false)
		for _, kv := range list[1:] {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list[1:] {
			node := avltree.Find(tree, IntKey(kv.Key))
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

func TestInsertAndDelete1(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) []int {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := New(false)
		for _, list := range lists[0:2] {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			var dv avltree.KeyAndValue
			tree, dv = avltree.Delete(tree, IntKey(kv.Key))
			if dv == nil || dv.Value().(int) != kv.Value {
				panic("wrong")
			}
		}
		for _, kv := range lists[2] {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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

func TestInsertAndDelete1Height(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) Node {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := New(false)
		for _, list := range lists[0:2] {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			tree, _ = avltree.Delete(tree, IntKey(kv.Key))
		}
		for _, kv := range lists[2] {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(ins1del3, ins2, ins4 []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndDelete1Balance(t *testing.T) {

	f := func(ins1del3, ins2, ins4 []keyAndValue) Node {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1del3, ins2, ins4,
		})
		tree := New(false)
		for _, list := range lists[0:2] {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			tree, _ = avltree.Delete(tree, IntKey(kv.Key))
		}
		for _, kv := range lists[2] {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(ins1del3, ins2, ins4 []keyAndValue) Node {
		return nil
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
		tree := New(false)
		for _, list := range lists {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			tree, _ = avltree.Delete(tree, IntKey(kv.Key))
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

func TestInsertAndDelete2Height(t *testing.T) {

	f := func(ins1, ins2del4, ins3 []keyAndValue) Node {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		tree := New(false)
		for _, list := range lists {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			tree, _ = avltree.Delete(tree, IntKey(kv.Key))
		}
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(ins1, ins2del4, ins3 []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndDelete2Balance(t *testing.T) {

	f := func(ins1, ins2del4, ins3 []keyAndValue) Node {
		lists := omitAllDuplicates([][]keyAndValue{
			ins1, ins2del4, ins3,
		})
		tree := New(false)
		for _, list := range lists {
			for _, kv := range list {
				tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			tree, _ = avltree.Delete(tree, IntKey(kv.Key))
		}
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(ins1, ins2del4, ins3 []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int(nil)
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
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

func TestDescIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int(nil)
		avltree.Iterate(tree, true, func(node Node) (breakIteration bool) {
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

func TestAscHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
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

func TestDescHalfIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		result := []int{}
		avltree.Iterate(tree, true, func(node Node) (breakIteration bool) {
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

func TestAscRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		appender := func(node Node) (breakIteration bool) {
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

func TestDescRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.RangeIterate(tree, true, lower, upper, func(node Node) (breakIteration bool) {
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

func TestDescHalfRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		result := []int{}
		lower := IntKey(k1)
		upper := IntKey(k2)
		stopKey := IntKey((k2 + k1) / 2)
		avltree.RangeIterate(tree, true, lower, upper, func(node Node) (breakIteration bool) {
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

func TestDuplicateKeyAscRangeIterate(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
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
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
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
			reversed(list)
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
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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

func TestDescRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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

func TestDuplicateKeyAscRange(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
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
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
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
			reversed(list)
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
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
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

func TestDeleteAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		beforeTree := tree
		afterTree, delValues := avltree.DeleteAll(tree, IntKey(delkey%keymax))
		result := make([][]int, keymax)
		avltree.Iterate(afterTree, false, func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			result[key] = append(result[key], node.Value().(int))
			return
		})
		result = append(result, toKeyValueInts(delValues))
		beforeValues := []int(nil)
		avltree.Iterate(beforeTree, false, func(node Node) (breakIteration bool) {
			key := int(node.Key().(IntKey))
			beforeValues = append(beforeValues, key, node.Value().(int))
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

func TestDeleteAllHeight(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) Node {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		tree, _ = avltree.DeleteAll(tree, IntKey(delkey%keymax))
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(list []keyAndValue, delkey int) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestFindAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := [][]int{}
		for key := 0; key < keymax; key++ {
			values := []int(nil)
			nodes := avltree.FindAll(tree, IntKey(key))
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

func TestMinAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		nodes := avltree.MinAll(tree)
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

func TestMaxAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue) []int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		result := []int(nil)
		nodes := avltree.MaxAll(tree)
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

func TestUpdateValueByFind(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		// Immutableなのでこの方法では更新できない
		for _, kv := range list {
			node := avltree.Find(tree, IntKey(kv.Key))
			value := node.Value().(int)
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

func TestUpdateValueByIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		// Immutableなのでこの方法では更新できない
		avltree.Iterate(tree, true, func(node Node) (breakIteration bool) {
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
		result := toKeyValueInts(list)
		return result
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree, _ = avltree.Update(tree, IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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

func TestKeepOldValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree, _ = avltree.Update(tree, IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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

func TestReplaceValue(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree, _ = avltree.Replace(tree, IntKey(kv.Key), kv.Value^value)
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

func TestAscUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.UpdateIterate(tree, false, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestAscHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		tree, _ = avltree.UpdateIterate(tree, false, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestDescUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.UpdateIterate(tree, true, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestDescHalfUpdateIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, true, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		tree, _ = avltree.UpdateIterate(tree, true, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestAscUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestAscHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestDescUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRangeIterate(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestDescHalfUpdateRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRangeIterate(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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

func TestAscUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRange(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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

func TestDescUpdateRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.UpdateRange(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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

func TestReplaceRange(t *testing.T) {

	const value int = 123456

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.ReplaceRange(tree, lower, upper, value)
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

func TestUpdateAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		count := 0
		tree, _ = avltree.UpdateAll(tree, IntKey(updkey%keymax), func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
			value := oldValue.(int)
			newValue = value ^ updkey
			count++
			keepOldValue = count%2 == 0
			return
		})
		result := make([][]int, keymax)
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
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

func TestReplaceAll(t *testing.T) {

	const keymax = 4
	const value int = 123456

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		tree, _ = avltree.ReplaceAll(tree, IntKey(updkey%keymax), value)
		result := make([][]int, keymax)
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
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

func TestAscDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		afterTree, values := avltree.DeleteIterate(tree, false, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		afterTree, values := avltree.DeleteIterate(tree, false, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		afterTree, values := avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		afterTree, values := avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscDeleteIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.DeleteIterate(tree, false, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
			value := oldValue.(int)
			deleteNode = value < 0
			return
		})
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestDescDeleteIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestAscDeleteIterateBalance(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.DeleteIterate(tree, false, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
			value := oldValue.(int)
			deleteNode = value < 0
			return
		})
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestDescDeleteIterateBalance(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _ = avltree.DeleteIterate(tree, true, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
			value := oldValue.(int)
			deleteNode = value < 0
			return
		})
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestAscDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			count--
			breakIteration = count <= 0
			return
		})
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRange(tree, false, lower, upper)
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values := avltree.DeleteRange(tree, true, lower, upper)
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscDeleteRangeIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) Node {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue, k1, k2 int) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescDeleteRangeIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) Node {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _ = avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
			deleteNode = value.(int) < 0
			return
		})
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue, k1, k2 int) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAlter(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			var delValue avltree.KeyAndValue
			tree, delValue, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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

func TestAlterHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree, _, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
		}
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAlterBalance(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			tree, _, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
		}
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		afterTree, delValues, _ := avltree.AlterIterate(tree, false, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		afterTree, delValues, _ := avltree.AlterIterate(tree, false, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		afterTree, delValues, _ := avltree.AlterIterate(tree, true, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescHalfAlterIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		afterTree, delValues, _ := avltree.AlterIterate(tree, true, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscAlterIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _, _ = avltree.AlterIterate(tree, false, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscAlterIterateBalance(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _, _ = avltree.AlterIterate(tree, false, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescAlterIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _, _ = avltree.AlterIterate(tree, true, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescAlterIterateBalance(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		tree, _, _ = avltree.AlterIterate(tree, true, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidBalanceNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAscAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRangeIterate(tree, false, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRangeIterate(tree, false, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRangeIterate(tree, true, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescHalfAlterRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRangeIterate(tree, true, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRange(tree, false, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestDescAlterRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		afterTree, values, _ := avltree.AlterRange(tree, true, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
		result := getAllAscKeyAndValues(afterTree)
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

func TestAscAlterRangeIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) Node {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _, _ = avltree.AlterRangeIterate(tree, false, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue, k1, k2 int) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescAlterRangeIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) Node {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		tree, _, _ = avltree.AlterRangeIterate(tree, true, lower, upper, func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
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
		invalidNode := takeInvalidHeightNode(tree)
		return invalidNode
	}

	g := func(listBase []keyAndValue, k1, k2 int) Node {
		return nil
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAlterAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, updkey int) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			tree, _ = avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		afterTree, values, _ := avltree.AlterAll(tree, IntKey(updkey%keymax), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
		avltree.Iterate(afterTree, false, func(node Node) (breakIteration bool) {
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
