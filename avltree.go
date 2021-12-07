package avltree

import "strings"

type NodeIteratorCallBack = func(node Node) (breakIteration bool)
type ValueUpdaterCallBack = func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue bool)
type UpdateIteratorCallBack = func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool)
type DeleteIteratorCallBack = func(key Key, oldValue interface{}) (deleteNode, breakIteration bool)

type NodeCounter interface{ NodeCount() int }
type NodeReleaser interface{ ReleaseNode(node RealNode) }
type ParentGetter interface {
	Node
	Parent() Node
}

type Tree interface {
	Root() Node
}

type RealTree interface {
	Tree
	NewNode(leftChild, rightChild Node, height int, key Key, value interface{}) RealNode
	SetRoot(newRoot RealNode) RealTree
	AllowDuplicateKeys() bool
}

type KeyAndValue interface {
	Key() Key
	Value() interface{}
}

type Node interface {
	KeyAndValue
	LeftChild() Node
	RightChild() Node
	SetValue(newValue interface{}) Node
}

type RealNode interface {
	Node
	Height() int
	SetChildren(newLeftChild, newRightChild Node, newHeight int) RealNode
	Set(newLeftChild, newRightChild Node, newHeight int, newValue interface{}) RealNode
}

type KeyOrdering int

const (
	LessThanOtherKey    KeyOrdering = -1
	EqualToOtherKey     KeyOrdering = 0
	GreaterThanOtherKey KeyOrdering = 1
)

type Key interface {
	CompareTo(other Key) KeyOrdering
}

type IntKey int

func (key IntKey) CompareTo(other Key) KeyOrdering {
	v1 := int(key)
	v2 := int(other.(IntKey))
	switch {
	case v1 < v2:
		return LessThanOtherKey
	case v1 > v2:
		return GreaterThanOtherKey
	default:
		return EqualToOtherKey
	}
	// return v1 - v2 は 算術オーバーフローがこわい
}

type StringKey string

func (key StringKey) CompareTo(other Key) KeyOrdering {
	s1 := string(key)
	s2 := string(other.(StringKey))
	return KeyOrdering(strings.Compare(s1, s2))
}

func Insert(tree Tree, replaceIfExists bool, key Key, value interface{}) (modified Tree, ok bool) {
	realTree := tree.(RealTree)
	helper := insertHelper{
		&realTree,
		replaceIfExists,
		&key,
		&value,
	}
	if newRoot, ok := helper.insertTo(tree.Root()); ok {
		return realTree.SetRoot(newRoot), true
	} else {
		return tree, false
	}
}

func Delete(tree Tree, key Key) (modified Tree, value interface{}, ok bool) {
	if newRoot, node, ok := removeNode(tree.Root(), key); ok {
		value = node.Value()
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
		realTree := tree.(RealTree)
		if root, ok := newRoot.(RealNode); ok {
			return realTree.SetRoot(root), value, true
		} else {
			return realTree.SetRoot(nil), value, true
		}
		// value が nil の場合もあるため、okによる削除判定が必要
	} else {
		return tree, nil, false
	}
}

func Update(tree Tree, key Key, callBack ValueUpdaterCallBack) (modified Tree, ok bool) {
	if newRoot, ok := updateValue(tree.Root(), key, callBack); ok {
		return tree.(RealTree).SetRoot(newRoot), true
	} else {
		return tree, false
	}
}

func Replace(tree Tree, key Key, newValue interface{}) (modified Tree, ok bool) {
	return Update(tree, key, func(key Key, oldValue interface{}) (interface{}, bool) {
		return newValue, false
	})
}

func Find(tree Tree, key Key) (node Node, ok bool) {
	node = tree.Root()
	for node != nil {
		cmp := key.CompareTo(node.Key())
		switch {
		case cmp.LessThan():
			node = node.LeftChild()
		case cmp.GreaterThan():
			node = node.RightChild()
		default:
			return node, true
		}
	}
	return nil, false
}

func Iterate(tree Tree, descOrder bool, callBack NodeIteratorCallBack) (ok bool) {
	if descOrder {
		return descIterateNode(tree.Root(), callBack)
	} else {
		return ascIterateNode(tree.Root(), callBack)
	}
}

func Range(tree Tree, descOrder bool, lower, upper Key) (nodes []Node) {
	RangeIterate(tree, descOrder, lower, upper, func(node Node) (breakIteration bool) {
		nodes = append(nodes, node)
		return
	})
	return
}

func RangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack NodeIteratorCallBack) (ok bool) {
	if lower == nil && upper == nil {
		return Iterate(tree, descOrder, callBack)
	}
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		return descRangeNode(tree.Root(), bounds, callBack)
	} else {
		return ascRangeNode(tree.Root(), bounds, callBack)
	}
}

func Count(tree Tree) int {
	if counter, ok := tree.(NodeCounter); ok {
		return counter.NodeCount()
	} else {
		return countNode(tree.Root())
	}
}

