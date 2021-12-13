package immutabletree

import (
	"math/rand"
	"testing"

	"github.com/neetsdkasu/avltree"
)

func genKeyAndValues(n int) []*keyAndValue {
	list := []*keyAndValue{}
	for i := 0; i < n; i++ {
		key := rand.Int()
		value := rand.Int()
		list = append(list, &keyAndValue{key, value})
	}
	return list
}

func BenchmarkInsert(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list[:b.N] {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[b.N:]
	b.ResetTimer()
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
}

func BenchmarkDelete(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	b.ResetTimer()
	for _, kv := range list {
		tree, _, _ = avltree.Delete(tree, IntKey(kv.Key))
	}
}

func BenchmarkUpdate(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	b.ResetTimer()
	for _, kv := range list {
		tree, _ = avltree.Update(tree, IntKey(kv.Key), func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
			newValue = oldValue.(int) >> 1
			return
		})
	}
}

func BenchmarkReplace(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	const value int = 12345
	b.ResetTimer()
	for _, kv := range list {
		tree, _ = avltree.Replace(tree, IntKey(kv.Key), value)
	}
}

func BenchmarkDeleteByAlter(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	b.ResetTimer()
	for _, kv := range list {
		tree, _, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
			request.Delete()
			return
		})
	}
}

func BenchmarkUpdateByAlter(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	b.ResetTimer()
	for _, kv := range list {
		tree, _, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
			request.Replace(node.Value().(int) >> 1)
			return
		})
	}
}

func BenchmarkReplaceByAlter(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	const value int = 12345
	b.ResetTimer()
	for _, kv := range list {
		tree, _, _ = avltree.Alter(tree, IntKey(kv.Key), func(node avltree.AlterNode) (request avltree.AlterRequest) {
			request.Replace(value)
			return
		})
	}
}

func BenchmarkFind1(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(2 * b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	list = list[:b.N]
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		for _, kv := range list {
			avltree.Find(tree, IntKey(kv.Key))
		}
	}
}

func BenchmarkFind2(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(4 * b.N)
	for i, kv := range list {
		if (i & 1) == 0 {
			tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
		}
	}
	list = list[:b.N]
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		for _, kv := range list {
			avltree.Find(tree, IntKey(kv.Key))
		}
	}
}

func BenchmarkAscIterate(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		avltree.Iterate(tree, false, func(node Node) (breakIteration bool) {
			return
		})
	}
}

func BenchmarkDescIterate(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		avltree.Iterate(tree, true, func(node Node) (breakIteration bool) {
			return
		})
	}
}

func BenchmarkAscRangeIterate(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		k1, k2 := rand.Int(), rand.Int()
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		lower, upper := IntKey(k1), IntKey(k2)
		avltree.RangeIterate(tree, false, lower, upper, func(node Node) (breakIteration bool) {
			return
		})
	}
}

func BenchmarkDescRangeIterate(b *testing.B) {
	tree := New(true)
	list := genKeyAndValues(b.N)
	for _, kv := range list {
		tree, _ = avltree.Insert(tree, false, IntKey(kv.Key), kv.Value)
	}
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		k1, k2 := rand.Int(), rand.Int()
		if k2 < k1 {
			k1, k2 = k2, k1
		}
		lower, upper := IntKey(k1), IntKey(k2)
		avltree.RangeIterate(tree, true, lower, upper, func(node Node) (breakIteration bool) {
			return
		})
	}
}
