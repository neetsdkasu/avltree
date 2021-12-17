// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのRealTree,RealNodeの実装例
// 実装に必要な最小限の構成になっている
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree"
//			. "github.com/neetsdkasu/avltree/intkey"
//			"github.com/neetsdkasu/avltree/simpletree"
//		)
//		func Example() {
//			tree := simpletree.New(false)
//			avltree.Insert(tree, false, IntKey(12), 345)
//			avltree.Insert(tree, false, IntKey(67), 890)
//			avltree.Insert(tree, false, IntKey(333), 666)
//			avltree.Insert(tree, false, IntKey(-5), 12345)
//			avltree.Delete(tree, IntKey(67))
//			avltree.Update(tree, IntKey(333), func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue bool) {
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
//			// Output:
//			// Find! 12 345
//			// Iterate! -5 12345
//			// Iterate! 12 345
//			// Iterate! 333 1998
//		}
//
package simpletree

import "github.com/neetsdkasu/avltree"

type SimpleTree struct {
	RootNode                avltree.Node
	AllowDuplicateKeysValue bool
}

type SimpleNode struct {
	LeftChildNode  avltree.Node
	RightChildNode avltree.Node
	HeightValue    int
	KeyData        avltree.Key
	ValueData      interface{}
}

func New(allowDuplicateKeys bool) avltree.Tree {
	return &SimpleTree{nil, allowDuplicateKeys}
}

func (tree *SimpleTree) Root() avltree.Node {
	return tree.RootNode
}

func (tree *SimpleTree) AllowDuplicateKeys() bool {
	return tree.AllowDuplicateKeysValue
}

func (tree *SimpleTree) NewNode(leftChild, rightChild avltree.Node, height int, key avltree.Key, value interface{}) avltree.RealNode {
	return &SimpleNode{leftChild, rightChild, height, key, value}
}

func (tree *SimpleTree) SetRoot(newRoot avltree.RealNode) avltree.RealTree {
	tree.RootNode = newRoot
	return tree
}

func (node *SimpleNode) Key() avltree.Key {
	return node.KeyData
}

func (node *SimpleNode) Value() interface{} {
	return node.ValueData
}

func (node *SimpleNode) LeftChild() avltree.Node {
	return node.LeftChildNode
}

func (node *SimpleNode) RightChild() avltree.Node {
	return node.RightChildNode
}

func (node *SimpleNode) SetValue(newValue interface{}) avltree.Node {
	node.ValueData = newValue
	return node
}

func (node *SimpleNode) Height() int {
	return node.HeightValue
}

func (node *SimpleNode) SetChildren(newLeftChild, newRightChild avltree.Node, newHeight int) avltree.RealNode {
	node.LeftChildNode = newLeftChild
	node.RightChildNode = newRightChild
	node.HeightValue = newHeight
	return node
}

func (node *SimpleNode) Set(newLeftChild, newRightChild avltree.Node, newHeight int, newValue interface{}) avltree.RealNode {
	node.LeftChildNode = newLeftChild
	node.RightChildNode = newRightChild
	node.HeightValue = newHeight
	node.ValueData = newValue
	return node
}
