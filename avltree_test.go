package avltree

import (
	"runtime"
	"testing"
	"testing/quick"
)

type assertion testing.T

func (a *assertion) Assert(expectTrue bool, args ...interface{}) {
	if !expectTrue {
		if _, file, line, ok := runtime.Caller(1); ok {
			args = append([]interface{}{
				"Fail Assert:", "[", file, line, "]",
			}, args...)
		}
		a.Fatal(args...)
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
	a.Assert(
		LessThanOtherKey.LessThan(),
		"LessThanOtherKey.LessThan()",
	)
	a.Assert(
		LessThanOtherKey.LessThanOrEqualTo(),
		"LessThanOtherKey.LessThanOrEqualTo()",
	)
	a.Assert(
		!LessThanOtherKey.EqualTo(),
		"!LessThanOtherKey.EqualTo()",
	)
	a.Assert(
		LessThanOtherKey.NotEqualTo(),
		"LessThanOtherKey.NotEqualTo()",
	)
	a.Assert(
		!LessThanOtherKey.GreaterThan(),
		"!LessThanOtherKey.GreaterThan()",
	)
	a.Assert(
		!LessThanOtherKey.GreaterThanOrEqualTo(),
		"!LessThanOtherKey.GreaterThanOrEqualTo()",
	)
}

func TestKeyOrderingEqualToOtherKey(t *testing.T) {
	a := (*assertion)(t)
	a.Assert(
		!EqualToOtherKey.LessThan(),
		"!EqualToOtherKey.LessThan()",
	)
	a.Assert(
		EqualToOtherKey.LessThanOrEqualTo(),
		"EqualToOtherKey.LessThanOrEqualTo()",
	)
	a.Assert(
		EqualToOtherKey.EqualTo(),
		"EqualToOtherKey.EqualTo()",
	)
	a.Assert(
		!EqualToOtherKey.NotEqualTo(),
		"!EqualToOtherKey.NotEqualTo()",
	)
	a.Assert(
		!EqualToOtherKey.GreaterThan(),
		"!EqualToOtherKey.GreaterThan()",
	)
	a.Assert(
		EqualToOtherKey.GreaterThanOrEqualTo(),
		"EqualToOtherKey.GreaterThanOrEqualTo()",
	)
}

func TestKeyOrderingGreaterThanOtherKey(t *testing.T) {
	a := (*assertion)(t)
	a.Assert(
		!GreaterThanOtherKey.LessThan(),
		"!GreaterThanOtherKey.LessThan()",
	)
	a.Assert(
		!GreaterThanOtherKey.LessThanOrEqualTo(),
		"!GreaterThanOtherKey.LessThanOrEqualTo()",
	)
	a.Assert(
		!GreaterThanOtherKey.EqualTo(),
		"!GreaterThanOtherKey.EqualTo()",
	)
	a.Assert(
		GreaterThanOtherKey.NotEqualTo(),
		"GreaterThanOtherKey.NotEqualTo()",
	)
	a.Assert(
		GreaterThanOtherKey.GreaterThan(),
		"GreaterThanOtherKey.GreaterThan()",
	)
	a.Assert(
		GreaterThanOtherKey.GreaterThanOrEqualTo(),
		"GreaterThanOtherKey.GreaterThanOrEqualTo()",
	)
}
