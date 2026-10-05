package zca

// AttachmentSource is either a file path (Path set) or in-memory Data with a Filename ("name.ext").
type AttachmentSource struct {
	Path     string
	Data     []byte
	Filename string
	Metadata AttachmentMetadata
}

type AttachmentMetadata struct {
	TotalSize int64
	Width     int64
	Height    int64
}
