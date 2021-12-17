// Author: Leonardone @ NEETSDKASU
// License: MIT

package examples

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/simpletree"
)

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
