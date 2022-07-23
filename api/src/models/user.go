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

package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// User represents an Social media User
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Pass      string    `json:"pass,omitempty"`
	CreatedAt time.Time `json:"createdat,omitempty"`
}

// Prepare will validate the User Struct.
func (user *User) Prepare(stage string) error {
	if erro := user.test(stage); erro != nil {
		return erro
	}

	if erro := user.format(stage); erro != nil {
		return erro
	}

	return nil
}

func (user *User) test(stage string) error {
	if user.Name == "" {
		return errors.New("the identity name cant be empty")
	}

	if user.Nick == "" {
		return errors.New("the identity nick cant be empty")
	}

	if user.Email == "" {
		return errors.New("the identity email cant be empty")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("the identity email is invalid")
	}

	if stage == "registration" && user.Pass == "" {
		return errors.New("the identity pass cant be empty")
	}

	return nil
}

func (user *User) format(stage string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if stage == "registration" {
		passHashed, erro := security.Hash(user.Pass)
		if erro != nil {
			return erro
		}

		user.Pass = string(passHashed)
	}

	return nil
}
