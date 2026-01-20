//go:build !cgo

package image

import (
	"fmt"
	"image"
	"io"
)

func decodeHEIF(_ io.Reader) (image.Image, error) {
	return nil, fmt.Errorf("HEIC/HEIF decoding is disabled in this build (requires CGO)")
}
