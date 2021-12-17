// Author: Leonardone @ NEETSDKASU
// License: MIT

package examples

import (
	"fmt"

	"github.com/neetsdkasu/avltree/intwrapper"
	"github.com/neetsdkasu/avltree/simpletree"
)

func Example_intwrapper() {
	tree := simpletree.New(false)
	w := intwrapper.New(tree)
	w.Insert(12, 345)
	w.Insert(67, 890)
	w.Insert(333, 666)
	w.Insert(-5, 12345)
	w.Delete(67)
	w.Update(333, func(key, oldValue int) (newValue int, keepOldValue bool) {
		newValue = oldValue * 3
		return
	})
	if node := w.Find(12); node != nil {
		fmt.Println("Find!", node.Key(), node.Value())
	}
	w.Iterate(func(node intwrapper.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
}