func CountRange(tree Tree, lower, upper Key) int {
	if lower == nil && upper == nil {
		return Count(tree)
	}
	if tree.(RealTree).AllowDuplicateKeys() {
		return countExtendedRange(tree.Root(), lower, upper)
	} else {
		return countRange(tree.Root(), lower, upper)
	}
}

func Min(tree Tree) (minimum Node, ok bool) {
	node := tree.Root()
	if node == nil {
		return nil, false
	}
	for {
		leftChild := node.LeftChild()
		if leftChild == nil {
			break
		}
		node = leftChild
	}
	return node, true
}

func Max(tree Tree) (maximum Node, ok bool) {
	node := tree.Root()
	if node == nil {
		return nil, false
	}
	for {
		rightChild := node.RightChild()
		if rightChild == nil {
			break
		}
		node = rightChild
	}
	return node, true
}

func DeleteAll(tree Tree, key Key) (modified Tree, values []interface{}, ok bool) {
	// ひとまず雑実装
	for {
		var value interface{}
		tree, value, ok = Delete(tree, key)
		if !ok {
			return tree, values, 0 < len(values)
		}
		values = append(values, value)
	}
}

func UpdateAll(tree Tree, key Key, callBack ValueUpdaterCallBack) (modified Tree, ok bool) {
	if key != nil {
		return UpdateRangeIterate(tree, false, key, key, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
			newValue, keepOldValue = callBack(key, oldValue)
			return
		})
	} else {
		return tree, false
	}
}

func ReplaceAll(tree Tree, key Key, newValue interface{}) (modified Tree, ok bool) {
	if key != nil {
		return UpdateRangeIterate(tree, false, key, key, func(key Key, oldValue interface{}) (interface{}, bool, bool) {
			return newValue, false, false
		})
	} else {
		return tree, false
	}
}

func FindAll(tree Tree, key Key) (nodes []Node, ok bool) {
	if key != nil {
		// FindAllを頻繁に呼び出すでもない限りは
		// Range呼び出しのオーバーヘッドなんて気にするほどのものではないはず
		RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
			nodes = append(nodes, node)
			return
		})
	}
	return nodes, 0 < len(nodes)
}

func MinAll(tree Tree) (minimums []Node, ok bool) {
	minimum, ok := Min(tree)
	if !ok {
		return nil, false
	}
	key := minimum.Key()
	RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
		minimums = append(minimums, node)
		return
	})
	return minimums, true
}

func MaxAll(tree Tree) (maximums []Node, ok bool) {
	maximum, ok := Max(tree)
	if !ok {
		return nil, false
	}
	key := maximum.Key()
	// 最大側は descOrder にすべきか…？
	// FindAllに揃えるなら ascOrder だが
	RangeIterate(tree, false, key, key, func(node Node) (breakIteration bool) {
		maximums = append(maximums, node)
		return
	})
	return maximums, true
}

type keyAndValue struct {
	key   Key
	value interface{}
}

func (kv *keyAndValue) Key() Key {
	return kv.key
}

func (kv *keyAndValue) Value() interface{} {
	return kv.value
}

func DeleteIterate(tree Tree, descOrder bool, callBack DeleteIteratorCallBack) (modified Tree, values []KeyAndValue) {
	var newRoot Node
	var deleted []Node
	if descOrder {
		newRoot, deleted, _ = descDeleteIterate(tree.Root(), callBack)
	} else {
		newRoot, deleted, _ = ascDeleteIterate(tree.Root(), callBack)
	}
	if len(deleted) == 0 {
		return tree, nil
	}
	for _, node := range deleted {
		values = append(values, &keyAndValue{
			node.Key(),
			node.Value(),
		})
		if releaser, ok := tree.(NodeReleaser); ok {
			releaser.ReleaseNode(node.(RealNode))
		}
	}
	if root, ok := newRoot.(RealNode); ok {
		return tree.(RealTree).SetRoot(root), values
	} else {
		return tree.(RealTree).SetRoot(nil), values
	}
}

func UpdateIterate(tree Tree, descOrder bool, callBack UpdateIteratorCallBack) (modified Tree, ok bool) {
	if descOrder {
		if newRoot, updated, _ := descUpdateIterate(tree.Root(), callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	} else {
		if newRoot, updated, _ := ascUpdateIterate(tree.Root(), callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	}
}

func UpdateRange(tree Tree, descOrder bool, lower, upper Key, callBack ValueUpdaterCallBack) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, descOrder, lower, upper, func(key Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
		newValue, keepOldValue = callBack(key, oldValue)
		return
	})
}

