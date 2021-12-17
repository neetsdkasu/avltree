// Author: Leonardone @ NEETSDKASU
// License: MIT

// github.com/neetsdkasu/avltreeのためのラッパーの実装例
// キーと値のやりとりをint型に制限する
// 内部での実際のキーはintkeyのIntKeyを使用している
//
// int型に制限するのはラッパーのメソッドの引数や戻り値だけではなく
// KeyAndValue,Node,AlterNode,AlterRequestもintwrapper独自でint型用に定義しなおされている
//
// コード例
//
//		import (
//			"fmt"
//			"github.com/neetsdkasu/avltree/intwrapper"
//			"github.com/neetsdkasu/avltree/simpletree"
//		)
//		func Example_intwrapper() {
//			tree := simpletree.New(false)
//			w := intwrapper.New(tree)
//			w.Insert(12, 345)
//			w.Insert(67, 890)
//			w.Insert(333, 666)
//			w.Insert(-5, 12345)
//			w.Delete(67)
//			w.Update(333, func(key, oldValue int) (newValue int, keepOldValue bool) {
//				newValue = oldValue * 3
//				return
//			})
//			if node := w.Find(12); node != nil {
//				fmt.Println("Find!", node.Key(), node.Value())
//			}
//			w.Iterate(func(node intwrapper.Node) (breakIteration bool) {
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
package intwrapper

import (
	"github.com/neetsdkasu/avltree"
	"github.com/neetsdkasu/avltree/intkey"
)

type IterateCallBack = func(node Node) (breakIteration bool)
type UpdateValueCallBack = func(key, oldValue int) (newValue int, keepOldValue bool)
type UpdateIterateCallBack = func(key, oldValue int) (newValue int, keepOldValue, breakIteration bool)
type DeleteIterateCallBack = func(key, value int) (deleteNode, breakIteration bool)
type AlterNodeCallBack = func(node AlterNode) (request AlterRequest)
type AlterIterateCallBack = func(node AlterNode) (request AlterRequest, breakIteration bool)

type IntAVLTree struct {
	Tree avltree.Tree
}

type KeyAndValue interface {
	Key() int
	Value() int
}

type keyAndValueWrapper struct {
	inner avltree.KeyAndValue
}

type Node interface {
	KeyAndValue
	LeftChild() Node
	RightChild() Node
	SetValue(newValue int)
}

type nodeWrapper struct {
	inner avltree.Node
}

type AlterNode interface {
	KeyAndValue
	Keep() AlterRequest
	Replace(newValue int) AlterRequest
	Delete() AlterRequest
}

type alterNodeWrapper struct {
	inner avltree.AlterNode
}

type AlterRequest struct {
	inner avltree.AlterRequest
}

func New(tree avltree.Tree) *IntAVLTree {
	return &IntAVLTree{tree}
}

func (tree *IntAVLTree) Insert(key, value int) (ok bool) {
	tree.Tree, ok = avltree.Insert(tree.Tree, false, intkey.IntKey(key), value)
	return
}

func (tree *IntAVLTree) InsertOrReplace(key, value int) (ok bool) {
	tree.Tree, ok = avltree.Insert(tree.Tree, true, intkey.IntKey(key), value)
	return
}

func (tree *IntAVLTree) Delete(key int) (deletedValue KeyAndValue) {
	var tempDeletedValue avltree.KeyAndValue
	tree.Tree, tempDeletedValue = avltree.Delete(tree.Tree, intkey.IntKey(key))
	deletedValue = wrapKeyAndValue(tempDeletedValue)
	return
}

func (tree *IntAVLTree) Update(key int, callBack UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.Update(tree.Tree, intkey.IntKey(key), wrapUpdateValueCallBack(callBack))
	return
}

func (tree *IntAVLTree) Replace(key int, value int) (ok bool) {
	tree.Tree, ok = avltree.Replace(tree.Tree, intkey.IntKey(key), value)
	return
}

func (tree *IntAVLTree) Alter(key int, callBack AlterNodeCallBack) (deletedValue KeyAndValue, ok bool) {
	var tempDeletedValue avltree.KeyAndValue
	tree.Tree, tempDeletedValue, ok = avltree.Alter(tree.Tree, intkey.IntKey(key), wrapAlterNodeCallBack(callBack))
	deletedValue = wrapKeyAndValue(tempDeletedValue)
	return
}

func (tree *IntAVLTree) Clear() {
	tree.Tree = avltree.Clear(tree.Tree)
}

func (tree *IntAVLTree) Release() {
	avltree.Release(&tree.Tree)
}

func (tree *IntAVLTree) Find(key int) (node Node) {
	return wrapNode(avltree.Find(tree.Tree, intkey.IntKey(key)))
}

func (tree *IntAVLTree) Iterate(callBack IterateCallBack) {
	avltree.Iterate(tree.Tree, false, wrapIterateCallBack(callBack))
}

func (tree *IntAVLTree) IterateRev(callBack IterateCallBack) {
	avltree.Iterate(tree.Tree, true, wrapIterateCallBack(callBack))
}

