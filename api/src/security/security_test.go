package security

import (
	"reflect"
	"testing"
)

var (
	passwords = []string{
		"1234567",
		"7654321",
	}

	fakeHash = []byte{7}
)

func TestMain(t *testing.T) {
	for _, password := range passwords {
		hash, erro := Hash(password)
		if erro != nil {
			t.Error(erro)
		}

		if reflect.TypeOf(hash) != reflect.TypeOf(fakeHash) {
			t.Error("the type of hashed password is wrong")
		}
	}
}
