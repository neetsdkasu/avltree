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
	t.Skip("avltree.Insert のテストは未実装")
}

func TestDelete(t *testing.T) {
	_ = Delete
	t.Skip("avltree.Delete のテストは未実装")
}

func TestUpdate(t *testing.T) {
	_ = Update
	t.Skip("avltree.Update のテストは未実装")
}

func TestReplace(t *testing.T) {
	_ = Replace
	t.Skip("avltree.Replace のテストは未実装")
}

func TestAlter(t *testing.T) {
	_ = Alter
	t.Skip("avltree.Alter のテストは未実装")
}

func TestClear(t *testing.T) {
	_ = Clear
	t.Skip("avltree.Clear のテストは未実装")
}

func TestRelease(t *testing.T) {
	_ = Release
	t.Skip("avltree.Release のテストは未実装")
}

func TestFind(t *testing.T) {
	_ = Find
	t.Skip("avltree.Find のテストは未実装")
}

func TestIterate(t *testing.T) {
	_ = Iterate
	t.Skip("avltree.Iterate のテストは未実装")
}

func TestRange(t *testing.T) {
	_ = Range
	t.Skip("avltree.Range のテストは未実装")
}

func TestRangeIterate(t *testing.T) {
	_ = RangeIterate
	t.Skip("avltree.RangeIterate のテストは未実装")
}

func TestCount(t *testing.T) {
	_ = Count
	t.Skip("avltree.Count のテストは未実装")
}

func TestCountRange(t *testing.T) {
	_ = CountRange
	t.Skip("avltree.CountRange のテストは未実装")
}

func TestMin(t *testing.T) {
	_ = Min
	t.Skip("avltree.Min のテストは未実装")
}

func TestMax(t *testing.T) {
	_ = Max
	t.Skip("avltree.Max のテストは未実装")
}

func TestDeleteAll(t *testing.T) {
	_ = DeleteAll
	t.Skip("avltree.DeleteAll のテストは未実装")
}

func TestUpdateAll(t *testing.T) {
	_ = UpdateAll
	t.Skip("avltree.UpdateAll のテストは未実装")
}

func TestReplaceAll(t *testing.T) {
	_ = ReplaceAll
	t.Skip("avltree.ReplaceAll のテストは未実装")
}

func TestAlterAll(t *testing.T) {
	_ = AlterAll
	t.Skip("avltree.AlterAll のテストは未実装")
}

func TestFindAll(t *testing.T) {
	_ = FindAll
	t.Skip("avltree.FindAll のテストは未実装")
}

func TestMinAll(t *testing.T) {
	_ = MinAll
	t.Skip("avltree.MinAll のテストは未実装")
}

func TestMaxAll(t *testing.T) {
	_ = MaxAll
	t.Skip("avltree.MaxAll のテストは未実装")
}

func TestDeleteIterate(t *testing.T) {
	_ = DeleteIterate
	t.Skip("avltree.DeleteIterate のテストは未実装")
}

func TestDeleteRange(t *testing.T) {
	_ = DeleteRange
	t.Skip("avltree.DeleteRange のテストは未実装")
}

func TestDeleteRangeIterate(t *testing.T) {
	_ = DeleteRangeIterate
	t.Skip("avltree.DeleteRangeIterate のテストは未実装")
}

func TestUpdateIterate(t *testing.T) {
	_ = UpdateIterate
	t.Skip("avltree.UpdateIterate のテストは未実装")
}

func TestUpdateRange(t *testing.T) {
	_ = UpdateRange
	t.Skip("avltree.UpdateRange のテストは未実装")
}

func TestUpdateRangeIterate(t *testing.T) {
	_ = UpdateRangeIterate
	t.Skip("avltree.UpdateRangeIterate のテストは未実装")
}

func TestReplaceRange(t *testing.T) {
	_ = ReplaceRange
	t.Skip("avltree.ReplaceRange のテストは未実装")
}

func TestAlterIterate(t *testing.T) {
	_ = AlterIterate
	t.Skip("avltree.AlterIterate のテストは未実装")
}

func TestAlterRange(t *testing.T) {
	_ = AlterRange
	t.Skip("avltree.AlterRange のテストは未実装")
}

func TestAlterRangeIterate(t *testing.T) {
	_ = AlterRangeIterate
	t.Skip("avltree.AlterRangeIterate のテストは未実装")
}

func TestInner_countNode(t *testing.T) {
	_ = countNode
	t.Skip("countNode のテストは未実装")
}

func TestInner_countRange(t *testing.T) {
	_ = countRange
	t.Skip("countRange のテストは未実装")
}

func TestInner_countExtendedRange(t *testing.T) {
	_ = countExtendedRange
	t.Skip("countExtendedRange のテストは未実装")
}

func TestInner_getHeight(t *testing.T) {
	_ = getHeight
	t.Skip("getHeight のテストは未実装")
}

func TestInner_checkBalance(t *testing.T) {
	_ = checkBalance
	t.Skip("checkBalance のテストは未実装")
}

func TestInner_compareNodeHeight(t *testing.T) {
	_ = compareNodeHeight
	t.Skip("compareNodeHeight のテストは未実装")
}

func TestInner_compareChildHeight(t *testing.T) {
	_ = compareChildHeight
	t.Skip("compareChildHeight のテストは未実装")
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
	t.Skip("calcNewHeight のテストは未実装")
}

func TestInner_setChildren(t *testing.T) {
	_ = setChildren
	t.Skip("setChildren のテストは未実装")
}

func TestInner_setLeftChild(t *testing.T) {
	_ = setLeftChild
	t.Skip("setLeftChild のテストは未実装")
}

