package linkedtree

import "avltree"

type (
	Key      = avltree.Key
	Node     = avltree.Node
	RealNode = avltree.RealNode
	Tree     = avltree.Tree
	RealTree = avltree.RealTree
)

type LinkedTree struct {
	RootNode                *LinkedTreeNode
	AllowDuplicateKeysValue bool
}

type LinkedTreeNode struct {
	LeftChildNode  *LinkedTreeNode
	RightChildNode *LinkedTreeNode
	HeightValue    int
	ParentNode     *LinkedTreeNode
	NodeCountValue int
	KeyData        Key
	ValueData      interface{}
}

func unwrap(node Node) *LinkedTreeNode {
	if ltNode, ok := node.(*LinkedTreeNode); ok {
		return ltNode
	} else {
		return nil
	}
}

func (node *LinkedTreeNode) toNode() Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func New(allowDuplicateKeys bool) Tree {
	return &LinkedTree{
		nil, // root
		allowDuplicateKeys,
	}
}

func (tree *LinkedTree) NodeCount() int {
	return tree.RootNode.NodeCount()
}

func (tree *LinkedTree) NewNode(leftChild, rightChild Node, height int, key Key, value interface{}) RealNode {
	node := &LinkedTreeNode{
		unwrap(leftChild),
		unwrap(rightChild),
		height,
		nil, // parent
		1,   // nodeCount
		key,
		value,
	}
	node.resetNodeCount()
	return node
}

func (tree *LinkedTree) Root() Node {
	return tree.RootNode.toNode()
}

func (tree *LinkedTree) SetRoot(newRoot RealNode) RealTree {
	tree.RootNode = unwrap(newRoot)
	tree.RootNode.setParent(nil)
	return tree
}

func (tree *LinkedTree) AllowDuplicateKeys() bool {
	return tree.AllowDuplicateKeysValue
}

func (node *LinkedTreeNode) Key() Key {
	return node.KeyData
}

func (node *LinkedTreeNode) Value() interface{} {
	return node.ValueData
}

func (node *LinkedTreeNode) Height() int {
	return node.HeightValue
}

func (node *LinkedTreeNode) LeftChild() Node {
	return node.LeftChildNode.toNode()
}

func (node *LinkedTreeNode) RightChild() Node {
	return node.RightChildNode.toNode()
}

func (node *LinkedTreeNode) SetValue(newValue interface{}) Node {
	node.ValueData = newValue
	return node
}

func (node *LinkedTreeNode) Parent() Node {
	if node == nil {
		return nil
	} else {
		return node.ParentNode.toNode()
	}
}

func (node *LinkedTreeNode) setParent(newParent Node) {
	if node != nil {
		node.ParentNode = unwrap(newParent)
	}
}

func (node *LinkedTreeNode) NodeCount() int {
	if node == nil {
		return 0
	} else {
		return node.NodeCountValue
	}
}

func (node *LinkedTreeNode) resetNodeCount() {
	node.NodeCountValue = 1 + node.LeftChildNode.NodeCount() + node.RightChildNode.NodeCount()
}

func (node *LinkedTreeNode) SetChildren(newLeftChild, newRightChild Node, newHeight int) RealNode {
	node.LeftChildNode = unwrap(newLeftChild)
	node.RightChildNode = unwrap(newRightChild)
	node.HeightValue = newHeight
	node.LeftChildNode.setParent(node)
	node.RightChildNode.setParent(node)
	node.resetNodeCount()
	return node
}

func (node *LinkedTreeNode) Set(newLeftChild, newRightChild Node, newHeight int, newValue interface{}) RealNode {
	node.ValueData = newValue
	return node.SetChildren(newLeftChild, newRightChild, newHeight)
}
