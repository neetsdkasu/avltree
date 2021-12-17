// Author: Leonardone @ NEETSDKASU
// License: MIT

package examples

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/simpletree"
	"github.com/neetsdkasu/avltree/simplewrapper"
)

func Example_wrapper() {
	tree := simpletree.New(false)
	w := simplewrapper.New(tree)
	w.Insert(IntKey(12), 345)
	w.Insert(IntKey(67), 890)
	w.Insert(IntKey(333), 666)
	w.Insert(IntKey(-5), 12345)
	w.Delete(IntKey(67))
	w.Update(IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
		newValue = oldValue.(int) * 3
		return
	})
	if node := w.Find(IntKey(12)); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	w.Iterate(func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
}
