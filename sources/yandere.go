package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/c032/go-booru"
)

const yandereURLStr = "https://yande.re/post"

var yandereURL *url.URL

var Yandere booru.Booru = yandere{}

func init() {
	var err error

	yandereURL, err = url.Parse(yandereURLStr)
	if err != nil {
		panic(err)
	}
}

type yandere struct{}

func (d yandere) Image(id int64) (*booru.Image, error) {
	var (
		err error

		result []*yandereImage
	)

	tags := fmt.Sprintf("id:%d", id)

	result, err = d.rawSearch(tags)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, booru.ErrImageNotFound
	}

	rawImg := result[0]

	return rawImg.Image()
}

func (d yandere) ParseURL(urlStr string) (int64, error) {
	var (
		err error

		imageID  int64
		imageURL *url.URL
	)

	const prefix = "/post/show/"

	imageURL, err = url.Parse(urlStr)
	if err != nil {
		return 0, err
	}

	if !strings.HasPrefix(imageURL.Path, prefix) {
		err = booru.ErrIncompatibleURL

		return 0, err
	}

	imageIDStr := strings.TrimPrefix(imageURL.Path, prefix)
	if strings.Contains(imageIDStr, "/") {
		imageIDStr = imageIDStr[:strings.Index(imageIDStr, "/")]
	}

	imageID, err = strconv.ParseInt(imageIDStr, 10, 64)

	return imageID, err
}

func (d yandere) rawSearch(tags string) ([]*yandereImage, error) {
	var (
		err error

		searchURL *url.URL
		resp      *http.Response
	)

	searchURL, err = url.Parse(yandereURLStr)
	if err != nil {
		return nil, err
	}

	searchURL.Path = "/post.json"

	q := url.Values{}
	q.Set("tags", tags)

	searchURL.RawQuery = q.Encode()

	resp, err = http.Get(searchURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	images := make([]*yandereImage, 0, 40)
	jd := json.NewDecoder(resp.Body)

	err = jd.Decode(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (d yandere) Search(tags string) ([]booru.ImageReference, error) {
	var (
		err    error
		images []*yandereImage
	)

	images, err = d.rawSearch(tags)
	if err != nil {
		return nil, err
	}

	result := make([]booru.ImageReference, len(images))
	for i, img := range images {
		result[i] = img
	}

	return result, nil
}

type yandereImage struct {
	ID        int64  `json:"id"`
	ParentID  int64  `json:"parent_id"`
	Source    string `json:"source"`
	MD5       string `json:"md5"`
	RatingStr string `json:"rating"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      uint64 `json:"file_size"`
	TagsStr   string `json:"tags"`
	FileURL   string `json:"file_url"`
	SampleURL string `json:"sample_url"`
	JPEGURL   string `json:"jpeg_url"`
}

func (yi *yandereImage) Image() (*booru.Image, error) {
	var (
		err     error
		fileURL *url.URL
	)

	fileURL, err = yandereURL.Parse(yi.FileURL)
	if err != nil {
		return nil, err
	}

	img := &booru.Image{
		FileURL: fileURL.String(),
		Rating:  yi.Rating(),
		Size:    yi.Size,
		Tags:    yi.Tags(),
	}

	return img, nil
}

func (yi *yandereImage) Rating() booru.ImageRating {
	switch yi.RatingStr {
	case "s":
		return booru.RatingSafe
	case "q":
		return booru.RatingQuestionable
	case "e":
		return booru.RatingExplicit
	}

	return booru.RatingUnknown
}

func (yi *yandereImage) Tags() []booru.Tag {
	var tags []booru.Tag

	for _, label := range strings.Split(yi.TagsStr, " ") {
		if label == "" {
			continue
		}

		tags = append(tags, booru.Tag{
			Label: label,
			Type:  booru.TagGeneral,
		})
	}

	return tags
}

func (yi *yandereImage) URL() string {
	return fmt.Sprintf("https://yande.re/post/show/%d", yi.ID)
}
