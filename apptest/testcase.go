package apptest

import (
	"testing"
)

type TestCase struct {
	t        *testing.T
	cleanups []func()
}

func NewTestCase(t *testing.T) *TestCase {
	t.Helper()
	return &TestCase{t: t}
}

func (tc *TestCase) T() *testing.T { return tc.t }

func (tc *TestCase) RegisterCleanup(fn func()) {
	tc.cleanups = append(tc.cleanups, fn)
}

func (tc *TestCase) Stop() {
	for i := len(tc.cleanups) - 1; i >= 0; i-- {
		tc.cleanups[i]()
	}
}
