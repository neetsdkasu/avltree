// Author: Leonardone @ NEETSDKASU
// License: MIT

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
