package mu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	//fmt.Println(DeriveKey("a", nil))
}

func TestRandBytes(t *testing.T) {
	assert.Len(t, RandBytes(1000), 1000)
	assert.Len(t, RandBytes(0), 0)
	assert.NotEqual(t, RandBytes(10), RandBytes(10))
}

func TestDeriveKey(t *testing.T) {
	keyOrig, salt := DeriveKey("test password", nil)
	assert.Len(t, keyOrig, 32)
	assert.Len(t, salt, 16)

	key, _ := DeriveKey("test password", salt)
	assert.Equal(t, keyOrig, key)

	key, _ = DeriveKey("test password ", salt)
	assert.NotEqual(t, keyOrig, key)

	salt[0]++
	key, _ = DeriveKey("test password", salt)
	assert.NotEqual(t, keyOrig, key)
	salt[0]--

	// Test that parameter changes have an effect

	params := DefaultArgon2Params
	key, _ = DeriveKey("test password", salt, params)
	assert.Equal(t, keyOrig, key)

	params.Time++
	key, _ = DeriveKey("test password", salt, params)
	assert.NotEqual(t, keyOrig, key)

	params = DefaultArgon2Params
	params.Memory++
	key, _ = DeriveKey("test password", salt, params)
	assert.NotEqual(t, keyOrig, key)

	params = DefaultArgon2Params
	params.Threads++
	key, _ = DeriveKey("test password", salt, params)
	assert.NotEqual(t, keyOrig, key)

	params = DefaultArgon2Params
	params.KeyLen++
	key, _ = DeriveKey("test password", salt, params)
	assert.NotEqual(t, keyOrig, key)
}

func TestEncryptDecrypt(t *testing.T) {
	key := RandBytes(32)
	plaintext := []byte("Attack at dawn!!!")

	ciphertext := Encrypt(plaintext, key)
	assert.NotEqual(t, plaintext, ciphertext)

	newplaintext, err := Decrypt(ciphertext, key)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, newplaintext)

	ciphertext[0]++
	_, err = Decrypt(ciphertext, key)
	assert.Error(t, err)
	ciphertext[0]--

	key[0]++
	_, err = Decrypt(ciphertext, key)
	assert.Error(t, err)
}

func TestRandString(t *testing.T) {
	s1 := RandString(20)
	s2 := RandString(20)
	assert.NotEqual(t, s1, s2)
	assert.Len(t, s1, 20)
}
