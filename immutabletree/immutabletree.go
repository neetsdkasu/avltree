// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのRealTree,RealNodeの実装例
// 不変(immutable)ぽい木を構築する
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			"github.com/neetsdkasu/avltree/immutabletree"
//			. "github.com/neetsdkasu/avltree/intkey"
//		)
//		func Example_immutabletree() {
//			tree := immutabletree.New(false)
//			tree, _ = avltree.Insert(tree, false, IntKey(12), 345)
//			tree, _ = avltree.Insert(tree, false, IntKey(67), 890)
//			tree, _ = avltree.Insert(tree, false, IntKey(333), 666)
//			tree, _ = avltree.Insert(tree, false, IntKey(-5), 12345)
//			saved := tree
//			tree, _ = avltree.Delete(tree, IntKey(67))
//			tree, _ = avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
//				newValue = oldValue.(int) * 3
//				return
//			})
//			if node := avltree.Find(tree, IntKey(12)); node != nil {
//				fmt.Println("Find!", node.Key(), node.Value())
//			}
//			avltree.Iterate(tree, false, func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Iterate!", node.Key(), node.Value())
//				return
//			})
//			avltree.Iterate(saved, false, func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Saved!", node.Key(), node.Value())
//				return
//			})
//			// Output:
//			// Find! 12 345
//			// Iterate! -5 12345
//			// Iterate! 12 345
//			// Iterate! 333 1998
//			// Saved! -5 12345
//			// Saved! 12 345
//			// Saved! 67 890
//			// Saved! 333 666
//		}
//
//
// ラッパーを用いたコード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			"github.com/neetsdkasu/avltree/immutabletree"
//			. "github.com/neetsdkasu/avltree/intkey"
//			"github.com/neetsdkasu/avltree/simplewrapper"
//		)
//		func Example_immutabletreeWithWrapper() {
//			tree := immutabletree.New(false)
//			w := simplewrapper.New(tree)
//			w.Insert(IntKey(12), 345)
//			w.Insert(IntKey(67), 890)
//			w.Insert(IntKey(333), 666)
//			w.Insert(IntKey(-5), 12345)
//			saved := *w
//			w.Delete(IntKey(67))
//			w.Update(IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
//				newValue = oldValue.(int) * 3
//				return
//			})
//			if node := w.Find(IntKey(12)); node != nil {
//				fmt.Println("Find!", node.Key(), node.Value())
//			}
//			w.Iterate(func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Iterate!", node.Key(), node.Value())
//				return
//			})
//			saved.Iterate(func(node avltree.Node) (breakIteration bool) {
//				fmt.Println("Saved!", node.Key(), node.Value())
//				return
//			})
//			// Output:
//			// Find! 12 345
//			// Iterate! -5 12345
//			// Iterate! 12 345
//			// Iterate! 333 1998
//			// Saved! -5 12345
//			// Saved! 12 345
//			// Saved! 67 890
//			// Saved! 333 666
//		}
//
package immutabletree

import "github.com/neetsdkasu/avltree"

type ImmutableTree struct {
	RootNode                *ImmutableTreeNode
	AllowDuplicateKeysValue bool
}

type ImmutableTreeNode struct {
	LeftChildNode  *ImmutableTreeNode
	RightChildNode *ImmutableTreeNode
	HeightValue    int
	NodeCountValue int
	KeyData        avltree.Key
	ValueData      interface{}
}

func New(allowDuplicateKeys bool) avltree.Tree {
	return &ImmutableTree{
		RootNode:                nil,
		AllowDuplicateKeysValue: allowDuplicateKeys,
	}
}

func unwrap(node avltree.Node) *ImmutableTreeNode {
	if node == nil {
		return nil
	} else {
		return node.(*ImmutableTreeNode)
	}
}

func (node *ImmutableTreeNode) toNode() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func (tree *ImmutableTree) Root() avltree.Node {
	return tree.RootNode.toNode()
}

func (tree *ImmutableTree) NewNode(leftChild, rightChild avltree.Node, height int, key avltree.Key, value interface{}) avltree.RealNode {
	newNode := &ImmutableTreeNode{
		LeftChildNode:  unwrap(leftChild),
		RightChildNode: unwrap(rightChild),
		HeightValue:    height,
		NodeCountValue: 1,
		KeyData:        key,
		ValueData:      value,
	}
	newNode.resetNodeCount()
	return newNode
}

func (tree *ImmutableTree) SetRoot(newRoot avltree.RealNode) avltree.RealTree {
	newTree := *tree
	newTree.RootNode = unwrap(newRoot)
	return &newTree
}

func (tree *ImmutableTree) AllowDuplicateKeys() bool {
	return tree.AllowDuplicateKeysValue
}

func (tree *ImmutableTree) NodeCount() int {
	return tree.RootNode.NodeCount()
}

func (node *ImmutableTreeNode) resetNodeCount() {
	if node != nil {
		node.NodeCountValue = 1 + node.LeftChildNode.NodeCount() + node.RightChildNode.NodeCount()
	}
}

func (node *ImmutableTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.NodeCountValue
	}
}

func (node *ImmutableTreeNode) Key() avltree.Key {
	return node.KeyData
}

func (node *ImmutableTreeNode) Value() interface{} {
	return node.ValueData
}

func (node *ImmutableTreeNode) LeftChild() avltree.Node {
	return node.LeftChildNode.toNode()
}

func (node *ImmutableTreeNode) RightChild() avltree.Node {
	return node.RightChildNode.toNode()
}

func (node *ImmutableTreeNode) SetValue(newValue interface{}) avltree.Node {
	newNode := *node
	newNode.ValueData = newValue
	return &newNode
}

func (node *ImmutableTreeNode) Height() int {
	return node.HeightValue
}

func (node *ImmutableTreeNode) SetChildren(newLeftChild, newRightChild avltree.Node, newHeight int) avltree.RealNode {
	newNode := *node
	newNode.LeftChildNode = unwrap(newLeftChild)
	newNode.RightChildNode = unwrap(newRightChild)
	newNode.HeightValue = newHeight
	newNode.resetNodeCount()
	return &newNode
}

func (node *ImmutableTreeNode) Set(newLeftChild, newRightChild avltree.Node, newHeight int, newValue interface{}) avltree.RealNode {
	newNode := *node
	newNode.LeftChildNode = unwrap(newLeftChild)
	newNode.RightChildNode = unwrap(newRightChild)
	newNode.HeightValue = newHeight
	newNode.ValueData = newValue
	newNode.resetNodeCount()
	return &newNode
}
