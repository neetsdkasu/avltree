package linkedtree

import (
	"sort"
	"testing"
	"testing/quick"

	"avltree"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

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

func ascIterateNode(node Node, callBack func(n Node)) {
	if node != nil {
		ascIterateNode(node.LeftChild(), callBack)
		callBack(node)
		ascIterateNode(node.RightChild(), callBack)
	}
}

func getAllAscKeyAndValues(tree Tree) (result []int) {
	// avltree.Iterateに頼らない方法のほうがいいかも？
	//	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
	//		result = append(result, int(node.Key().(IntKey)))
	//		result = append(result, node.Value().(int))
	//		return
	//	})
	ascIterateNode(tree.Root(), func(node Node) {
		result = append(result, int(node.Key().(IntKey)))
		result = append(result, node.Value().(int))
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
	// avltree.Iterateに頼らない方法のほうがいいかも？
	//	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
	//		if !checkHeight(node) {
	//			invalidNode = node
	//			breakIteration = true
	//		}
	//		return
	//	})
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
	// avltree.Iterateに頼らない方法のほうがいいかも？
	//	avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
	//		if !checkBalance(node) {
	//			invalidNode = node
	//			breakIteration = true
	//		}
	//		return
	//	})
	return takeInvalidNode(tree, checkBalance)
}

func TestInsertOneEntry(t *testing.T) {

	f := func(k, v int) Tree {
		tree := New(false)
		avltree.Insert(tree, false, IntKey(k), v)
		return tree
	}

	g := func(k, v int) Tree {
		root := &LinkedTreeNode{nil, nil, 1, nil, 1, IntKey(k), v}
		tree := &LinkedTree{root, false}
		return tree
	}

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInsertTwoEntries(t *testing.T) {

	f := func(k1, v1, k2, v2 int) Tree {
		if k1 == k2 {
			return nil
		}
		tree := New(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		return tree
	}

	g := func(k1, v1, k2, v2 int) Tree {
		if k1 == k2 {
			return nil
		}
		root := &LinkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v1}
		child := &LinkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.LeftChildNode = child
		} else {
			root.RightChildNode = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestRejectDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, false, IntKey(k1), v3)
		avltree.Insert(tree, false, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &LinkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v1}
		child := &LinkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v2}
		if k2 < k1 {
			root.LeftChildNode = child
		} else {
			root.RightChildNode = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestReplaceDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(false)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, true, IntKey(k1), v3)
		avltree.Insert(tree, true, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		root := &LinkedTreeNode{nil, nil, 2, nil, 2, IntKey(k1), v3}
		child := &LinkedTreeNode{nil, nil, 1, root, 1, IntKey(k2), v4}
		if k2 < k1 {
			root.LeftChildNode = child
		} else {
			root.RightChildNode = child
		}
		tree := &LinkedTree{root, false}
		return tree
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestAllowDuplicateKey(t *testing.T) {

	f := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		tree := New(true)
		avltree.Insert(tree, false, IntKey(k1), v1)
		avltree.Insert(tree, false, IntKey(k2), v2)
		avltree.Insert(tree, false, IntKey(k1), v3)
		avltree.Insert(tree, false, IntKey(k2), v4)
		return tree
	}

	g := func(k1, v1, k2, v2, v3, v4 int) Tree {
		if k1 == k2 || v1 == v3 || v2 == v4 {
			return nil
		}
		var root *LinkedTreeNode
		if k2 < k1 {
			root = &LinkedTreeNode{nil, nil, 3, nil, 4, IntKey(k1), v1}
			lChild := &LinkedTreeNode{nil, nil, 2, root, 2, IntKey(k2), v2}
			rChild := &LinkedTreeNode{nil, nil, 1, root, 1, IntKey(k1), v3}
			lrChild := &LinkedTreeNode{nil, nil, 1, lChild, 1, IntKey(k2), v4}
			lChild.RightChildNode = lrChild
			root.LeftChildNode = lChild
			root.RightChildNode = rChild
		} else {
			root = &LinkedTreeNode{nil, nil, 3, nil, 4, IntKey(k1), v3}
			lChild := &LinkedTreeNode{nil, nil, 1, root, 1, IntKey(k1), v1}
			rChild := &LinkedTreeNode{nil, nil, 2, root, 2, IntKey(k2), v2}
			rrChild := &LinkedTreeNode{nil, nil, 1, rChild, 1, IntKey(k2), v4}
			rChild.RightChildNode = rrChild
			root.LeftChildNode = lChild
			root.RightChildNode = rChild
		}
		tree := &LinkedTree{root, true}
		return tree
	}

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestNodeCount(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := New(true)
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
		tree := New(true)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		if len(list) > 1 {
			avltree.Delete(tree, IntKey(list[0].Key))
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

func TestParent(t *testing.T) {

	f := func(list []keyAndValue) Node {
		tree := New(true)
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

	if err := quick.CheckEqual(f, g, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestMin(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
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
		tree := New(false)
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
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
		tree := New(false)
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
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			avltree.Delete(tree, IntKey(kv.Key))
		}
		for _, kv := range lists[2] {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[0] {
			avltree.Delete(tree, IntKey(kv.Key))
		}
		for _, kv := range lists[2] {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			avltree.Delete(tree, IntKey(kv.Key))
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
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			avltree.Delete(tree, IntKey(kv.Key))
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
				avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
			}
		}
		for _, kv := range lists[1] {
			avltree.Delete(tree, IntKey(kv.Key))
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
		tree := New(true)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
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
		tree := New(true)
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
		tree := New(true)
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

func TestDeleteAll(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) [][]int {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		avltree.DeleteAll(tree, IntKey(delkey%keymax))
		result := make([][]int, keymax)
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
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

func TestDeleteAllHeight(t *testing.T) {

	const keymax = 4

	f := func(list []keyAndValue, delkey int) Node {
		tree := New(true)
		for _, kv := range list {
			key := kv.Key
			if key < 0 {
				key ^= -1
			}
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if delkey < 0 {
			delkey ^= -1
		}
		avltree.DeleteAll(tree, IntKey(delkey%keymax))
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
		tree := New(true)
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
		tree := New(true)
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

func TestUpdateValueByFind(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			node, _ := avltree.Find(tree, IntKey(kv.Key))
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

func TestUpdateValueByIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
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

func TestUpdateValue(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			avltree.Update(tree, IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			avltree.Update(tree, IntKey(kv.Key), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		for _, kv := range list {
			avltree.Replace(tree, IntKey(kv.Key), kv.Value^value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.UpdateIterate(tree, false, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		avltree.UpdateIterate(tree, false, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.UpdateIterate(tree, true, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, true, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		avltree.UpdateIterate(tree, true, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRangeIterate(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRangeIterate(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRange(tree, false, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.UpdateRange(tree, true, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.ReplaceRange(tree, lower, upper, value)
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
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		count := 0
		avltree.UpdateAll(tree, IntKey(updkey%keymax), func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
			avltree.Insert(tree, false, IntKey(key%keymax), kv.Value)
		}
		if updkey < 0 {
			updkey ^= -1
		}
		avltree.ReplaceAll(tree, IntKey(updkey%keymax), value)
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		_, values := avltree.DeleteIterate(tree, false, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestAscHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		_, values := avltree.DeleteIterate(tree, false, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestDescDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		_, values := avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestDescHalfDeleteIterate(t *testing.T) {

	f := func(listBase []keyAndValue) []int {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		_, values := avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestAscDeleteIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue) Node {
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.DeleteIterate(tree, false, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.DeleteIterate(tree, true, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.DeleteIterate(tree, false, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		avltree.DeleteIterate(tree, true, func(key Key, oldValue interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestAscHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestDescDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestDescHalfDeleteRangeIterate(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		count := (len(list) + 1) / 2
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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

func TestAscDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRange(tree, false, lower, upper)
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

func TestDescDeleteRange(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) []int {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		_, values := avltree.DeleteRange(tree, true, lower, upper)
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

func TestAscDeleteRangeIterateHeight(t *testing.T) {

	f := func(listBase []keyAndValue, k1, k2 int) Node {
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		list := omitDuplicates(listBase)
		tree := New(false)
		for _, kv := range list {
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.DeleteRangeIterate(tree, false, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		lower := IntKey(k1)
		upper := IntKey(k2)
		avltree.DeleteRangeIterate(tree, true, lower, upper, func(key Key, value interface{}) (deleteNode, breakIteration bool) {
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			_, delvalue, _ := avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
			if delvalue != nil {
				values = append(values, delvalue.(int))
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			_, delvalue, _ := avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
			if delvalue != nil {
				values = append(values, delvalue.(int))
			}
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
			avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
		values := []int{}
		for _, kv := range list {
			_, delvalue, _ := avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
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
			if delvalue != nil {
				values = append(values, delvalue.(int))
			}
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
