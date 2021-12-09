package linkedtree

import "avltree"

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
	KeyData        avltree.Key
	ValueData      interface{}
}

func unwrap(node avltree.Node) *LinkedTreeNode {
	if ltNode, ok := node.(*LinkedTreeNode); ok {
		return ltNode
	} else {
		return nil
	}
}

func (node *LinkedTreeNode) toNode() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node
	}
}

func New(allowDuplicateKeys bool) avltree.Tree {
	return &LinkedTree{
		nil, // RootNode
		allowDuplicateKeys,
	}
}

func (tree *LinkedTree) NodeCount() int {
	return tree.RootNode.NodeCount()
}

func (tree *LinkedTree) NewNode(
	leftChild,
	rightChild avltree.Node,
	height int,
	key avltree.Key,
	value interface{},
) avltree.RealNode {
	node := &LinkedTreeNode{
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

func (tree *LinkedTree) Root() avltree.Node {
	return tree.RootNode.toNode()
}

func (tree *LinkedTree) SetRoot(newRoot avltree.RealNode) avltree.RealTree {
	tree.RootNode = unwrap(newRoot)
	tree.RootNode.setParent(nil)
	return tree
}

func (tree *LinkedTree) AllowDuplicateKeys() bool {
	return tree.AllowDuplicateKeysValue
}

func (node *LinkedTreeNode) Key() avltree.Key {
	return node.KeyData
}

func (node *LinkedTreeNode) Value() interface{} {
	return node.ValueData
}

func (node *LinkedTreeNode) Height() int {
	return node.HeightValue
}

func (node *LinkedTreeNode) LeftChild() avltree.Node {
	return node.LeftChildNode.toNode()
}

func (node *LinkedTreeNode) RightChild() avltree.Node {
	return node.RightChildNode.toNode()
}

func (node *LinkedTreeNode) SetValue(newValue interface{}) avltree.Node {
	node.ValueData = newValue
	return node
}

func (node *LinkedTreeNode) Parent() avltree.Node {
	if node == nil {
		return nil
	} else {
		return node.ParentNode.toNode()
	}
}

func (node *LinkedTreeNode) setParent(newParent avltree.Node) {
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
	node.NodeCountValue = 1 +
		node.LeftChildNode.NodeCount() +
		node.RightChildNode.NodeCount()
}

func (node *LinkedTreeNode) SetChildren(
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

func (node *LinkedTreeNode) Set(
	newLeftChild,
	newRightChild avltree.Node,
	newHeight int,
	newValue interface{},
) avltree.RealNode {
	node.ValueData = newValue
	return node.SetChildren(newLeftChild, newRightChild, newHeight)
}