func (tree *IntAVLTree) Range(lower, upper int) (nodes []Node) {
	return wrapNodes(avltree.Range(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper)))
}

func (tree *IntAVLTree) RangeRev(lower, upper int) (nodes []Node) {
	return wrapNodes(avltree.Range(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper)))
}

func (tree *IntAVLTree) RangeIterate(lower, upper int, callBack IterateCallBack) {
	avltree.RangeIterate(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapIterateCallBack(callBack))
}

func (tree *IntAVLTree) RangeIterateRev(lower, upper int, callBack IterateCallBack) {
	avltree.RangeIterate(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapIterateCallBack(callBack))
}

func (tree *IntAVLTree) Count() int {
	return avltree.Count(tree.Tree)
}

func (tree *IntAVLTree) CountRange(lower, upper int) int {
	return avltree.CountRange(tree.Tree, intkey.IntKey(lower), intkey.IntKey(upper))
}

func (tree *IntAVLTree) Min() (node Node) {
	return wrapNode(avltree.Min(tree.Tree))
}

func (tree *IntAVLTree) Max() (node Node) {
	return wrapNode(avltree.Max(tree.Tree))
}

func (tree *IntAVLTree) DeleteAll(key int) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteAll(tree.Tree, intkey.IntKey(key))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) UpdateAll(key int, callBack UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateAll(tree.Tree, intkey.IntKey(key), wrapUpdateValueCallBack(callBack))
	return
}

func (tree *IntAVLTree) ReplaceAll(key int, value int) (ok bool) {
	tree.Tree, ok = avltree.ReplaceAll(tree.Tree, intkey.IntKey(key), value)
	return
}

