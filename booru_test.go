package booru_test

import (
	"fmt"

	booru "github.com/c032/go-booru"
	boorusources "github.com/c032/go-booru/sources"
)

func ExampleYandere() {
	var (
		b booru.Booru = boorusources.Yandere

		err       error
		sfwImages []booru.ImageReference
	)

	tags := "rating:safe"

	sfwImages, err = b.Search(tags)
	if err != nil {
		panic(err)
	}

	for _, ref := range sfwImages {
		fmt.Printf("Image URL: %s\n", ref.URL())
	}
}
