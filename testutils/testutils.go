package testutils

import (
	"runtime/debug"
	"testing"
)

// FailNowOnErr is just like FailNow but only on error
func FailNowOnErr(t *testing.T, err error) {
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Log(err)
		t.FailNow()
	}
}
