package archive

import "io"

type Bundle struct {
	content  map[string][]byte
	manifest Manifest
}

// Manifest is created and loaded into a bundle which contains information on the files
type Manifest struct {
	Created time.Time                 `json:"createdAt"`
	Version string                    `json:"version"`
	Files   map[string]fileDescriptor `json:"files"`
}

type fileDescriptor struct {
	Added      time.Time         `json:"addedAt"`
	Properties map[string]string `json:"properties"`
	Digest     string            `json:"digest"`
}

func NewBundle() *Bundle {
	return &Bundle{}
}

func EncodeBundle(dst io.Writer, bundle *Bundle) error {
	return nil
}

func DecodeBundle(src io.Reader, bundle *Bundle) error {
	return nil
}
