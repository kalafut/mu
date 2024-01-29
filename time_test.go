package mu

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	mt := Time{}
	assert.InDelta(time.Now().UnixMilli(), mt.Now().UnixMilli(), 20)

	v := time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC)
	mt.Set(v)
	assert.Equal(v.UnixMilli(), mt.Now().UnixMilli())
	mt.Clear()

	assert.InDelta(time.Now().UnixMilli(), mt.Now().UnixMilli(), 20)
}

func TestAdvance(t *testing.T) {
	assert := assert.New(t)

	mt := Time{}
	mt.Advance(time.Hour)
	assert.InDelta(time.Until(mt.Now()).Milliseconds(), time.Hour.Milliseconds(), 2)

	mt = Time{}
	mt.Rewind(time.Hour)
	assert.InDelta(time.Since(mt.Now()).Milliseconds(), time.Hour.Milliseconds(), 2)

	// mt.Advance(time.Hour)
	// assert.InDelta(time.Now().Add(time.Hour).UnixMilli(), mt.Now().UnixMilli(), 20)

	// mt.Advance(-time.Hour)
	// assert.InDelta(time.Now().UnixMilli(), mt.Now().UnixMilli(), 20)
}
