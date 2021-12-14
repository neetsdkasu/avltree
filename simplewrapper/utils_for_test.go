package simplewrapper

import (
	"sort"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/intkey"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

type (
	IntKey = intkey.IntKey
	Key    = avltree.Key
	Node   = avltree.Node
	Tree   = avltree.Tree
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

func getAllAscKeyAndValues(tree *AVLTree) (result []int) {
	tree.Iterate(func(node Node) (breakIteration bool) {
		result = append(result, int(node.Key().(IntKey)))
		result = append(result, node.Value().(int))
		return
	})
	return
}
