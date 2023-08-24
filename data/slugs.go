package data

import "fmt"

// SlugNotFound is an error raised when a slug can not be found in the database
var SlugNotFound = fmt.Errorf("slug not found")

// Slug defines the structure for an API slug
type Slug struct {
	// the id for the slug
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the slug

	// the name for the slug
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`
}

var slugs []Slug = []Slug{
	{
		ID:   1,
		Name: "AVITO_VOICE_MESSAGES",
	},
	{
		ID:   2,
		Name: "AVITO_PERFORMANCE_VAS",
	},
	{
		ID:   3,
		Name: "AVITO_DISCOUNT_30",
	},
}

func AddSlug(slug Slug) {
	if len(slugs) == 0 {
		slug.ID = 1
	} else {
		slug.ID = slugs[len(slugs)-1].ID + 1
	}

	slugs = append(slugs, slug)
}

func GetSlugs() []Slug {
	return slugs
}

func GetSlugByID(id int) (*Slug, error) {
	for _, slug := range slugs {
		if slug.ID == id {
			return &slug, nil
		}
	}

	return nil, SlugNotFound
}
