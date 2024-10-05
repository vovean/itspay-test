package testkit

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// AssertEqualCmp asserts that two objects are equal.
func AssertEqualCmp(t *testing.T, expected any, actual any, options ...cmp.Option) bool {
	t.Helper()

	if diff := cmp.Diff(expected, actual, options...); diff != "" {
		return assert.Fail(t, fmt.Sprintf(
			"Not equal: \n"+
				"expected: %s\n"+
				"actual  : %s\n"+
				"\n"+
				"Diff:\n%s",
			expected,
			actual,
			diff,
		))
	}

	return true
}

// RequireEqualCmp asserts that two objects are equal and fail now if not.
func RequireEqualCmp(t *testing.T, expected any, actual any, options ...cmp.Option) {
	t.Helper()

	if AssertEqualCmp(t, expected, actual, options...) {
		return
	}

	t.FailNow()
}
