package booru

import (
	"errors"
)

var (
	ErrImageNotFound = errors.New("image not found")
)

type ImageRating int

const (
	RatingSafe ImageRating = iota
	RatingQuestionable
	RatingExplicit

	RatingUnknown ImageRating = -1
)

type Image struct {
	FileURL string
	Rating  ImageRating
	Size    uint64
	Tags    []Tag
}

type ImageReference interface {
	URL() string
	Image() (*Image, error)
}
