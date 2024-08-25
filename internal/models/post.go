package models

import (
	"fmt"
	"time"

	"github.com/fdemchenko/arcus/internal/validator"
)

const PostTitleMaxLength = 70
const PostTagMaxLength = 40
const PostTagsMaxAmount = 6

type Post struct {
	ID        int
	Title     string
	Content   *string
	Tags      []string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Post) Validate(v validator.Validator) {
	v.Check(p.Title != "", "title", "must not be empty")
	v.Check(len(p.Title) <= PostTitleMaxLength, "title", fmt.Sprintf("must not be greater than %d characters long", PostTitleMaxLength))

	if p.Content != nil {
		v.Check(*p.Content != "", "description", "must not be empty")
	}

	v.Check(len(p.Tags) <= PostTagsMaxAmount, "tags", fmt.Sprintf("max tags amount is %d", PostTagsMaxAmount))

	for _, tag := range p.Tags {
		v.Check(tag != "", "tag", "must not be empty")
		v.Check(len(tag) <= PostTagMaxLength, "tag", fmt.Sprintf("must not be greater than %d characters long", PostTagMaxLength))
	}
}
