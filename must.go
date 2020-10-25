package mu

import (
	"io/ioutil"
	"os"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(filename string) []byte {
	result, err := ioutil.ReadFile(filename)
	PanicErr(err)

	return result
}

func WriteFile(filename string, data []byte, perm os.FileMode) {
	err := ioutil.WriteFile(filename, data, perm)
	PanicErr(err)
}
