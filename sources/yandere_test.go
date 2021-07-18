package sources_test

import (
	"testing"

	"github.com/c032/go-booru/sources"
)

func TestYandere_Image(t *testing.T) {
	const imageID = 818148

	image, err := sources.Yandere.Image(imageID)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatalf("Yandere.Image(%d) = nil", imageID)
	}

	if image.FileURL == "" {
		t.Errorf("Yandere.Image(%d).FileURL = %q; want non-empty string", imageID, image.FileURL)
	}
	if len(image.Tags) == 0 {
		t.Errorf("len(Yandere.Image(%d).Tags) = 0; want at least one tag", imageID)
	} else {
		for i, tag := range image.Tags {
			if tag.Label == "" {
				t.Errorf("Yandere.Image(%d).Tags[%d] = %q; want non-empty string", imageID, i, tag.Label)
			}
		}
	}
}

func TestYandere_ParseURL(t *testing.T) {
	const (
		imageURL        = "https://yande.re/post/show/818148"
		expectedImageID = 818148
	)

	imageID, err := sources.Yandere.ParseURL(imageURL)
	if err != nil {
		t.Fatal(err)
	}

	if imageID != expectedImageID {
		t.Errorf("Yandere.ParseURL(%q) = %d; want %d", imageURL, imageID, expectedImageID)
	}
}

func TestYandere_Search(t *testing.T) {
	tags := "rating:safe"

	result, err := sources.Yandere.Search(tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("len(Yandere.Search(%q)) = 0; want at least one result", tags)
	}
	for i, r := range result {
		if got := r.URL(); got == "" {
			t.Errorf("Yandere.Search(%q)[%d].URL() = %q; want non-empty string", tags, i, got)
		}
	}
}