func UpdateRangeIterate(tree Tree, descOrder bool, lower, upper Key, callBack UpdateIteratorCallBack) (modified Tree, ok bool) {
	if lower == nil && upper == nil {
		return UpdateIterate(tree, descOrder, callBack)
	}
	bounds := newKeyBounds(lower, upper, tree.(RealTree).AllowDuplicateKeys())
	if descOrder {
		if newRoot, updated, _ := descUpdateRange(tree.Root(), bounds, callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	} else {
		if newRoot, updated, _ := ascUpdateRange(tree.Root(), bounds, callBack); updated {
			return tree.(RealTree).SetRoot(newRoot), true
		} else {
			return tree, false
		}
	}
}

func ReplaceRange(tree Tree, lower, upper Key, newValue interface{}) (modified Tree, ok bool) {
	return UpdateRangeIterate(tree, false, lower, upper, func(key Key, oldValue interface{}) (interface{}, bool, bool) {
		return newValue, false, false
	})
}

func (ordering KeyOrdering) LessThan() bool {
	return int(ordering) < 0
}

func (ordering KeyOrdering) LessThanOrEqualTo() bool {
	return int(ordering) <= 0
}

func (ordering KeyOrdering) EqualTo() bool {
	return ordering == EqualToOtherKey
}

func (ordering KeyOrdering) NotEqualTo() bool {
	return ordering != EqualToOtherKey
}

func (ordering KeyOrdering) GreaterThan() bool {
	return 0 < int(ordering)
}

func (ordering KeyOrdering) GreaterThanOrEqualTo() bool {
	return 0 <= int(ordering)
}

func (ordering KeyOrdering) Less(orEqual bool) bool {
	if orEqual {
		return ordering.LessThanOrEqualTo()
	} else {
		return ordering.LessThan()
	}
}

func (ordering KeyOrdering) Greater(orEqual bool) bool {
	if orEqual {
		return ordering.GreaterThanOrEqualTo()
	} else {
		return ordering.GreaterThan()
	}
}

func countNode(node Node) int {
	if node == nil {
		return 0
	}
	// countNodeの再帰での呼び出し頻度を考えると
	// この判定は無駄が多く、コストが高くつく、かも
	// Treeを構成する一部のNodeにだけNodeCounterが実装されている可能性は低いと思う
	if counter, ok := node.(NodeCounter); ok {
		return counter.NodeCount()
	}
	return 1 + countNode(node.LeftChild()) + countNode(node.RightChild())
}

func countRange(node Node, lower, upper Key) int {
	if node == nil {
		return 0
	}
	// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
	// lower == nil, upper < key    ... leftChild
	// lower == nil, key == upper   ... all(leftChild) key
	// lower == nil, key < upper    ... all(leftChild) key rightChild
	// upper == nil, lower < key    ... leftChild key all(rightChild)
	// upper == nil, lower == key   ... key all(rightChild)
	// upper == nil, key < lower    ... rightChild
	// lower < upper < key      ... leftChild
	// lower < upper == key     ... leftChild{upper->nil} key
	// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
	// key == lower < upper     ... key rightChild{lower->nil}
	// key < lower < upper      ... rightChild
	// lower == key == upper    ... key
	if lower == nil {
		if upper == nil {
			// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
			return countNode(node)
		}
		cmp := node.Key().CompareTo(upper)
		switch {
		case cmp.GreaterThan():
			// lower == nil, upper < key    ... leftChild
			return countRange(node.LeftChild(), lower, upper)
		case cmp.EqualTo():
			// lower == nil, key == upper   ... all(leftChild) key
			return countNode(node.LeftChild()) + 1
		case cmp.LessThan():
			// lower == nil, key < upper    ... all(leftChild) key rightChild
			return countNode(node.LeftChild()) + 1 + countRange(node.RightChild(), lower, upper)
		default:
			// ここには到達しないはず
			panic("unreachable?")
		}
	}
	if upper == nil {
		cmp := node.Key().CompareTo(lower)
		switch {
		case cmp.GreaterThan():
			// upper == nil, lower < key    ... leftChild key all(rightChild)
			return countRange(node.LeftChild(), lower, upper) + 1 + countNode(node.RightChild())
		case cmp.EqualTo():
			// upper == nil, lower == key   ... key all(rightChild)
			return 1 + countNode(node.RightChild())
		case cmp.LessThan():
			// upper == nil, key < lower    ... rightChild
			return countRange(node.RightChild(), lower, upper)
		default:
			// ここには到達しないはず
			panic("unreachable?")
		}
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	cmpUpper := key.CompareTo(upper)
	switch {
	case cmpUpper.GreaterThan():
		// lower < upper < key      ... leftChild
		return countRange(node.LeftChild(), lower, upper)
	case cmpLower.GreaterThan() && cmpUpper.EqualTo():
		// lower < upper == key     ... leftChild{upper->nil} key
		return countRange(node.LeftChild(), lower, nil) + 1
	case cmpLower.GreaterThan() && cmpUpper.LessThan():
		// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
		return countRange(node.LeftChild(), lower, nil) + 1 + countRange(node.RightChild(), nil, upper)
	case cmpLower.EqualTo() && cmpUpper.LessThan():
		// key == lower < upper     ... key rightChild{lower->nil}
		return 1 + countRange(node.RightChild(), nil, upper)
	case cmpLower.LessThan():
		// key < lower < upper      ... rightChild
		return countRange(node.RightChild(), lower, upper)
	case cmpLower.EqualTo() && cmpUpper.EqualTo():
		// lower == key == upper    ... key
		return 1
	}
	// 条件漏れが無ければ、ここには到達しないと思う…
	panic("unreachable?")
}

func countExtendedRange(node Node, lower, upper Key) int {
	if node == nil {
		return 0
	}
	// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
	// lower == nil, upper < key    ... leftChild
	// lower == nil, key == upper   ... all(leftChild) key rightChild
	// lower == nil, key < upper    ... all(leftChild) key rightChild
	// upper == nil, lower < key    ... leftChild key all(rightChild)
	// upper == nil, lower == key   ... leftChild key all(rightChild)
	// upper == nil, key < lower    ... rightChild
	// lower < upper < key      ... leftChild
	// lower < upper == key     ... leftChild{upper->nil} key rightChild{lower->nil}
	// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
	// key == lower < upper     ... leftChild{upper->nil} key rightChild{lower->nil}
	// key < lower < upper      ... rightChild
	// lower == key == upper    ... leftChild{upper->nil} key rightChild{lower->nil}
	if lower == nil {
		if upper == nil {
			// lower == nil, upper == nil   ... all(leftChild) key all(rightChild)
			return countNode(node)
		}
		if node.Key().CompareTo(upper).GreaterThan() {
			// lower == nil, upper < key    ... leftChild
			return countExtendedRange(node.LeftChild(), lower, upper)
		} else {
			// lower == nil, key == upper   ... all(leftChild) key rightChild
			// lower == nil, key < upper    ... all(leftChild) key rightChild
			return countNode(node.LeftChild()) + 1 + countExtendedRange(node.RightChild(), lower, upper)
		}
	}
	if upper == nil {
		if node.Key().CompareTo(lower).GreaterThanOrEqualTo() {
			// upper == nil, lower < key    ... leftChild key all(rightChild)
			// upper == nil, lower == key   ... leftChild key all(rightChild)
			return countExtendedRange(node.LeftChild(), lower, upper) + 1 + countNode(node.RightChild())
		} else {
			// upper == nil, key < lower    ... rightChild
			return countExtendedRange(node.RightChild(), lower, upper)
		}
	}
	key := node.Key()
	cmpLower := key.CompareTo(lower)
	cmpUpper := key.CompareTo(upper)
	switch {
	case cmpUpper.GreaterThan():
		// lower < upper < key      ... leftChild
		return countExtendedRange(node.LeftChild(), lower, upper)
	case cmpLower.LessThan():
		// key < lower < upper      ... rightChild
		return countExtendedRange(node.RightChild(), lower, upper)
	default:
		// lower < upper == key     ... leftChild{upper->nil} key rightChild{lower->nil}
		// lower < key < upper      ... leftChild{upper->nil} key rightChild{lower->nil}
		// key == lower < upper     ... leftChild{upper->nil} key rightChild{lower->nil}
		// lower == key == upper    ... leftChild{upper->nil} key rightChild{lower->nil}
		return countExtendedRange(node.LeftChild(), lower, nil) + 1 + countExtendedRange(node.RightChild(), nil, upper)
	}
}

func getHeight(node Node) int {
	if node == nil {
		return 0
	} else {
		return node.(RealNode).Height()
	}
}

type balance int

const (
	balanced      balance = 0
	leftIsHigher  balance = -1
	rightIsHigher balance = 1
)

func checkBalance(node RealNode) balance {
	if node == nil {
		return balanced
	}
	heightL := getHeight(node.LeftChild())
	heightR := getHeight(node.RightChild())
	switch {
	case heightL < heightR && heightL+1 < heightR:
		return rightIsHigher
	case heightL > heightR && heightL > heightR+1:
		return leftIsHigher
	default:
		return balanced
	}
	// 算術オーバーフローが怖いのか？
	// heightL - heightR
	// heightL + 1
	// heightR + 1
	// だが、現実的なAVL木の高さを考えるとせいぜい高くても33～34くらいでは？
	// 高さ10で完全二分木のノード総数は(2の10乗-1)個=1023個
	// 高さ32で4億個？最小構成のノードでもノード1つ14bytesくらいかもだろうし
	// (最小構成、leftChild,rightChild　32bitsサイズのアドレス値, height,data int8, key int32)
	// 4億個もあったら最小構成でも56GB以上のメモリ・ストレージを必要とするわけで…
	// それ以上だとメモリアドレスが64bitsサイズ、必要な容量がグっと増えるし・・・
	// そもオーバーフローするとしたら内部データが意図せず破壊された場合のみで
	// その場合は正常動作を保証する必要もないわけで…
}

func compareNodeHeight(leftNode, rightNode Node) balance {
	heightL := getHeight(leftNode)
	heightR := getHeight(rightNode)
	// 算術オーバーフローが怖いのかい？
	// heightL - heightR
	switch {
	case heightL < heightR:
		return rightIsHigher
	case heightL > heightR:
		return leftIsHigher
	default:
		return balanced
	}
}

func compareChildHeight(node Node) balance {
	if node == nil {
		return balanced
	} else {
		return compareNodeHeight(node.LeftChild(), node.RightChild())
	}
}

func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func calcNewHeight(leftChild, rightChild Node) int {
	leftHeight := getHeight(leftChild)
	rightHeight := getHeight(rightChild)
	return 1 + intMax(leftHeight, rightHeight)
}

func setChildren(root RealNode, leftChild, rightChild Node) RealNode {
	newHeight := calcNewHeight(leftChild, rightChild)
	return root.SetChildren(leftChild, rightChild, newHeight)
}

func setLeftChild(root RealNode, newLeftChild Node) RealNode {
	return setChildren(root, newLeftChild, root.RightChild())
}

func setRightChild(root RealNode, newRightChild Node) RealNode {
	return setChildren(root, root.LeftChild(), newRightChild)
}

func resetNode(root RealNode, newLeftChild, newRightChild Node, newValue interface{}) RealNode {
	newHeight := calcNewHeight(newLeftChild, newRightChild)
	return root.Set(newLeftChild, newRightChild, newHeight, newValue)
}

type insertHelper struct {
	tree            *RealTree
	replaceIfExists bool
	key             *Key
	value           *interface{}
}

func (helper *insertHelper) newNode() RealNode {
	return (*helper.tree).NewNode(nil, nil, 1, *helper.key, *helper.value)
}

func (helper *insertHelper) compareKey(node Node) KeyOrdering {
	return (*helper.key).CompareTo(node.Key())
}

func (helper *insertHelper) allowDuplicateKeys() bool {
	return (*helper.tree).AllowDuplicateKeys()
}

func (helper *insertHelper) insertTo(root Node) (newRoot RealNode, ok bool) {
	if root == nil {
		return helper.newNode(), true
	}
	cmp := helper.compareKey(root)
	switch {
	case cmp.LessThan(): // newKey < root.key
		if newLeftChild, ok := helper.insertTo(root.LeftChild()); ok {
			newRoot = setLeftChild(root.(RealNode), newLeftChild)
		} else {
			return root.(RealNode), false
		}
	case cmp.GreaterThan(): // root.key < newKey
		if newRightChild, ok := helper.insertTo(root.RightChild()); ok {
			newRoot = setRightChild(root.(RealNode), newRightChild)
		} else {
			return root.(RealNode), false
		}
	default:
		if helper.replaceIfExists {
			newRoot = root.SetValue(*helper.value).(RealNode)
			return newRoot, true
		}
		if !helper.allowDuplicateKeys() {
			return root.(RealNode), false
		}
		if newRightChild, ok := helper.insertTo(root.RightChild()); ok {
			newRoot = setRightChild(root.(RealNode), newRightChild)
		} else {
			return root.(RealNode), false
		}
	}

	newRoot = rotate(newRoot)

	return newRoot, true
}

func rotate(root RealNode) RealNode {
	// 無限ループは不要な気がする
	// 複数のノードをまとめて削除するとバランス崩れが発生しそうだからその時必要？
	for {
		switch checkBalance(root) {
		case leftIsHigher:
			root = rotateRight(root)
		case rightIsHigher:
			root = rotateLeft(root)
		case balanced:
			return root
		default:
			panic("unreachable")
		}
	}
}

func rotateRight(root RealNode) RealNode {
	oldRootLeftChild := root.LeftChild().(RealNode)
	if compareChildHeight(oldRootLeftChild) == rightIsHigher {
		newRoot := oldRootLeftChild.RightChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootRightChild := rotate(setLeftChild(root, tempRightChild))
		newRootLeftChild := rotate(setRightChild(oldRootLeftChild, tempLeftChild))
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempRightChild := oldRootLeftChild.RightChild()
		newRootRightChild := rotate(setLeftChild(root, tempRightChild))
		return setRightChild(oldRootLeftChild, newRootRightChild)
	}
}

func rotateLeft(root RealNode) RealNode {
	oldRootRightChild := root.RightChild().(RealNode)
	if compareChildHeight(oldRootRightChild) == leftIsHigher {
		newRoot := oldRootRightChild.LeftChild().(RealNode)
		tempLeftChild := newRoot.LeftChild()
		tempRightChild := newRoot.RightChild()
		newRootLeftChild := rotate(setRightChild(root, tempLeftChild))
		newRootRightChild := rotate(setLeftChild(oldRootRightChild, tempRightChild))
		return setChildren(newRoot, newRootLeftChild, newRootRightChild)
	} else {
		tempLeftChild := oldRootRightChild.LeftChild()
		newLeftChild := rotate(setRightChild(root, tempLeftChild))
		return setLeftChild(oldRootRightChild, newLeftChild)
	}
}

func removeNode(root Node, key Key) (newRoot, removed Node, ok bool) {
	if root == nil {
		return nil, nil, false
	}
	cmp := key.CompareTo(root.Key())
	switch {
	case cmp.LessThan(): // key < root.Key()
		if tempLeftChild, node, ok := removeNode(root.LeftChild(), key); ok {
			removed = node
			newRoot = setLeftChild(root.(RealNode), tempLeftChild)
		} else {
			return nil, nil, false
		}
	case cmp.GreaterThan(): // root.Key() < key
		if tempRightChild, node, ok := removeNode(root.RightChild(), key); ok {
			removed = node
			newRoot = setRightChild(root.(RealNode), tempRightChild)
		} else {
			return nil, nil, false
		}
	default: // just target node
		removed = root
		leftChild := root.LeftChild()
		rightChild := root.RightChild()
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot == nil {
			// no children
			// height(root) == 1
			// root.LeftChild() == nil
			// root.RightChild() == nil
			return nil, removed, true
		}
		newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
	}

	newRoot = rotate(newRoot.(RealNode))

	return newRoot, removed, true
}

func removeMin(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	leftChild := root.LeftChild()
	if leftChild == nil {
		newRoot = root.RightChild()
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, root
	}
	leftChild, removed = removeMin(leftChild)
	newRoot = setLeftChild(root.(RealNode), leftChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

func removeMax(root Node) (newRoot, removed Node) {
	if root == nil {
		return nil, nil
	}
	rightChild := root.RightChild()
	if rightChild == nil {
		newRoot = root.LeftChild()
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, root
	}
	rightChild, removed = removeMax(rightChild)
	newRoot = setRightChild(root.(RealNode), rightChild)
	newRoot = rotate(newRoot.(RealNode))
	return newRoot, removed
}

func ascIterateNode(node Node, callBack NodeIteratorCallBack) (ok bool) {
	if node == nil {
		return true
	}
	if !ascIterateNode(node.LeftChild(), callBack) {
		return false
	}
	if breakIteration := callBack(node); breakIteration {
		return false
	}
	return ascIterateNode(node.RightChild(), callBack)
}

func descIterateNode(node Node, callBack NodeIteratorCallBack) (ok bool) {
	if node == nil {
		return true
	}
	if !descIterateNode(node.RightChild(), callBack) {
		return false
	}
	if breakIteration := callBack(node); breakIteration {
		return false
	}
	return descIterateNode(node.LeftChild(), callBack)
}

type boundsChecker interface {
	includeLower() bool
	includeKey() bool
	includeUpper() bool
}

type keyBounds interface {
	checkLower(key Key) boundsChecker
	checkUpper(key Key) boundsChecker
}

func newKeyBounds(lower, upper Key, extended bool) keyBounds {
	if lower == nil {
		return &upperBound{upper, extended}
	} else if upper == nil {
		return &lowerBound{lower, extended}
	} else {
		return &bothBounds{lower, upper, extended}
	}
}

type bothBounds struct {
	lower, upper Key
	ext          bool
}

func (bounds *bothBounds) checkLower(key Key) boundsChecker {
	return &lowerBoundsChecker{key.CompareTo(bounds.lower), bounds.ext}
}

func (bounds *bothBounds) checkUpper(key Key) boundsChecker {
	return &upperBoundsChecker{key.CompareTo(bounds.upper), bounds.ext}
}

type upperBound struct {
	upper Key
	ext   bool
}

func (bounds *upperBound) checkLower(key Key) boundsChecker {
	return noBoundsChecker{}
}

func (bounds *upperBound) checkUpper(key Key) boundsChecker {
	return &upperBoundsChecker{key.CompareTo(bounds.upper), bounds.ext}
}

type lowerBound struct {
	lower Key
	ext   bool
}

func (bounds *lowerBound) checkLower(key Key) boundsChecker {
	return &lowerBoundsChecker{key.CompareTo(bounds.lower), bounds.ext}
}

func (bounds *lowerBound) checkUpper(key Key) boundsChecker {
	return noBoundsChecker{}
}

type noBoundsChecker struct{}

func (noBoundsChecker) includeLower() bool { return true }
func (noBoundsChecker) includeKey() bool   { return true }
func (noBoundsChecker) includeUpper() bool { return true }

type upperBoundsChecker struct {
	cmpUpper KeyOrdering
	ext      bool
}

func (checker *upperBoundsChecker) includeLower() bool {
	return true
}

func (checker *upperBoundsChecker) includeKey() bool {
	return checker.cmpUpper.LessThanOrEqualTo()
}

func (checker *upperBoundsChecker) includeUpper() bool {
	return checker.cmpUpper.Less(checker.ext)
}

type lowerBoundsChecker struct {
	cmpLower KeyOrdering
	ext      bool
}

func (checker *lowerBoundsChecker) includeLower() bool {
	return checker.cmpLower.Greater(checker.ext)
}

func (checker *lowerBoundsChecker) includeKey() bool {
	return checker.cmpLower.GreaterThanOrEqualTo()
}

func (checker *lowerBoundsChecker) includeUpper() bool {
	return true
}

func ascRangeNode(node Node, bounds keyBounds, callBack NodeIteratorCallBack) (ok bool) {
	if node == nil {
		return true
	}
	key := node.Key()
	lower := bounds.checkLower(key)
	if lower.includeLower() {
		if !ascRangeNode(node.LeftChild(), bounds, callBack) {
			return false
		}
	}
	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		if breakIteration := callBack(node); breakIteration {
			return false
		}
	}
	if upper.includeUpper() {
		return ascRangeNode(node.RightChild(), bounds, callBack)
	} else {
		return true
	}
}

func descRangeNode(node Node, bounds keyBounds, callBack NodeIteratorCallBack) (ok bool) {
	if node == nil {
		return true
	}
	key := node.Key()
	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		if !descRangeNode(node.RightChild(), bounds, callBack) {
			return false
		}
	}
	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		if breakIteration := callBack(node); breakIteration {
			return false
		}
	}
	if lower.includeLower() {
		return descRangeNode(node.LeftChild(), bounds, callBack)
	} else {
		return true
	}
}

func updateValue(node Node, key Key, callBack ValueUpdaterCallBack) (newNode RealNode, ok bool) {
	if node == nil {
		return nil, false
	}
	nodeKey := node.Key()
	cmp := key.CompareTo(nodeKey)
	switch {
	case cmp.LessThan():
		if leftChild, ok := updateValue(node.LeftChild(), key, callBack); ok {
			return setLeftChild(node.(RealNode), leftChild), true
		} else {
			return nil, false
		}
	case cmp.EqualTo():
		if newValue, keepOldValue := callBack(nodeKey, node.Value()); !keepOldValue {
			return node.SetValue(newValue).(RealNode), true
		} else {
			return nil, false
		}
	case cmp.GreaterThan():
		if rightChild, ok := updateValue(node.RightChild(), key, callBack); ok {
			return setRightChild(node.(RealNode), rightChild), true
		} else {
			return nil, false
		}
	default:
		panic("unreachable")
	}
}

func ascUpdateIterate(root Node, callBack UpdateIteratorCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}

	leftChild, leftUpdated, breakIteration := ascUpdateIterate(root.LeftChild(), callBack)
	if breakIteration {
		if leftUpdated {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root.(RealNode)
		}
		return newRoot, leftUpdated, breakIteration
	}

	newValue, keepOldValue, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		switch {
		case !leftUpdated && keepOldValue:
			newRoot = root.(RealNode)
		case !leftUpdated && !keepOldValue:
			newRoot = root.SetValue(newValue).(RealNode)
		case leftUpdated && keepOldValue:
			newRoot = setLeftChild(root.(RealNode), leftChild)
		case leftUpdated && !keepOldValue:
			newRoot = resetNode(root.(RealNode), leftChild, root.RightChild(), newValue)
		}
		updated = leftUpdated || !keepOldValue
		return newRoot, updated, breakIteration
	}

	rightChild, rightUpdated, breakIteration := ascUpdateIterate(root.RightChild(), callBack)
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

