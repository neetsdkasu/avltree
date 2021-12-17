// Author: Leonardone @ NEETSDKASU
// License: MIT

package examples

import (
	"fmt"

	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/immutabletree"
	. "github.com/neetsdkasu/avltree/intkey"
	"github.com/neetsdkasu/avltree/simplewrapper"
)

func Example_immutabletree() {
	tree := immutabletree.New(false)
	tree, _ = avltree.Insert(tree, false, IntKey(12), 345)
	tree, _ = avltree.Insert(tree, false, IntKey(67), 890)
	tree, _ = avltree.Insert(tree, false, IntKey(333), 666)
	tree, _ = avltree.Insert(tree, false, IntKey(-5), 12345)
	saved := tree
	tree, _ = avltree.Delete(tree, IntKey(67))
	tree, _ = avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
	avltree.Iterate(saved, false, func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Saved!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
	// Saved! -5 12345
	// Saved! 12 345
	// Saved! 67 890
	// Saved! 333 666
}

func Example_immutabletreeWithWrapper() {
	tree := immutabletree.New(false)
	w := simplewrapper.New(tree)
	w.Insert(IntKey(12), 345)
	w.Insert(IntKey(67), 890)
	w.Insert(IntKey(333), 666)
	w.Insert(IntKey(-5), 12345)
	saved := *w
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
	saved.Iterate(func(node avltree.Node) (breakIteration bool) {
		fmt.Println("Saved!", node.Key(), node.Value())
		return
	})
	// Output:
	// Find! 12 345
	// Iterate! -5 12345
	// Iterate! 12 345
	// Iterate! 333 1998
	// Saved! -5 12345
	// Saved! 12 345
	// Saved! 67 890
	// Saved! 333 666
}
