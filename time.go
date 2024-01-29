package mu

import (
	"time"
)

type Time struct {
	time.Time

	override time.Time
}

func (t *Time) Set(override time.Time) {
	t.override = override
}

func (t *Time) Clear() {
	t.override = time.Time{}
}

func (t *Time) Now() time.Time {
	if !t.override.IsZero() {
		return t.override
	}

	return time.Now()
}

func (t *Time) Advance(d time.Duration) {
	orig := time.Now()
	if !t.override.IsZero() {
		orig = t.override
	}
	t.override = orig.Add(d)
}

func (t *Time) Rewind(d time.Duration) {
	t.Advance(-d)
}

func (t *Time) NowTrunc() time.Time {
	return t.Now().UTC().Truncate(time.Second)
}

func NowTrunc() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
