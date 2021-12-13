package avltree

import (
	"runtime"
	"testing"
	"testing/quick"
)

var cfg1000 = &quick.Config{MaxCount: 1000}

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
		var key1 Key = IntKey(k1)
		var key2 Key = IntKey(k2)
		switch key1.CompareTo(key2) {
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

	if err := quick.Check(f, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestStringKey(t *testing.T) {
	f := func(k1, k2 string) bool {
		var key1 Key = StringKey(k1)
		var key2 Key = StringKey(k2)
		switch key1.CompareTo(key2) {
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

	if err := quick.Check(f, cfg1000); err != nil {
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
	t.Skip("avltree.Insert のテストをまだ実装していない")
}

func TestDelete(t *testing.T) {
	_ = Delete
	t.Skip("avltree.Delete のテストをまだ実装していない")
}

func TestUpdate(t *testing.T) {
	_ = Update
	t.Skip("avltree.Update のテストをまだ実装していない")
}

func TestReplace(t *testing.T) {
	_ = Replace
	t.Skip("avltree.Replace のテストをまだ実装していない")
}

func TestAlter(t *testing.T) {
	_ = Alter
	t.Skip("avltree.Alter のテストをまだ実装していない")
}

func TestClear(t *testing.T) {
	_ = Clear
	t.Skip("avltree.Clear のテストをまだ実装していない")
}

func TestRelease(t *testing.T) {
	_ = Release
	t.Skip("avltree.Release のテストをまだ実装していない")
}

func TestFind(t *testing.T) {
	_ = Find
	t.Skip("avltree.Find のテストをまだ実装していない")
}

func TestIterate(t *testing.T) {
	_ = Iterate
	t.Skip("avltree.Iterate のテストをまだ実装していない")
}

func TestRange(t *testing.T) {
	_ = Range
	t.Skip("avltree.Range のテストをまだ実装していない")
}

func TestRangeIterate(t *testing.T) {
	_ = RangeIterate
	t.Skip("avltree.RangeIterate のテストをまだ実装していない")
}

func TestCount(t *testing.T) {
	_ = Count
	t.Skip("avltree.Count のテストをまだ実装していない")
}

func TestCountRange(t *testing.T) {
	_ = CountRange
	t.Skip("avltree.CountRange のテストをまだ実装していない")
}

func TestMin(t *testing.T) {
	_ = Min
	t.Skip("avltree.Min のテストをまだ実装していない")
}

func TestMax(t *testing.T) {
	_ = Max
	t.Skip("avltree.Max のテストをまだ実装していない")
}

func TestDeleteAll(t *testing.T) {
	_ = DeleteAll
	t.Skip("avltree.DeleteAll のテストをまだ実装していない")
}

func TestUpdateAll(t *testing.T) {
	_ = UpdateAll
	t.Skip("avltree.UpdateAll のテストをまだ実装していない")
}

func TestReplaceAll(t *testing.T) {
	_ = ReplaceAll
	t.Skip("avltree.ReplaceAll のテストをまだ実装していない")
}

func TestAlterAll(t *testing.T) {
	_ = AlterAll
	t.Skip("avltree.AlterAll のテストをまだ実装していない")
}

func TestFindAll(t *testing.T) {
	_ = FindAll
	t.Skip("avltree.FindAll のテストをまだ実装していない")
}

func TestMinAll(t *testing.T) {
	_ = MinAll
	t.Skip("avltree.MinAll のテストをまだ実装していない")
}

func TestMaxAll(t *testing.T) {
	_ = MaxAll
	t.Skip("avltree.MaxAll のテストをまだ実装していない")
}

func TestDeleteIterate(t *testing.T) {
	_ = DeleteIterate
	t.Skip("avltree.DeleteIterate のテストをまだ実装していない")
}

func TestDeleteRange(t *testing.T) {
	_ = DeleteRange
	t.Skip("avltree.DeleteRange のテストをまだ実装していない")
}

func TestDeleteRangeIterate(t *testing.T) {
	_ = DeleteRangeIterate
	t.Skip("avltree.DeleteRangeIterate のテストをまだ実装していない")
}

func TestUpdateIterate(t *testing.T) {
	_ = UpdateIterate
	t.Skip("avltree.UpdateIterate のテストをまだ実装していない")
}

func TestUpdateRange(t *testing.T) {
	_ = UpdateRange
	t.Skip("avltree.UpdateRange のテストをまだ実装していない")
}

func TestUpdateRangeIterate(t *testing.T) {
	_ = UpdateRangeIterate
	t.Skip("avltree.UpdateRangeIterate のテストをまだ実装していない")
}

func TestReplaceRange(t *testing.T) {
	_ = ReplaceRange
	t.Skip("avltree.ReplaceRange のテストをまだ実装していない")
}

func TestAlterIterate(t *testing.T) {
	_ = AlterIterate
	t.Skip("avltree.AlterIterate のテストをまだ実装していない")
}

func TestAlterRange(t *testing.T) {
	_ = AlterRange
	t.Skip("avltree.AlterRange のテストをまだ実装していない")
}

func TestAlterRangeIterate(t *testing.T) {
	_ = AlterRangeIterate
	t.Skip("avltree.AlterRangeIterate のテストをまだ実装していない")
}

func TestInner_countNode(t *testing.T) {
	_ = countNode
	t.Skip("countNode のテストをまだ実装していない")
}

func TestInner_countRange(t *testing.T) {
	_ = countRange
	t.Skip("countRange のテストをまだ実装していない")
}

func TestInner_countExtendedRange(t *testing.T) {
	_ = countExtendedRange
	t.Skip("countExtendedRange のテストをまだ実装していない")
}

func TestInner_getHeight(t *testing.T) {
	_ = getHeight
	t.Skip("getHeight のテストをまだ実装していない")
}

func TestInner_checkBalance(t *testing.T) {
	_ = checkBalance
	t.Skip("checkBalance のテストをまだ実装していない")
}

func TestInner_compareNodeHeight(t *testing.T) {
	_ = compareNodeHeight
	t.Skip("compareNodeHeight のテストをまだ実装していない")
}

func TestInner_compareChildHeight(t *testing.T) {
	_ = compareChildHeight
	t.Skip("compareChildHeight のテストをまだ実装していない")
}

func TestInner_intMax(t *testing.T) {
	f := func(a, b int) bool {
		m := intMax(a, b)
		return (m >= a) && (m >= b) && (m == a || m == b)
	}

	if err := quick.Check(f, cfg1000); err != nil {
		t.Fatal(err)
	}
}

func TestInner_calcNewHeight(t *testing.T) {
	_ = calcNewHeight
	t.Skip("calcNewHeight のテストをまだ実装していない")
}

func TestInner_setChildren(t *testing.T) {
	_ = setChildren
	t.Skip("setChildren のテストをまだ実装していない")
}

func TestInner_setLeftChild(t *testing.T) {
	_ = setLeftChild
	t.Skip("setLeftChild のテストをまだ実装していない")
}

func TestInner_setRightChild(t *testing.T) {
	_ = setRightChild
	t.Skip("setRightChild のテストをまだ実装していない")
}

func TestInner_resetNode(t *testing.T) {
	_ = resetNode
	t.Skip("resetNode のテストをまだ実装していない")
}

func TestInner_insertHelper_newNode(t *testing.T) {
	var helper insertHelper
	_ = helper.newNode
	t.Skip("insertHelper.newNode のテストをまだ実装していない")
}

func TestInner_insertHelper_compareKey(t *testing.T) {
	var helper insertHelper
	_ = helper.compareKey
	t.Skip("insertHelper.compareKey のテストをまだ実装していない")
}

func TestInner_insertHelper_allowDuplicateKeys(t *testing.T) {
	var helper insertHelper
	_ = helper.allowDuplicateKeys
	t.Skip("insertHelper.allowDuplicateKeys のテストをまだ実装していない")
}

func TestInner_insertHelper_insertTo(t *testing.T) {
	var helper insertHelper
	_ = helper.insertTo
	t.Skip("insertHelper.insertTo のテストをまだ実装していない")
}

func TestInner_rotate(t *testing.T) {
	_ = rotate
	t.Skip("rotate のテストをまだ実装していない")
}

func TestInner_rotateRight(t *testing.T) {
	_ = rotateRight
	t.Skip("rotateRight のテストをまだ実装していない")
}

func TestInner_rotateLeft(t *testing.T) {
	_ = rotateLeft
	t.Skip("rotateLeft のテストをまだ実装していない")
}

func TestInner_removeNode(t *testing.T) {
	_ = removeNode
	t.Skip("removeNode のテストをまだ実装していない")
}

func TestInner_removeMin(t *testing.T) {
	_ = removeMin
	t.Skip("removeMin のテストをまだ実装していない")
}

func TestInner_removeMax(t *testing.T) {
	_ = removeMax
	t.Skip("removeMax のテストをまだ実装していない")
}

func TestInner_ascIterateNode(t *testing.T) {
	_ = ascIterateNode
	t.Skip("ascIterateNode のテストをまだ実装していない")
}

func TestInner_descIterateNode(t *testing.T) {
	_ = descIterateNode
	t.Skip("descIterateNode のテストをまだ実装していない")
}

func TestInner_newKeyBounds(t *testing.T) {
	_ = newKeyBounds
	t.Skip("newKeyBounds のテストをまだ実装していない")
}

func TestInner_upperBound(t *testing.T) {
	var _ upperBound
	t.Skip("upperBound のテストをまだ実装していない")
}

func TestInner_lowerBound(t *testing.T) {
	var _ lowerBound
	t.Skip("lowerBound のテストをまだ実装していない")
}

func TestInner_bothBounds(t *testing.T) {
	var _ bothBounds
	t.Skip("bothBounds のテストをまだ実装していない")
}

func TestInner_noBoundsChecker(t *testing.T) {
	a := (*assertion)(t)
	var nbc noBoundsChecker
	var checker boundsChecker = nbc
	a.IsTrue(
		checker.includeLower(),
		"checker.includeLower()",
	)
	a.IsTrue(
		checker.includeKey(),
		"checker.includeKey()",
	)
	a.IsTrue(
		checker.includeUpper(),
		"checker.includeUpper()",
	)
}

func TestInner_upperBoundsChecker(t *testing.T) {
	var _ upperBoundsChecker
	t.Skip("upperBoundsChecker のテストをまだ実装していない")
}

func TestInner_lowerBoundsChecker(t *testing.T) {
	var _ lowerBoundsChecker
	t.Skip("lowerBoundsChecker のテストをまだ実装していない")
}

func TestInner_ascRangeNode(t *testing.T) {
	_ = ascRangeNode
	t.Skip("ascRangeNode のテストをまだ実装していない")
}

func TestInner_descRangeNode(t *testing.T) {
	_ = descRangeNode
	t.Skip("descRangeNode のテストをまだ実装していない")
}

func TestInner_updateValue(t *testing.T) {
	_ = updateValue
	t.Skip("updateValue のテストをまだ実装していない")
}

func TestInner_ascUpdateIterate(t *testing.T) {
	_ = ascUpdateIterate
	t.Skip("ascUpdateIterate のテストをまだ実装していない")
}

func TestInner_descUpdateIterate(t *testing.T) {
	_ = descUpdateIterate
	t.Skip("descUpdateIterate のテストをまだ実装していない")
}

func TestInner_ascUpdateRange(t *testing.T) {
	_ = ascUpdateRange
	t.Skip("ascUpdateRange のテストをまだ実装していない")
}

func TestInner_descUpdateRange(t *testing.T) {
	_ = descUpdateRange
	t.Skip("descUpdateRange のテストをまだ実装していない")
}

func TestInner_ascDeleteIterate(t *testing.T) {
	_ = ascDeleteIterate
	t.Skip("ascDeleteIterate のテストをまだ実装していない")
}

func TestInner_descDeleteIterate(t *testing.T) {
	_ = descDeleteIterate
	t.Skip("descDeleteIterate のテストをまだ実装していない")
}

func TestInner_ascDeleteRange(t *testing.T) {
	_ = ascDeleteRange
	t.Skip("ascDeleteRange のテストをまだ実装していない")
}

func TestInner_descDeleteRange(t *testing.T) {
	_ = descDeleteRange
	t.Skip("descDeleteRange のテストをまだ実装していない")
}

func TestInner_alterNode(t *testing.T) {
	var _ alterNode
	t.Skip("alterNode のテストをまだ実装していない")
}

func TestAlterRequest(t *testing.T) {
	var _ AlterRequest
	t.Skip("AlterRequest のテストをまだ実装していない")
}

func TestInner_alter(t *testing.T) {
	_ = alter
	t.Skip("alter のテストをまだ実装していない")
}

func TestInner_ascAlterIterat(t *testing.T) {
	_ = ascAlterIterate
	t.Skip("ascAlterIterate のテストをまだ実装していない")
}

func TestInner_descAlterIterat(t *testing.T) {
	_ = descAlterIterate
	t.Skip("descAlterIterate のテストをまだ実装していない")
}

func TestInner_ascAlterRange(t *testing.T) {
	_ = ascAlterRange
	t.Skip("ascAlterRange のテストをまだ実装していない")
}

func TestInner_descAlterRange(t *testing.T) {
	_ = descAlterRange
	t.Skip("descAlterRange のテストをまだ実装していない")
}
