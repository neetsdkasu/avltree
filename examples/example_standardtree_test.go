// Author: Leonardone @ NEETSDKASU
// License: MIT

package examples

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/standardtree"
)

func Example_standardtree() {
	tree := standardtree.New(false)
	avltree.Insert(tree, false, IntKey(12), 345)
	avltree.Insert(tree, false, IntKey(67), 890)
	avltree.Insert(tree, false, IntKey(333), 666)
	avltree.Insert(tree, false, IntKey(-5), 12345)
	avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Iterate!", node.Key(), node.Value())
		if counter, ok := node.(avltree.NodeCounter); ok {
			fmt.Println("  Node Count", counter.NodeCount())
		}
		if getter, ok := node.(avltree.ParentGetter); ok {
			parent := getter.Parent()
			if parent != nil {
				fmt.Println("  Parent is", parent.Key())
			} else {
				fmt.Println("  No Parent")
			}
		}
		return
	})
	// Output:
	// Iterate! -5 12345
	//   Node Count 1
	//   Parent is 12
	// Iterate! 12 345
	//   Node Count 2
	//   Parent is 67
	// Iterate! 67 890
	//   Node Count 4
	//   No Parent
	// Iterate! 333 666
	//   Node Count 1
	//   Parent is 67
}
