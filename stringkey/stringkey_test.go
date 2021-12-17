// Author: Leonardone @ NEETSDKASU
// License: MIT

package stringkey

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/simpletree"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

func TestStringKey(t *testing.T) {
	f := func(k1, k2 string) bool {
		var key1 avltree.Key = StringKey(k1)
		var key2 avltree.Key = StringKey(k2)
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
	avltree.Insert(tree, false, StringKey("ABC"), 1234)
	avltree.Insert(tree, false, StringKey("XYZ"), 9999)
	avltree.Insert(tree, false, StringKey("FOOBAR"), -500)
	avltree.Insert(tree, false, StringKey("FizzBuzz"), 3515)
	avltree.Delete(tree, StringKey("ABC"))
	avltree.Update(tree, StringKey("FOOBAR"), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
		newValue = oldValue.(int) * 4
		return
	})
	if node := avltree.Find(tree, StringKey("XYZ")); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! XYZ 9999
	// Iterate! FOOBAR -2000
	// Iterate! FizzBuzz 3515
	// Iterate! XYZ 9999
}
