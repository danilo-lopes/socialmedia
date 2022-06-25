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
