package models

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/fdemchenko/arcus/internal/validator"
)

const PostTitleMaxLength = 70
const PostTagMaxLength = 40
const PostTagsMaxAmount = 6

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   *string   `json:"content"`
	Tags      []string  `json:"tags"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) Validate(v validator.Validator) {
	v.Check(p.Title != "", "title", "must not be empty")
	v.Check(utf8.RuneCountInString(p.Title) <= PostTitleMaxLength, "title", fmt.Sprintf("must not be greater than %d characters long", PostTitleMaxLength))

	if p.Content != nil {
		v.Check(*p.Content != "", "description", "must not be empty")
	}

	v.Check(len(p.Tags) <= PostTagsMaxAmount, "tags", fmt.Sprintf("max tags amount is %d", PostTagsMaxAmount))

	for _, tag := range p.Tags {
		v.Check(tag != "", "tag", "must not be empty")
		v.Check(utf8.RuneCountInString(tag) <= PostTagMaxLength, "tag", fmt.Sprintf("must not be greater than %d characters long", PostTagMaxLength))
	}
}
