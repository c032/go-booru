package booru

type TagType int

const (
	TagGeneral TagType = iota
	TagArtist
	TagCharacter
	TagCopyright
)

type Tag struct {
	Label string
	Type  TagType
}
