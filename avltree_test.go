package avltree

import (
	"runtime"
	"testing"
	"testing/quick"
)

type assertion testing.T

func (a *assertion) IsTrue(expectTrue bool, args ...interface{}) {
	if !expectTrue {
		msg := []interface{}{"Fail IsTrue:"}
		if _, file, line, ok := runtime.Caller(1); ok {
			msg = append(msg, "[", file, line, "]")
		}
		a.Fatal(append(msg, args...)...)
	}
}

func (a *assertion) IsFalse(expectFalse bool, args ...interface{}) {
	if expectFalse {
		msg := []interface{}{"Fail IsFalse:"}
		if _, file, line, ok := runtime.Caller(1); ok {
			msg = append(msg, "[", file, line, "]")
		}
		a.Fatal(append(msg, args...)...)
	}
}

func TestIntKey(t *testing.T) {
	f := func(k1, k2 int) bool {
		switch IntKey(k1).CompareTo(IntKey(k2)) {
		case LessThanOtherKey:
			return k1 < k2
		case EqualToOtherKey:
			return k1 == k2
		case GreaterThanOtherKey:
			return k1 > k2
		default:
			return false
		}
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestStringKey(t *testing.T) {
	f := func(k1, k2 string) bool {
		switch StringKey(k1).CompareTo(StringKey(k2)) {
		case LessThanOtherKey:
			return k1 < k2
		case EqualToOtherKey:
			return k1 == k2
		case GreaterThanOtherKey:
			return k1 > k2
		default:
			return false
		}
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestKeyOrdering_LessThanOtherKey(t *testing.T) {
	a := (*assertion)(t)
	a.IsTrue(
		LessThanOtherKey.LessThan(),
		"LessThanOtherKey.LessThan()",
	)
	a.IsTrue(
		LessThanOtherKey.LessThanOrEqualTo(),
		"LessThanOtherKey.LessThanOrEqualTo()",
	)
	a.IsFalse(
		LessThanOtherKey.EqualTo(),
		"LessThanOtherKey.EqualTo()",
	)
	a.IsTrue(
		LessThanOtherKey.NotEqualTo(),
		"LessThanOtherKey.NotEqualTo()",
	)
	a.IsFalse(
		LessThanOtherKey.GreaterThan(),
		"LessThanOtherKey.GreaterThan()",
	)
	a.IsFalse(
		LessThanOtherKey.GreaterThanOrEqualTo(),
		"LessThanOtherKey.GreaterThanOrEqualTo()",
	)
	a.IsTrue(
		LessThanOtherKey.Less(false),
		"LessThanOtherKey.Less(false)",
	)
	a.IsTrue(
		LessThanOtherKey.Less(true),
		"LessThanOtherKey.Less(true)",
	)
	a.IsFalse(
		LessThanOtherKey.Greater(false),
		"LessThanOtherKey.Greater(false)",
	)
	a.IsFalse(
		LessThanOtherKey.Greater(true),
		"LessThanOtherKey.Greater(true)",
	)
}

func TestKeyOrdering_EqualToOtherKey(t *testing.T) {
	a := (*assertion)(t)
	a.IsFalse(
		EqualToOtherKey.LessThan(),
		"EqualToOtherKey.LessThan()",
	)
	a.IsTrue(
		EqualToOtherKey.LessThanOrEqualTo(),
		"EqualToOtherKey.LessThanOrEqualTo()",
	)
	a.IsTrue(
		EqualToOtherKey.EqualTo(),
		"EqualToOtherKey.EqualTo()",
	)
	a.IsFalse(
		EqualToOtherKey.NotEqualTo(),
		"EqualToOtherKey.NotEqualTo()",
	)
	a.IsFalse(
		EqualToOtherKey.GreaterThan(),
		"EqualToOtherKey.GreaterThan()",
	)
	a.IsTrue(
		EqualToOtherKey.GreaterThanOrEqualTo(),
		"EqualToOtherKey.GreaterThanOrEqualTo()",
	)
	a.IsFalse(
		EqualToOtherKey.Less(false),
		"EqualToOtherKey.Less(false)",
	)
	a.IsTrue(
		EqualToOtherKey.Less(true),
		"EqualToOtherKey.Less(true)",
	)
	a.IsFalse(
		EqualToOtherKey.Greater(false),
		"EqualToOtherKey.Greater(false)",
	)
	a.IsTrue(
		EqualToOtherKey.Greater(true),
		"EqualToOtherKey.Greater(true)",
	)
}

func TestKeyOrdering_GreaterThanOtherKey(t *testing.T) {
	a := (*assertion)(t)
	a.IsFalse(
		GreaterThanOtherKey.LessThan(),
		"GreaterThanOtherKey.LessThan()",
	)
	a.IsFalse(
		GreaterThanOtherKey.LessThanOrEqualTo(),
		"GreaterThanOtherKey.LessThanOrEqualTo()",
	)
	a.IsFalse(
		GreaterThanOtherKey.EqualTo(),
		"GreaterThanOtherKey.EqualTo()",
	)
	a.IsTrue(
		GreaterThanOtherKey.NotEqualTo(),
		"GreaterThanOtherKey.NotEqualTo()",
	)
	a.IsTrue(
		GreaterThanOtherKey.GreaterThan(),
		"GreaterThanOtherKey.GreaterThan()",
	)
	a.IsTrue(
		GreaterThanOtherKey.GreaterThanOrEqualTo(),
		"GreaterThanOtherKey.GreaterThanOrEqualTo()",
	)
	a.IsFalse(
		GreaterThanOtherKey.Less(false),
		"GreaterThanOtherKey.Less(false)",
	)
	a.IsFalse(
		GreaterThanOtherKey.Less(true),
		"GreaterThanOtherKey.Less(true)",
	)
	a.IsTrue(
		GreaterThanOtherKey.Greater(false),
		"GreaterThanOtherKey.Greater(false)",
	)
	a.IsTrue(
		GreaterThanOtherKey.Greater(true),
		"GreaterThanOtherKey.Greater(true)",
	)
}

func TestInsert(t *testing.T) {
	_ = Insert
	t.Log("avltree.Insert のテストをまだ実装していない")
}

func TestDelete(t *testing.T) {
	_ = Delete
	t.Log("avltree.Delete のテストをまだ実装していない")
}

func TestUpdate(t *testing.T) {
	_ = Update
	t.Log("avltree.Update のテストをまだ実装していない")
}

func TestReplace(t *testing.T) {
	_ = Replace
	t.Log("avltree.Replace のテストをまだ実装していない")
}

func TestAlter(t *testing.T) {
	_ = Alter
	t.Log("avltree.Alter のテストをまだ実装していない")
}

func TestClear(t *testing.T) {
	_ = Clear
	t.Log("avltree.Clear のテストをまだ実装していない")
}

func TestRelease(t *testing.T) {
	_ = Release
	t.Log("avltree.Release のテストをまだ実装していない")
}

func TestFind(t *testing.T) {
	_ = Find
	t.Log("avltree.Find のテストをまだ実装していない")
}

func TestIterate(t *testing.T) {
	_ = Iterate
	t.Log("avltree.Iterate のテストをまだ実装していない")
}

func TestRange(t *testing.T) {
	_ = Range
	t.Log("avltree.Range のテストをまだ実装していない")
}

func TestRangeIterate(t *testing.T) {
	_ = RangeIterate
	t.Log("avltree.RangeIterate のテストをまだ実装していない")
}

func TestCount(t *testing.T) {
	_ = Count
	t.Log("avltree.Count のテストをまだ実装していない")
}

func TestCountRange(t *testing.T) {
	_ = CountRange
	t.Log("avltree.CountRange のテストをまだ実装していない")
}

func TestMin(t *testing.T) {
	_ = Min
	t.Log("avltree.Min のテストをまだ実装していない")
}

func TestMax(t *testing.T) {
	_ = Max
	t.Log("avltree.Max のテストをまだ実装していない")
}

func TestDeleteAll(t *testing.T) {
	_ = DeleteAll
	t.Log("avltree.DeleteAll のテストをまだ実装していない")
}

func TestUpdateAll(t *testing.T) {
	_ = UpdateAll
	t.Log("avltree.UpdateAll のテストをまだ実装していない")
}

func TestReplaceAll(t *testing.T) {
	_ = ReplaceAll
	t.Log("avltree.ReplaceAll のテストをまだ実装していない")
}

func TestAlterAll(t *testing.T) {
	_ = AlterAll
	t.Log("avltree.AlterAll のテストをまだ実装していない")
}

func TestFindAll(t *testing.T) {
	_ = FindAll
	t.Log("avltree.FindAll のテストをまだ実装していない")
}

func TestMinAll(t *testing.T) {
	_ = MinAll
	t.Log("avltree.MinAll のテストをまだ実装していない")
}

func TestMaxAll(t *testing.T) {
	_ = MaxAll
	t.Log("avltree.MaxAll のテストをまだ実装していない")
}

func TestDeleteIterate(t *testing.T) {
	_ = DeleteIterate
	t.Log("avltree.DeleteIterate のテストをまだ実装していない")
}

func TestDeleteRange(t *testing.T) {
	_ = DeleteRange
	t.Log("avltree.DeleteRange のテストをまだ実装していない")
}

func TestDeleteRangeIterate(t *testing.T) {
	_ = DeleteRangeIterate
	t.Log("avltree.DeleteRangeIterate のテストをまだ実装していない")
}

func TestUpdateIterate(t *testing.T) {
	_ = UpdateIterate
	t.Log("avltree.UpdateIterate のテストをまだ実装していない")
}

func TestUpdateRange(t *testing.T) {
	_ = UpdateRange
	t.Log("avltree.UpdateRange のテストをまだ実装していない")
}

func TestUpdateRangeIterate(t *testing.T) {
	_ = UpdateRangeIterate
	t.Log("avltree.UpdateRangeIterate のテストをまだ実装していない")
}

func TestReplaceRange(t *testing.T) {
	_ = ReplaceRange
	t.Log("avltree.ReplaceRange のテストをまだ実装していない")
}

func TestAlterIterate(t *testing.T) {
	_ = AlterIterate
	t.Log("avltree.AlterIterate のテストをまだ実装していない")
}

func TestAlterRange(t *testing.T) {
	_ = AlterRange
	t.Log("avltree.AlterRange のテストをまだ実装していない")
}

func TestAlterRangeIterate(t *testing.T) {
	_ = AlterRangeIterate
	t.Log("avltree.AlterRangeIterate のテストをまだ実装していない")
}

func TestInner_countNode(t *testing.T) {
	_ = countNode
	t.Log("countNode のテストをまだ実装していない")
}

func TestInner_countRange(t *testing.T) {
	_ = countRange
	t.Log("countRange のテストをまだ実装していない")
}

func TestInner_countExtendedRange(t *testing.T) {
	_ = countExtendedRange
	t.Log("countExtendedRange のテストをまだ実装していない")
}

func TestInner_getHeight(t *testing.T) {
	_ = getHeight
	t.Log("getHeight のテストをまだ実装していない")
}

func TestInner_checkBalance(t *testing.T) {
	_ = checkBalance
	t.Log("checkBalance のテストをまだ実装していない")
}

func TestInner_compareNodeHeight(t *testing.T) {
	_ = compareNodeHeight
	t.Log("compareNodeHeight のテストをまだ実装していない")
}

func TestInner_compareChildHeight(t *testing.T) {
	_ = compareChildHeight
	t.Log("compareChildHeight のテストをまだ実装していない")
}

func TestInner_intMax(t *testing.T) {
	_ = intMax
	t.Log("intMax のテストをまだ実装していない")
}

func TestInner_calcNewHeight(t *testing.T) {
	_ = calcNewHeight
	t.Log("calcNewHeight のテストをまだ実装していない")
}

func TestInner_setChildren(t *testing.T) {
	_ = setChildren
	t.Log("setChildren のテストをまだ実装していない")
}

func TestInner_setLeftChild(t *testing.T) {
	_ = setLeftChild
	t.Log("setLeftChild のテストをまだ実装していない")
}

func TestInner_setRightChild(t *testing.T) {
	_ = setRightChild
	t.Log("setRightChild のテストをまだ実装していない")
}

func TestInner_resetNode(t *testing.T) {
	_ = resetNode
	t.Log("resetNode のテストをまだ実装していない")
}

func TestInner_insertHelper_newNode(t *testing.T) {
	var helper insertHelper
	_ = helper.newNode
	t.Log("insertHelper.newNode のテストをまだ実装していない")
}

func TestInner_insertHelper_compareKey(t *testing.T) {
	var helper insertHelper
	_ = helper.compareKey
	t.Log("insertHelper.compareKey のテストをまだ実装していない")
}

func TestInner_insertHelper_allowDuplicateKeys(t *testing.T) {
	var helper insertHelper
	_ = helper.allowDuplicateKeys
	t.Log("insertHelper.allowDuplicateKeys のテストをまだ実装していない")
}

func TestInner_insertHelper_insertTo(t *testing.T) {
	var helper insertHelper
	_ = helper.insertTo
	t.Log("insertHelper.insertTo のテストをまだ実装していない")
}

func TestInner_rotate(t *testing.T) {
	_ = rotate
	t.Log("rotate のテストをまだ実装していない")
}

func TestInner_rotateRight(t *testing.T) {
	_ = rotateRight
	t.Log("rotateRight のテストをまだ実装していない")
}

func TestInner_rotateLeft(t *testing.T) {
	_ = rotateLeft
	t.Log("rotateLeft のテストをまだ実装していない")
}

func TestInner_removeNode(t *testing.T) {
	_ = removeNode
	t.Log("removeNode のテストをまだ実装していない")
}

func TestInner_removeMin(t *testing.T) {
	_ = removeMin
	t.Log("removeMin のテストをまだ実装していない")
}

func TestInner_removeMax(t *testing.T) {
	_ = removeMax
	t.Log("removeMax のテストをまだ実装していない")
}

func TestInner_ascIterateNode(t *testing.T) {
	_ = ascIterateNode
	t.Log("ascIterateNode のテストをまだ実装していない")
}

func TestInner_descIterateNode(t *testing.T) {
	_ = descIterateNode
	t.Log("descIterateNode のテストをまだ実装していない")
}

func TestInner_newKeyBounds(t *testing.T) {
	_ = newKeyBounds
	t.Log("newKeyBounds のテストをまだ実装していない")
}

func TestInner_upperBound(t *testing.T) {
	var _ upperBound
	t.Log("upperBound のテストをまだ実装していない")
}

func TestInner_lowerBound(t *testing.T) {
	var _ lowerBound
	t.Log("lowerBound のテストをまだ実装していない")
}

func TestInner_bothBounds(t *testing.T) {
	var _ bothBounds
	t.Log("bothBounds のテストをまだ実装していない")
}

func TestInner_noBoundsChecker(t *testing.T) {
	var _ noBoundsChecker
	t.Log("noBoundsChecker のテストをまだ実装していない")
}

func TestInner_upperBoundsChecker(t *testing.T) {
	var _ upperBoundsChecker
	t.Log("upperBoundsChecker のテストをまだ実装していない")
}

func TestInner_lowerBoundsChecker(t *testing.T) {
	var _ lowerBoundsChecker
	t.Log("lowerBoundsChecker のテストをまだ実装していない")
}

func TestInner_ascRangeNode(t *testing.T) {
	_ = ascRangeNode
	t.Log("ascRangeNode のテストをまだ実装していない")
}

func TestInner_descRangeNode(t *testing.T) {
	_ = descRangeNode
	t.Log("descRangeNode のテストをまだ実装していない")
}

func TestInner_updateValue(t *testing.T) {
	_ = updateValue
	t.Log("updateValue のテストをまだ実装していない")
}

func TestInner_ascUpdateIterate(t *testing.T) {
	_ = ascUpdateIterate
	t.Log("ascUpdateIterate のテストをまだ実装していない")
}

func TestInner_descUpdateIterate(t *testing.T) {
	_ = descUpdateIterate
	t.Log("descUpdateIterate のテストをまだ実装していない")
}

func TestInner_ascUpdateRange(t *testing.T) {
	_ = ascUpdateRange
	t.Log("ascUpdateRange のテストをまだ実装していない")
}

func TestInner_descUpdateRange(t *testing.T) {
	_ = descUpdateRange
	t.Log("descUpdateRange のテストをまだ実装していない")
}

func TestInner_ascDeleteIterate(t *testing.T) {
	_ = ascDeleteIterate
	t.Log("ascDeleteIterate のテストをまだ実装していない")
}

func TestInner_descDeleteIterate(t *testing.T) {
	_ = descDeleteIterate
	t.Log("descDeleteIterate のテストをまだ実装していない")
}

func TestInner_ascDeleteRange(t *testing.T) {
	_ = ascDeleteRange
	t.Log("ascDeleteRange のテストをまだ実装していない")
}

func TestInner_descDeleteRange(t *testing.T) {
	_ = descDeleteRange
	t.Log("descDeleteRange のテストをまだ実装していない")
}

func TestInner_alterNode(t *testing.T) {
	var _ alterNode
	t.Log("alterNode のテストをまだ実装していない")
}

func TestAlterRequest(t *testing.T) {
	var _ AlterRequest
	t.Log("AlterRequest のテストをまだ実装していない")
}

func TestInner_alter(t *testing.T) {
	_ = alter
	t.Log("alter のテストをまだ実装していない")
}

func TestInner_ascAlterIterat(t *testing.T) {
	_ = ascAlterIterate
	t.Log("ascAlterIterate のテストをまだ実装していない")
}

func TestInner_descAlterIterat(t *testing.T) {
	_ = descAlterIterate
	t.Log("descAlterIterate のテストをまだ実装していない")
}

func TestInner_ascAlterRange(t *testing.T) {
	_ = ascAlterRange
	t.Log("ascAlterRange のテストをまだ実装していない")
}

func TestInner_descAlterRange(t *testing.T) {
	_ = descAlterRange
	t.Log("descAlterRange のテストをまだ実装していない")
}