func (tree *IntAVLTree) AlterAll(key int, callBack AlterNodeCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterAll(tree.Tree, intkey.IntKey(key), wrapAlterNodeCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) FindAll(key int) (nodes []Node) {
	return wrapNodes(avltree.FindAll(tree.Tree, intkey.IntKey(key)))
}

func (tree *IntAVLTree) MinAll() (nodea []Node) {
	return wrapNodes(avltree.MinAll(tree.Tree))
}

func (tree *IntAVLTree) MaxAll() (nodes []Node) {
	return wrapNodes(avltree.MaxAll(tree.Tree))
}

func (tree *IntAVLTree) DeleteIterate(callBack DeleteIterateCallBack) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteIterate(tree.Tree, false, wrapDeleteIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) DeleteIterateRev(callBack DeleteIterateCallBack) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteIterate(tree.Tree, true, wrapDeleteIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) DeleteRange(lower, upper int) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteRange(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) DeleteRangeRev(lower, upper int) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteRange(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) DeleteRangeIterate(lower, upper int, callBack DeleteIterateCallBack) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteRangeIterate(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapDeleteIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) DeleteRangeIterateRev(lower, upper int, callBack DeleteIterateCallBack) (deletedValues []KeyAndValue) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues = avltree.DeleteRangeIterate(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapDeleteIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) UpdateIterate(callBack UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateIterate(tree.Tree, false, wrapUpdateIterateCallBack(callBack))
	return
}

func (tree *IntAVLTree) UpdateIterateRev(callBack UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateIterate(tree.Tree, true, wrapUpdateIterateCallBack(callBack))
	return
}

func (tree *IntAVLTree) UpdateRange(lower, upper int, callBack UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRange(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapUpdateValueCallBack(callBack))
	return
}

func (tree *IntAVLTree) UpdateRangeRev(lower, upper int, callBack UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRange(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapUpdateValueCallBack(callBack))
	return
}

func (tree *IntAVLTree) UpdateRangeIterate(lower, upper int, callBack UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRangeIterate(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapUpdateIterateCallBack(callBack))
	return
}

func (tree *IntAVLTree) UpdateRangeIterateRev(lower, upper int, callBack UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRangeIterate(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapUpdateIterateCallBack(callBack))
	return
}

func (tree *IntAVLTree) ReplaceRange(lower, upper int, value int) (ok bool) {
	tree.Tree, ok = avltree.ReplaceRange(tree.Tree, intkey.IntKey(lower), intkey.IntKey(upper), value)
	return
}

func (tree *IntAVLTree) AlterIterate(callBack AlterIterateCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterIterate(tree.Tree, false, wrapAlterIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) AlterIterateRev(callBack AlterIterateCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterIterate(tree.Tree, true, wrapAlterIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) AlterRange(lower, upper int, callBack AlterNodeCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterRange(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapAlterNodeCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) AlterRangeRev(lower, upper int, callBack AlterNodeCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterRange(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapAlterNodeCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) AlterRangeIterate(lower, upper int, callBack AlterIterateCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterRangeIterate(tree.Tree, false, intkey.IntKey(lower), intkey.IntKey(upper), wrapAlterIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func (tree *IntAVLTree) AlterRangeIterateRev(lower, upper int, callBack AlterIterateCallBack) (deletedValues []KeyAndValue, ok bool) {
	var tempDeletedValues []avltree.KeyAndValue
	tree.Tree, tempDeletedValues, ok = avltree.AlterRangeIterate(tree.Tree, true, intkey.IntKey(lower), intkey.IntKey(upper), wrapAlterIterateCallBack(callBack))
	deletedValues = wrapKeyAndValues(tempDeletedValues)
	return
}

func wrapKeyAndValue(kv avltree.KeyAndValue) KeyAndValue {
	if kv == nil {
		return nil
	} else {
		return &keyAndValueWrapper{kv}
	}
}

func wrapKeyAndValues(kvs []avltree.KeyAndValue) []KeyAndValue {
	if kvs == nil {
		return nil
	} else {
		wrapped := make([]KeyAndValue, len(kvs))
		for i, kv := range kvs {
			wrapped[i] = &keyAndValueWrapper{kv}
		}
		return wrapped
	}
}

func wrapNode(node avltree.Node) Node {
	if node == nil {
		return nil
	} else {
		return &nodeWrapper{node}
	}
}

func wrapNodes(nodes []avltree.Node) []Node {
	if nodes == nil {
		return nil
	} else {
		wrapped := make([]Node, len(nodes))
		for i, node := range nodes {
			wrapped[i] = &nodeWrapper{node}
		}
		return wrapped
	}
}

func wrapIterateCallBack(callBack IterateCallBack) avltree.IterateCallBack {
	return func(node avltree.Node) (breakIteration bool) {
		return callBack(&nodeWrapper{node})
	}
}

func wrapUpdateValueCallBack(callBack UpdateValueCallBack) avltree.UpdateValueCallBack {
	return func(key avltree.Key, value interface{}) (newValue interface{}, keepOldValue bool) {
		newValue, keepOldValue = callBack(int(key.(intkey.IntKey)), value.(int))
		return
	}
}
func wrapUpdateIterateCallBack(callBack UpdateIterateCallBack) avltree.UpdateIterateCallBack {
	return func(key avltree.Key, oldValue interface{}) (newValue interface{}, keepOldValue, breakIteration bool) {
		newValue, keepOldValue, breakIteration = callBack(int(key.(intkey.IntKey)), oldValue.(int))
		return
	}
}

func wrapDeleteIterateCallBack(callBack DeleteIterateCallBack) avltree.DeleteIterateCallBack {
	return func(key avltree.Key, value interface{}) (deleteNode, breakIteration bool) {
		return callBack(int(key.(intkey.IntKey)), value.(int))
	}
}

func wrapAlterNodeCallBack(callBack AlterNodeCallBack) avltree.AlterNodeCallBack {
	return func(node avltree.AlterNode) (request avltree.AlterRequest) {
		return callBack(&alterNodeWrapper{node}).inner
	}
}

func wrapAlterIterateCallBack(callBack AlterIterateCallBack) avltree.AlterIterateCallBack {
	return func(node avltree.AlterNode) (request avltree.AlterRequest, breakIteration bool) {
		req, breakIteration := callBack(&alterNodeWrapper{node})
		return req.inner, breakIteration
	}
}

func (kv *keyAndValueWrapper) Key() int {
	return int(kv.inner.Key().(intkey.IntKey))
}

func (kv *keyAndValueWrapper) Value() int {
	return kv.inner.Value().(int)
}

func (node *nodeWrapper) Key() int {
	return int(node.inner.Key().(intkey.IntKey))
}

func (node *nodeWrapper) Value() int {
	return node.inner.Value().(int)
}

func (node *nodeWrapper) LeftChild() Node {
	return wrapNode(node.inner.LeftChild())
}

func (node *nodeWrapper) RightChild() Node {
	return wrapNode(node.inner.RightChild())
}

func (node *nodeWrapper) SetValue(newValue int) {
	node.inner.SetValue(newValue)
}

// 内部で保持している実際のノード(avltree.Node)へアクセスするためのメソッド(バックドア？)
func (node *nodeWrapper) Node() avltree.Node {
	return node.inner
}

func (node *alterNodeWrapper) Key() int {
	return int(node.inner.Key().(intkey.IntKey))
}

func (node *alterNodeWrapper) Value() int {
	return node.inner.Value().(int)
}

func (*alterNodeWrapper) Keep() (request AlterRequest) {
	return
}

func (*alterNodeWrapper) Replace(newValue int) (request AlterRequest) {
	request.inner.Replace(newValue)
	return
}

func (*alterNodeWrapper) Delete() (request AlterRequest) {
	request.inner.Delete()
	return
}

// AlterRequest内部で保持しているノードにアクセスするメソッド
func (node *alterNodeWrapper) Node() Node {
	if nodeGetter, ok := node.inner.(interface{ Node() avltree.Node }); ok {
		return wrapNode(nodeGetter.Node())
	} else {
		return nil
	}
}

func (request *AlterRequest) Keep() AlterRequest {
	request.inner.Keep()
	return *request
}

func (request *AlterRequest) Replace(newValue int) AlterRequest {
	request.inner.Replace(newValue)
	return *request
}

func (request *AlterRequest) Delete() AlterRequest {
	request.inner.Delete()
	return *request
}