func descUpdateIterate(root Node, callBack UpdateIteratorCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}

	rightChild, rightUpdated, breakIteration := descUpdateIterate(root.RightChild(), callBack)
	if breakIteration {
		if rightUpdated {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root.(RealNode)
		}
		return newRoot, rightUpdated, breakIteration
	}

	newValue, keepOldValue, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		switch {
		case !rightUpdated && keepOldValue:
			newRoot = root.(RealNode)
		case !rightUpdated && !keepOldValue:
			newRoot = root.SetValue(newValue).(RealNode)
		case rightUpdated && keepOldValue:
			newRoot = setRightChild(root.(RealNode), rightChild)
		case rightUpdated && !keepOldValue:
			newRoot = resetNode(root.(RealNode), root.LeftChild(), rightChild, newValue)
		default:
			panic("unreachable")
		}
		updated = rightUpdated || !keepOldValue
		return newRoot, updated, breakIteration
	}

	leftChild, leftUpdated, breakIteration := descUpdateIterate(root.LeftChild(), callBack)
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

func ascUpdateRange(root Node, bounds keyBounds, callBack UpdateIteratorCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}
	var leftUpdated, keepOldValue, rightUpdated bool
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	lower := bounds.checkLower(key)
	if lower.includeLower() {
		leftChild, leftUpdated, breakIteration = ascUpdateRange(leftChild, bounds, callBack)
		if breakIteration {
			if leftUpdated {
				newRoot = setLeftChild(root.(RealNode), leftChild)
			} else {
				newRoot = root.(RealNode)
			}
			return newRoot, leftUpdated, breakIteration
		}
	}

	var newValue interface{}
	upper := bounds.checkUpper(key)
	if lower.includeKey() && upper.includeKey() {
		newValue, keepOldValue, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			switch {
			case !leftUpdated && keepOldValue:
				newRoot = root.(RealNode)
			case !leftUpdated && !keepOldValue:
				newRoot = root.SetValue(newValue).(RealNode)
			case leftUpdated && keepOldValue:
				newRoot = setLeftChild(root.(RealNode), leftChild)
			case leftUpdated && !keepOldValue:
				newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
			default:
				panic("unreachable")
			}
			updated = leftUpdated || !keepOldValue
			return newRoot, updated, breakIteration
		}
	} else {
		keepOldValue = true
	}

	if upper.includeUpper() {
		rightChild, rightUpdated, breakIteration = ascUpdateRange(rightChild, bounds, callBack)
	}
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

