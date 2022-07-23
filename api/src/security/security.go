/*
Copyright 2022 Danilo S. Lopes.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hash receive password and hash it
func Hash(pass string) ([]byte, error) {
	// password lenth cant be higher then 72 bytes
	if len([]byte(pass)) > 72 {
		return nil, errors.New("the password cant be higher then 72 bytes")
	}
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

// ValidatePass validates password and the hash are equal
func ValidatePass(passHashed, passString string) error {
	return bcrypt.CompareHashAndPassword([]byte(passHashed), []byte(passString))
}
