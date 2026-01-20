//go:build cgo

package image

import (
	"image"
	"io"

	"github.com/jdeng/goheif"
)

func decodeHEIF(r io.Reader) (image.Image, error) {
	return goheif.Decode(r)
}
