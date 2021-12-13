package simplewrapper

import "github.com/neetsdkasu/avltree"

type AVLTree struct {
	Tree avltree.Tree
}

func New(tree avltree.Tree) *AVLTree {
	return &AVLTree{tree}
}

func (tree *AVLTree) Insert(key avltree.Key, value interface{}) (ok bool) {
	tree.Tree, ok = avltree.Insert(tree.Tree, false, key, value)
	return
}

func (tree *AVLTree) InsertOrReplace(key avltree.Key, value interface{}) (ok bool) {
	tree.Tree, ok = avltree.Insert(tree.Tree, true, key, value)
	return
}

func (tree *AVLTree) Delete(key avltree.Key) (deletedValue avltree.KeyAndValue) {
	tree.Tree, deletedValue = avltree.Delete(tree.Tree, key)
	return
}

func (tree *AVLTree) Update(key avltree.Key, callBack avltree.UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.Update(tree.Tree, key, callBack)
	return
}

func (tree *AVLTree) Replace(key avltree.Key, value interface{}) (ok bool) {
	tree.Tree, ok = avltree.Replace(tree.Tree, key, value)
	return
}

func (tree *AVLTree) Alter(key avltree.Key, callBack avltree.AlterNodeCallBack) (deletedValue avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValue, ok = avltree.Alter(tree.Tree, key, callBack)
	return
}

func (tree *AVLTree) Clear() {
	tree.Tree = avltree.Clear(tree.Tree)
}

func (tree *AVLTree) Release() {
	avltree.Release(&tree.Tree)
}

func (tree *AVLTree) Find(key avltree.Key) (node avltree.Node) {
	return avltree.Find(tree.Tree, key)
}

func (tree *AVLTree) Iterate(callBack avltree.IterateCallBack) {
	avltree.Iterate(tree.Tree, false, callBack)
}

func (tree *AVLTree) IterateRev(callBack avltree.IterateCallBack) {
	avltree.Iterate(tree.Tree, true, callBack)
}

func (tree *AVLTree) Range(lower, upper avltree.Key) (nodes []avltree.Node) {
	return avltree.Range(tree.Tree, false, lower, upper)
}

func (tree *AVLTree) RangeRev(lower, upper avltree.Key) (nodes []avltree.Node) {
	return avltree.Range(tree.Tree, true, lower, upper)
}

func (tree *AVLTree) RangeIterate(lower, upper avltree.Key, callBack avltree.IterateCallBack) {
	avltree.RangeIterate(tree.Tree, false, lower, upper, callBack)
}

func (tree *AVLTree) RangeIterateRev(lower, upper avltree.Key, callBack avltree.IterateCallBack) {
	avltree.RangeIterate(tree.Tree, true, lower, upper, callBack)
}

func (tree *AVLTree) Count() int {
	return avltree.Count(tree.Tree)
}

func (tree *AVLTree) CountRange(lower, upper avltree.Key) int {
	return avltree.CountRange(tree.Tree, lower, upper)
}

func (tree *AVLTree) Min() (node avltree.Node) {
	return avltree.Min(tree.Tree)
}

func (tree *AVLTree) Max() (node avltree.Node) {
	return avltree.Max(tree.Tree)
}

func (tree *AVLTree) DeleteAll(key avltree.Key) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteAll(tree.Tree, key)
	return
}

func (tree *AVLTree) UpdateAll(key avltree.Key, callBack avltree.UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateAll(tree.Tree, key, callBack)
	return
}

func (tree *AVLTree) ReplaceAll(key avltree.Key, value interface{}) (ok bool) {
	tree.Tree, ok = avltree.ReplaceAll(tree.Tree, key, value)
	return
}

func (tree *AVLTree) AlterAll(key avltree.Key, callBack avltree.AlterNodeCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterAll(tree.Tree, key, callBack)
	return
}

func (tree *AVLTree) FindAll(key avltree.Key) (nodes []avltree.Node) {
	return avltree.FindAll(tree.Tree, key)
}

