// Author: Leonardone @ NEETSDKASU
// License: MIT

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