func TestInner_setRightChild(t *testing.T) {
	_ = setRightChild
	t.Skip("setRightChild のテストは未実装")
}

func TestInner_resetNode(t *testing.T) {
	_ = resetNode
	t.Skip("resetNode のテストは未実装")
}

func TestInner_insertHelper_newNode(t *testing.T) {
	var helper insertHelper
	_ = helper.newNode
	t.Skip("insertHelper.newNode のテストは未実装")
}

func TestInner_insertHelper_compareKey(t *testing.T) {
	var helper insertHelper
	_ = helper.compareKey
	t.Skip("insertHelper.compareKey のテストは未実装")
}

func TestInner_insertHelper_allowDuplicateKeys(t *testing.T) {
	var helper insertHelper
	_ = helper.allowDuplicateKeys
	t.Skip("insertHelper.allowDuplicateKeys のテストは未実装")
}

func TestInner_insertHelper_insertTo(t *testing.T) {
	var helper insertHelper
	_ = helper.insertTo
	t.Skip("insertHelper.insertTo のテストは未実装")
}

func TestInner_rotate(t *testing.T) {
	_ = rotate
	t.Skip("rotate のテストは未実装")
}

func TestInner_rotateRight(t *testing.T) {
	_ = rotateRight
	t.Skip("rotateRight のテストは未実装")
}

func TestInner_rotateLeft(t *testing.T) {
	_ = rotateLeft
	t.Skip("rotateLeft のテストは未実装")
}

func TestInner_removeNode(t *testing.T) {
	_ = removeNode
	t.Skip("removeNode のテストは未実装")
}

func TestInner_removeMin(t *testing.T) {
	_ = removeMin
	t.Skip("removeMin のテストは未実装")
}

func TestInner_removeMax(t *testing.T) {
	_ = removeMax
	t.Skip("removeMax のテストは未実装")
}

func TestInner_ascIterateNode(t *testing.T) {
	_ = ascIterateNode
	t.Skip("ascIterateNode のテストは未実装")
}

func TestInner_descIterateNode(t *testing.T) {
	_ = descIterateNode
	t.Skip("descIterateNode のテストは未実装")
}

func TestInner_newKeyBounds(t *testing.T) {
	_ = newKeyBounds
	t.Skip("newKeyBounds のテストは未実装")
}

func TestInner_upperBound(t *testing.T) {
	var _ upperBound
	t.Skip("upperBound のテストは未実装")
}

func TestInner_lowerBound(t *testing.T) {
	var _ lowerBound
	t.Skip("lowerBound のテストは未実装")
}

func TestInner_bothBounds(t *testing.T) {
	var _ bothBounds
	t.Skip("bothBounds のテストは未実装")
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
	t.Skip("upperBoundsChecker のテストは未実装")
}

func TestInner_lowerBoundsChecker(t *testing.T) {
	var _ lowerBoundsChecker
	t.Skip("lowerBoundsChecker のテストは未実装")
}

func TestInner_ascRangeNode(t *testing.T) {
	_ = ascRangeNode
	t.Skip("ascRangeNode のテストは未実装")
}

func TestInner_descRangeNode(t *testing.T) {
	_ = descRangeNode
	t.Skip("descRangeNode のテストは未実装")
}

func TestInner_updateValue(t *testing.T) {
	_ = updateValue
	t.Skip("updateValue のテストは未実装")
}

func TestInner_ascUpdateIterate(t *testing.T) {
	_ = ascUpdateIterate
	t.Skip("ascUpdateIterate のテストは未実装")
}

func TestInner_descUpdateIterate(t *testing.T) {
	_ = descUpdateIterate
	t.Skip("descUpdateIterate のテストは未実装")
}

func TestInner_ascUpdateRange(t *testing.T) {
	_ = ascUpdateRange
	t.Skip("ascUpdateRange のテストは未実装")
}

func TestInner_descUpdateRange(t *testing.T) {
	_ = descUpdateRange
	t.Skip("descUpdateRange のテストは未実装")
}

func TestInner_ascDeleteIterate(t *testing.T) {
	_ = ascDeleteIterate
	t.Skip("ascDeleteIterate のテストは未実装")
}

func TestInner_descDeleteIterate(t *testing.T) {
	_ = descDeleteIterate
	t.Skip("descDeleteIterate のテストは未実装")
}

func TestInner_ascDeleteRange(t *testing.T) {
	_ = ascDeleteRange
	t.Skip("ascDeleteRange のテストは未実装")
}

func TestInner_descDeleteRange(t *testing.T) {
	_ = descDeleteRange
	t.Skip("descDeleteRange のテストは未実装")
}

func TestInner_alterNode(t *testing.T) {
	var _ alterNode
	t.Skip("alterNode のテストは未実装")
}

func TestAlterRequest(t *testing.T) {
	var _ AlterRequest
	t.Skip("AlterRequest のテストは未実装")
}

func TestInner_alter(t *testing.T) {
	_ = alter
	t.Skip("alter のテストは未実装")
}

func TestInner_ascAlterIterat(t *testing.T) {
	_ = ascAlterIterate
	t.Skip("ascAlterIterate のテストは未実装")
}

func TestInner_descAlterIterat(t *testing.T) {
	_ = descAlterIterate
	t.Skip("descAlterIterate のテストは未実装")
}

func TestInner_ascAlterRange(t *testing.T) {
	_ = ascAlterRange
	t.Skip("ascAlterRange のテストは未実装")
}

func TestInner_descAlterRange(t *testing.T) {
	_ = descAlterRange
	t.Skip("descAlterRange のテストは未実装")
}
