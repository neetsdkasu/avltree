// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのRealTree,RealNodeの実装例
// simpletreeと比較した場合standardtreeは以下の機能を持っている
// + 各ノードが親ノードの情報を持っている (avltree.ParentGetterを実装している)
// + 各ノードがそのノードをサブツリーとしたときのノード総数の情報を持っている(avltree.NodeCounterを実装している)
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			. "github.com/neetsdkasu/avltree/intkey"
//			"github.com/neetsdkasu/avltree/standardtree"
//		)
//		func Example_standardtree() {
//			tree := standardtree.New(false)
//			avltree.Insert(tree, false, IntKey(12), 345)
//			avltree.Insert(tree, false, IntKey(67), 890)
//			avltree.Insert(tree, false, IntKey(333), 666)
//			avltree.Insert(tree, false, IntKey(-5), 12345)
//			avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Iterate!", node.Key(), node.Value())
//				if counter, ok := node.(avltree.NodeCounter); ok {
//					fmt.Println("  Node Count", counter.NodeCount())
//				}
//				if getter, ok := node.(avltree.ParentGetter); ok {
//					parent := getter.Parent()
//					if parent != nil {
//						fmt.Println("  Parent is", parent.Key())
//					} else {
//						fmt.Println("  No Parent")
//					}
//				}
//				return
//			})
//			// Output:
//			// Iterate! -5 12345
//			//   Node Count 1
//			//   Parent is 12
//			// Iterate! 12 345
//			//   Node Count 2
//			//   Parent is 67
//			// Iterate! 67 890
//			//   Node Count 4
//			//   No Parent
//			// Iterate! 333 666
//			//   Node Count 1
//			//   Parent is 67
//		}
//
package standardtree

import "github.com/neetsdkasu/avltree"

type StandardTree struct {
	RootNode                *StandardTreeNode
	AllowDuplicateKeysValue bool
}

type StandardTreeNode struct {
	LeftChildNode  *StandardTreeNode
	RightChildNode *StandardTreeNode
	HeightValue    int
	ParentNode     *StandardTreeNode
	NodeCountValue int
	KeyData        avltree.Key
	ValueData      interface{}
}

func unwrap(node avltree.Node) *StandardTreeNode {
	if ltNode, ok := node.(*StandardTreeNode); ok {
		return ltNode
	} else {
		return nil
	}
}

func (node *StandardTreeNode) toNode() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func New(allowDuplicateKeys bool) avltree.Tree {
	return &StandardTree{
		nil, // RootNode
		allowDuplicateKeys,
	}
}

func (tree *StandardTree) NodeCount() int {
	return tree.RootNode.NodeCount()
}

func (tree *StandardTree) NewNode(
	leftChild,
	rightChild avltree.Node,
	height int,
	key avltree.Key,
	value interface{},
) avltree.RealNode {
	node := &StandardTreeNode{
		unwrap(leftChild),
		unwrap(rightChild),
		height,
		nil, // ParentNode
		1,   // NodeCountValue
		key,
		value,
	}
	node.resetNodeCount()
	return node
}

func (tree *StandardTree) Root() avltree.Node {
	return tree.RootNode.toNode()
}

func (tree *StandardTree) SetRoot(newRoot avltree.RealNode) avltree.RealTree {
	tree.RootNode = unwrap(newRoot)
	tree.RootNode.setParent(nil)
	return tree
}

func (tree *StandardTree) AllowDuplicateKeys() bool {
	return tree.AllowDuplicateKeysValue
}

func (node *StandardTreeNode) Key() avltree.Key {
	return node.KeyData
}

func (node *StandardTreeNode) Value() interface{} {
	return node.ValueData
}

func (node *StandardTreeNode) Height() int {
	return node.HeightValue
}

func (node *StandardTreeNode) LeftChild() avltree.Node {
	return node.LeftChildNode.toNode()
}

func (node *StandardTreeNode) RightChild() avltree.Node {
	return node.RightChildNode.toNode()
}

func (node *StandardTreeNode) SetValue(newValue interface{}) avltree.Node {
	node.ValueData = newValue
	return node
}

func (node *StandardTreeNode) Parent() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node.ParentNode.toNode()
	}
}

func (node *StandardTreeNode) setParent(newParent avltree.Node) {
	if node != nil {
		node.ParentNode = unwrap(newParent)
	}
}

func (node *StandardTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.NodeCountValue
	}
}

func (node *StandardTreeNode) resetNodeCount() {
	node.NodeCountValue = 1 +
		node.LeftChildNode.NodeCount() +
		node.RightChildNode.NodeCount()
}

func (node *StandardTreeNode) SetChildren(
	newLeftChild,
	newRightChild avltree.Node,
	newHeight int,
) avltree.RealNode {
	node.LeftChildNode = unwrap(newLeftChild)
	node.RightChildNode = unwrap(newRightChild)
	node.HeightValue = newHeight
	node.LeftChildNode.setParent(node)
	node.RightChildNode.setParent(node)
	node.resetNodeCount()
	return node
}

func (node *StandardTreeNode) Set(
	newLeftChild,
	newRightChild avltree.Node,
	newHeight int,
	newValue interface{},
) avltree.RealNode {
	node.ValueData = newValue
	return node.SetChildren(newLeftChild, newRightChild, newHeight)
}