func descUpdateRange(root Node, bounds keyBounds, callBack UpdateIteratorCallBack) (newRoot RealNode, updated, breakIteration bool) {
	if root == nil {
		return nil, false, false
	}
	var leftUpdated, keepOldValue, rightUpdated bool
	leftChild := root.LeftChild()
	rightChild := root.RightChild()
	key := root.Key()

	upper := bounds.checkUpper(key)
	if upper.includeUpper() {
		rightChild, rightUpdated, breakIteration = descUpdateRange(rightChild, bounds, callBack)
		if breakIteration {
			if rightUpdated {
				newRoot = setRightChild(root.(RealNode), rightChild)
			} else {
				newRoot = root.(RealNode)
			}
			return newRoot, rightUpdated, breakIteration
		}
	}

	var newValue interface{}
	lower := bounds.checkLower(key)
	if lower.includeKey() && upper.includeKey() {
		newValue, keepOldValue, breakIteration = callBack(root.Key(), root.Value())
		if breakIteration {
			switch {
			case !rightUpdated && keepOldValue:
				newRoot = root.(RealNode)
			case !rightUpdated && !keepOldValue:
				newRoot = root.SetValue(newValue).(RealNode)
			case rightUpdated && keepOldValue:
				newRoot = setRightChild(root.(RealNode), rightChild)
			case rightUpdated && !keepOldValue:
				newRoot = resetNode(root.(RealNode), root.LeftChild(), rightChild, newValue)
			default:
				panic("unreachable")
			}
			updated = rightUpdated || !keepOldValue
			return newRoot, updated, breakIteration
		}
	} else {
		keepOldValue = true
	}

	if lower.includeLower() {
		leftChild, leftUpdated, breakIteration = descUpdateRange(leftChild, bounds, callBack)
	}
	switch {
	case !leftUpdated && keepOldValue && !rightUpdated:
		newRoot = root.(RealNode)
	case !leftUpdated && !keepOldValue && !rightUpdated:
		newRoot = root.SetValue(newValue).(RealNode)
	case keepOldValue:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case !keepOldValue:
		newRoot = resetNode(root.(RealNode), leftChild, rightChild, newValue)
	default:
		panic("unreachable")
	}
	updated = leftUpdated || !keepOldValue || rightUpdated
	return newRoot, updated, breakIteration
}

