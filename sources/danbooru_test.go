package sources_test

import (
	"testing"

	"github.com/c032/go-booru/sources"
)

func TestDanbooru_Image(t *testing.T) {
	const imageID = 2867318

	image, err := sources.Danbooru.Image(imageID)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatalf("Danbooru.Image(%d) = nil", imageID)
	}

	if image.FileURL == "" {
		t.Errorf("Danbooru.Image(%d).FileURL = %q; want non-empty string", imageID, image.FileURL)
	}
	if len(image.Tags) == 0 {
		t.Errorf("len(Danbooru.Image(%d).Tags) = 0; want at least one tag", imageID)
	} else {
		for i, tag := range image.Tags {
			if tag.Label == "" {
				t.Errorf("Danbooru.Image(%d).Tags[%d] = %q; want non-empty string", imageID, i, tag.Label)
			}
		}
	}
}

func TestDanbooru_ParseURL(t *testing.T) {
	const (
		imageURL        = "https://danbooru.donmai.us/posts/2867318"
		expectedImageID = 2867318
	)

	imageID, err := sources.Danbooru.ParseURL(imageURL)
	if err != nil {
		t.Fatal(err)
	}

	if imageID != expectedImageID {
		t.Errorf("Danbooru.ParseURL(%q) = %d; want %d", imageURL, imageID, expectedImageID)
	}
}

func TestDanbooru_Search(t *testing.T) {
	tags := "rating:safe"

	result, err := sources.Danbooru.Search(tags)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("len(Danbooru.Search(%q)) = 0; want at least one result", tags)
	}
	for i, r := range result {
		if got := r.URL(); got == "" {
			t.Errorf("Danbooru.Search(%q)[%d].URL() = %q; want non-empty string", tags, i, got)
		}
	}
}
