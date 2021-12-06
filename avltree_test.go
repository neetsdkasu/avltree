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

func TestKeyOrderingLessThanOtherKey(t *testing.T) {
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
}

func TestKeyOrderingEqualToOtherKey(t *testing.T) {
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
}

func TestKeyOrderingGreaterThanOtherKey(t *testing.T) {
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
}
