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
	"errors"
	"strings"
	"time"
)

// Publications represents en publication made by user
type Publication struct {
	ID          uint64    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Content     string    `json:"content,omitempty"`
	AuthorID    uint64    `json:"authorid,omitempty"`
	AuthorNick  string    `json:"authornick,omitempty"`
	AuthorEmail string    `json:"authoremail,omitempty"`
	Likes       uint64    `json:"likes"`
	CreatedAt   time.Time `json:"createdat,omitempty"`
}

// Prepare will prepate the publication
func (publication *Publication) Prepare() error {
	if erro := publication.validate(); erro != nil {
		return erro
	}

	publication.format()

	return nil
}

func (publication *Publication) validate() error {
	if publication.Title == "" {
		return errors.New("the title cant be empty")
	}

	if publication.Content == "" {
		return errors.New("the content cant be empty")
	}

	return nil
}

func (publication *Publication) format() {
	publication.Title = strings.TrimSpace(publication.Title)
	publication.Content = strings.TrimSpace(publication.Content)
}
