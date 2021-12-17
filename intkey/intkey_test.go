// Author: Leonardone @ NEETSDKASU
// License: MIT

package intkey

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/simpletree"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

func TestIntKey(t *testing.T) {
	f := func(k1, k2 int) bool {
		var key1 avltree.Key = IntKey(k1)
		var key2 avltree.Key = IntKey(k2)
		switch key1.CompareTo(key2) {
		case avltree.LessThanOtherKey:
			return k1 < k2
		case avltree.EqualToOtherKey:
			return k1 == k2
		case avltree.GreaterThanOtherKey:
			return k1 > k2
		default:
			return false
		}
	}

	if err := quick.Check(f, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func Example() {
	tree := simpletree.New(false)
	avltree.Insert(tree, false, IntKey(12), 345)
	avltree.Insert(tree, false, IntKey(67), 890)
	avltree.Insert(tree, false, IntKey(333), 666)
	avltree.Insert(tree, false, IntKey(-5), 12345)
	avltree.Delete(tree, IntKey(67))
	avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
		newValue = oldValue.(int) * 3
		return
	})
	if node := avltree.Find(tree, IntKey(12)); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
}