func (tree *AVLTree) MinAll() (nodea []avltree.Node) {
	return avltree.MinAll(tree.Tree)
}

func (tree *AVLTree) MaxAll() (nodes []avltree.Node) {
	return avltree.MaxAll(tree.Tree)
}

func (tree *AVLTree) DeleteIterate(callBack avltree.DeleteIterateCallBack) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteIterate(tree.Tree, false, callBack)
	return
}

func (tree *AVLTree) DeleteIterateRev(callBack avltree.DeleteIterateCallBack) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteIterate(tree.Tree, true, callBack)
	return
}

func (tree *AVLTree) DeleteRange(lower, upper avltree.Key) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteRange(tree.Tree, false, lower, upper)
	return
}

func (tree *AVLTree) DeleteRangeRev(lower, upper avltree.Key) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteRange(tree.Tree, true, lower, upper)
	return
}

func (tree *AVLTree) DeleteRangeIterate(lower, upper avltree.Key, callBack avltree.DeleteIterateCallBack) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteRangeIterate(tree.Tree, false, lower, upper, callBack)
	return
}

func (tree *AVLTree) DeleteRangeIterateRev(lower, upper avltree.Key, callBack avltree.DeleteIterateCallBack) (deletedValues []avltree.KeyAndValue) {
	tree.Tree, deletedValues = avltree.DeleteRangeIterate(tree.Tree, true, lower, upper, callBack)
	return
}

func (tree *AVLTree) UpdateIterate(callBack avltree.UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateIterate(tree.Tree, false, callBack)
	return
}

func (tree *AVLTree) UpdateIterateRev(callBack avltree.UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateIterate(tree.Tree, true, callBack)
	return
}

func (tree *AVLTree) UpdateRange(lower, upper avltree.Key, callBack avltree.UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRange(tree.Tree, false, lower, upper, callBack)
	return
}

func (tree *AVLTree) UpdateRangeRev(lower, upper avltree.Key, callBack avltree.UpdateValueCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRange(tree.Tree, true, lower, upper, callBack)
	return
}

func (tree *AVLTree) UpdateRangeIterate(lower, upper avltree.Key, callBack avltree.UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRangeIterate(tree.Tree, false, lower, upper, callBack)
	return
}

func (tree *AVLTree) UpdateRangeIterateRev(lower, upper avltree.Key, callBack avltree.UpdateIterateCallBack) (ok bool) {
	tree.Tree, ok = avltree.UpdateRangeIterate(tree.Tree, true, lower, upper, callBack)
	return
}

func (tree *AVLTree) ReplaceRange(lower, upper avltree.Key, value interface{}) (ok bool) {
	tree.Tree, ok = avltree.ReplaceRange(tree.Tree, lower, upper, value)
	return
}

func (tree *AVLTree) AlterIterate(callBack avltree.AlterIterateCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterIterate(tree.Tree, false, callBack)
	return
}

func (tree *AVLTree) AlterIterateRev(callBack avltree.AlterIterateCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterIterate(tree.Tree, true, callBack)
	return
}

func (tree *AVLTree) AlterRange(lower, upper avltree.Key, callBack avltree.AlterNodeCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterRange(tree.Tree, false, lower, upper, callBack)
	return
}

func (tree *AVLTree) AlterRangeRev(lower, upper avltree.Key, callBack avltree.AlterNodeCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterRange(tree.Tree, true, lower, upper, callBack)
	return
}

func (tree *AVLTree) AlterRangeIterate(lower, upper avltree.Key, callBack avltree.AlterIterateCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterRangeIterate(tree.Tree, false, lower, upper, callBack)
	return
}

func (tree *AVLTree) AlterRangeIterateRev(lower, upper avltree.Key, callBack avltree.AlterIterateCallBack) (deletedValues []avltree.KeyAndValue, ok bool) {
	tree.Tree, deletedValues, ok = avltree.AlterRangeIterate(tree.Tree, true, lower, upper, callBack)
	return
}
