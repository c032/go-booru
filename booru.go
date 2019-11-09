package booru

import (
	"errors"
)

var (
	ErrIncompatibleURL = errors.New("incompatible url")
)

type Booru interface {
	Image(id int64) (*Image, error)
	ParseURL(urlStr string) (imageID int64, err error)
	Search(tags string) ([]ImageReference, error)
}