func ascDeleteIterate(root Node, callBack DeleteIteratorCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}

	leftChild, leftDeleted, breakIteration := ascDeleteIterate(root.LeftChild(), callBack)
	if breakIteration {
		if len(leftDeleted) > 0 {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root
		}
		newRoot = rotate(newRoot.(RealNode))
		return newRoot, leftDeleted, breakIteration
	}

	deleteRoot, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		deleted = leftDeleted
		if deleteRoot {
			deleted = append(deleted, root)
			rightChild := root.RightChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
			}
		} else if len(leftDeleted) > 0 {
			newRoot = setLeftChild(root.(RealNode), leftChild)
		} else {
			newRoot = root
		}
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, deleted, breakIteration
	}

	rightChild, rightDeleted, breakIteration := ascDeleteIterate(root.RightChild(), callBack)
	deleted = leftDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, rightDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}

func descDeleteIterate(root Node, callBack DeleteIteratorCallBack) (newRoot Node, deleted []Node, breakIteration bool) {
	if root == nil {
		return nil, nil, false
	}

	rightChild, rightDeleted, breakIteration := descDeleteIterate(root.RightChild(), callBack)
	if breakIteration {
		if len(rightDeleted) > 0 {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root
		}
		return rotate(newRoot.(RealNode)), rightDeleted, breakIteration
	}

	deleteRoot, breakIteration := callBack(root.Key(), root.Value())
	if breakIteration {
		deleted = rightDeleted
		if deleteRoot {
			deleted = append(deleted, root)
			leftChild := root.LeftChild()
			if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
				leftChild, newRoot = removeMax(leftChild)
			} else {
				rightChild, newRoot = removeMin(rightChild)
			}
			if newRoot != nil {
				newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
			}
		} else if len(rightDeleted) > 0 {
			newRoot = setRightChild(root.(RealNode), rightChild)
		} else {
			newRoot = root
		}
		if newRoot != nil {
			newRoot = rotate(newRoot.(RealNode))
		}
		return newRoot, deleted, breakIteration
	}

	leftChild, leftDeleted, breakIteration := descDeleteIterate(root.LeftChild(), callBack)

	deleted = rightDeleted
	switch {
	case len(leftDeleted) == 0 && !deleteRoot && len(rightDeleted) == 0:
		newRoot = root
	case !deleteRoot:
		newRoot = setChildren(root.(RealNode), leftChild, rightChild)
	case deleteRoot:
		deleted = append(deleted, root)
		if compareNodeHeight(leftChild, rightChild) == leftIsHigher {
			leftChild, newRoot = removeMax(leftChild)
		} else {
			rightChild, newRoot = removeMin(rightChild)
		}
		if newRoot != nil {
			newRoot = setChildren(newRoot.(RealNode), leftChild, rightChild)
		}
	default:
		panic("unreachable")
	}
	deleted = append(deleted, leftDeleted...)
	if newRoot != nil {
		newRoot = rotate(newRoot.(RealNode))
	}
	return newRoot, deleted, breakIteration
}
