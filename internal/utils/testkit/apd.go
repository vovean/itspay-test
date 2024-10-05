package testkit

import (
	"fmt"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/go-cmp/cmp"
)

var DecimalComparer = cmp.Comparer(func(l, r apd.Decimal) bool { return l.Cmp(&r) == 0 })

func MustDecimalFromString(value string) *apd.Decimal {
	decimal, _, err := apd.NewFromString(value)
	if err != nil {
		panic(fmt.Errorf("can't convert string %q to apd.Decimal: %w", value, err))
	}

	return decimal
}
