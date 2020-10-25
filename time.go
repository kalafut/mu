package mu

import (
	"time"
)

func NowTrunc() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
