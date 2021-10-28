package mu

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			in := make([]byte, i)
			rand.Read(in)
			enc := Base64Encode(in)
			dec, err := Base64Decode(enc)

			assert.NoError(t, err)
			assert.Equal(t, in, dec)
		}
	})

	t.Run("bad decode", func(t *testing.T) {
		_, err := Base64Decode([]byte("a"))
		assert.Error(t, err)
	})
}
