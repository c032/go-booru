package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/c032/booru-go"
)

const danbooruURLStr = "https://danbooru.donmai.us/"

var danbooruURL *url.URL

var Danbooru booru.Booru = danbooru{}

func init() {
	var err error

	danbooruURL, err = url.Parse(danbooruURLStr)
	if err != nil {
		panic(err)
	}
}

type danbooru struct{}

func (d danbooru) Image(id int64) (*booru.Image, error) {
	var (
		err error

		result []*danbooruImage
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

func (d danbooru) ParseURL(urlStr string) (int64, error) {
	var (
		err error

		imageID  int64
		imageURL *url.URL
	)

	imageURL, err = url.Parse(urlStr)
	if err != nil {
		return 0, err
	}

	if !strings.HasPrefix(imageURL.Path, "/posts/") {
		err = booru.ErrIncompatibleURL

		return 0, err
	}

	imageIDStr := strings.TrimPrefix(imageURL.Path, "/posts/")
	if strings.Contains(imageIDStr, "/") {
		imageIDStr = imageIDStr[:strings.Index(imageIDStr, "/")]
	}

	imageID, err = strconv.ParseInt(imageIDStr, 10, 64)

	return imageID, err
}

func (d danbooru) rawSearch(tags string) ([]*danbooruImage, error) {
	var (
		err error

		searchURL *url.URL
		resp      *http.Response
	)

	searchURL, err = url.Parse(danbooruURLStr)
	if err != nil {
		return nil, err
	}

	searchURL.Path = "/posts.json"

	q := url.Values{}
	q.Set("tags", tags)

	searchURL.RawQuery = q.Encode()

	resp, err = http.Get(searchURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	images := make([]*danbooruImage, 0, 8)
	jd := json.NewDecoder(resp.Body)

	err = jd.Decode(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (d danbooru) Search(tags string) ([]booru.ImageReference, error) {
	var (
		err    error
		images []*danbooruImage
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

type danbooruImage struct {
	ID                 int64  `json:"id"`
	Source             string `json:"source"`
	MD5                string `json:"md5"`
	RatingStr          string `json:"rating"`
	Width              int    `json:"image_width"`
	Height             int    `json:"image_height"`
	Size               uint64 `json:"file_size"`
	PixivID            int64  `json:"pixiv_id"`
	TagString          string `json:"tag_string"`
	TagStringArtist    string `json:"tag_string_artist"`
	TagStringCharacter string `json:"tag_string_character"`
	TagStringCopyright string `json:"tag_string_copyright"`
	TagStringGeneral   string `json:"tag_string_general"`
	FileURL            string `json:"file_url"`
	LargeFileURL       string `json:"large_file_url"`
	PreviewFileURL     string `json:"preview_file_url"`
}

func (di *danbooruImage) Image() (*booru.Image, error) {
	var (
		err     error
		fileURL *url.URL
	)

	fileURL, err = danbooruURL.Parse(di.FileURL)
	if err != nil {
		return nil, err
	}

	img := &booru.Image{
		FileURL: fileURL.String(),
		Rating:  di.Rating(),
		Size:    di.Size,
		Tags:    di.Tags(),
	}

	return img, nil
}

func (di *danbooruImage) Rating() booru.ImageRating {
	switch di.RatingStr {
	case "s":
		return booru.RatingSafe
	case "q":
		return booru.RatingQuestionable
	case "e":
		return booru.RatingExplicit
	}

	return booru.RatingUnknown
}

func (di *danbooruImage) Tags() []booru.Tag {
	var tags []booru.Tag

	for _, label := range strings.Split(di.TagStringArtist, " ") {
		if label == "" {
			continue
		}

		tags = append(tags, booru.Tag{
			Label: label,
			Type:  booru.TagArtist,
		})
	}

	for _, label := range strings.Split(di.TagStringCharacter, " ") {
		if label == "" {
			continue
		}

		tags = append(tags, booru.Tag{
			Label: label,
			Type:  booru.TagCharacter,
		})
	}

	for _, label := range strings.Split(di.TagStringCopyright, " ") {
		if label == "" {
			continue
		}

		tags = append(tags, booru.Tag{
			Label: label,
			Type:  booru.TagCopyright,
		})
	}

	for _, label := range strings.Split(di.TagStringGeneral, " ") {
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

func (di *danbooruImage) URL() string {
	return fmt.Sprintf("https://danbooru.donmai.us/posts/%d", di.ID)
}
