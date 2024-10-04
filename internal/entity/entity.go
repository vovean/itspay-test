package entity

import (
	"time"

	"github.com/cockroachdb/apd/v3"
)

type Rate struct {
	Ask        *apd.Decimal
	Bid        *apd.Decimal
	ReceivedAt time.Time
}
