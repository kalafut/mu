package mu

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/dchest/uniuri"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

// Results in about 40ms generation time on a 2020 MBA
// Only 1MB, much less than the 64MB standard
var DefaultArgon2Params = Argon2Params{
	Time:    20,
	Memory:  4 * 1024,
	Threads: 1,
	KeyLen:  32,
}

func RandBytes(length int) []byte {
	d := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, d); err != nil {
		panic(err)
	}
	return d
}

func RandString(length int) string {
	return uniuri.NewLen(length)
}

type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func DeriveKey(password string, salt []byte, customParams ...Argon2Params) ([]byte, []byte) {
	if salt == nil {
		salt = RandBytes(16)
	}

	params := DefaultArgon2Params
	if len(customParams) > 0 {
		params = customParams[0]
	}

	//t := time.Now()
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)
	//fmt.Println(time.Now().Sub(t))

	return key, salt
}

func Encrypt(data, key []byte) []byte {
	if len(key) != 32 {
		panic("invalid key length")
	}

	var secretKey [32]byte
	var nonce [24]byte

	copy(secretKey[:], key)
	copy(nonce[:], RandBytes(24))

	return secretbox.Seal(nonce[:], []byte(data), &nonce, &secretKey)
}

func Decrypt(data, key []byte) ([]byte, error) {
	if len(key) != 32 {
		panic("invalid key length")
	}

	var secretKey [32]byte
	var nonce [24]byte

	copy(secretKey[:], key)

	copy(nonce[:], data[:24])
	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &secretKey)
	if !ok {
		return nil, errors.New("decryption failure")
	}

	return decrypted, nil
}
