package image

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/jackmordaunt/icns"
	"github.com/kpfaulkner/jxl-go"
	"github.com/sergeymakinen/go-bmp"
	"github.com/sergeymakinen/go-ico"
	"golang.org/x/image/tiff"
)

// ResizeMode defines how the image should be resized
type ResizeMode string

const (
	// ResizeModeFit resizes the image to fit within the specified dimensions while maintaining aspect ratio
	ResizeModeFit ResizeMode = "fit"
	// ResizeModeFill resizes the image to fill the specified dimensions while maintaining aspect ratio and crops any excess
	ResizeModeFill ResizeMode = "fill"
	// ResizeModeStretch resizes the image to the specified dimensions without maintaining aspect ratio
	ResizeModeStretch ResizeMode = "stretch"
)

// ProcessOptions contains all options for image processing
type ProcessOptions struct {
	Width      int        // Target width
	Height     int        // Target height
	ResizeMode ResizeMode // How to resize the image
	Quality    int        // Output quality (1-100, only for JPEG)
	OutputFormat string   // Output format (jpg, png, gif)
	PadColor   [3]uint8   // RGB color to use for padding
}

// DefaultOptions returns the default processing options
func DefaultOptions() ProcessOptions {
	return ProcessOptions{
		Width:        800,
		Height:       600,
		ResizeMode:   ResizeModeFit,
		Quality:      85,
		OutputFormat: "",
		PadColor:     [3]uint8{255, 255, 255}, // White
	}
}

// OpenImage opens an image file and decodes it based on its format
func OpenImage(filename string) (image.Image, error) {
	// Get file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the image based on its format
	var img image.Image
	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "tiff", "tif":
		// Use imaging library for standard formats
		return imaging.Open(filename)
	case "webp":
		img, err = webp.Decode(file)
	case "avif":
		img, err = avif.Decode(file)
	case "ico":
		img, err = ico.Decode(file)
	case "icns":
		img, err = icns.Decode(file)
	case "heic", "heif":
		img, err = decodeHEIF(file)
	case "jxl":
		// Reset file pointer to beginning
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to reset file pointer: %w", err)
		}
		img, err = jxl_go.Decode(file)
	case "jp2":
		// JP2 is not directly supported by any Go library
		// We could potentially use an external tool or library for this
		return nil, fmt.Errorf("JPEG 2000 (.jp2) format is not supported for decoding")
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// ProcessImage processes an image according to the provided options
func ProcessImage(inputPath, outputPath string, options ProcessOptions) error {
	// Open the input file using our custom function that supports more formats
	src, err := OpenImage(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Determine output format if not specified
	if options.OutputFormat == "" {
		options.OutputFormat = strings.TrimPrefix(filepath.Ext(outputPath), ".")
		if options.OutputFormat == "" {
			// Default to JPEG if no extension is provided
			options.OutputFormat = "jpg"
		}
	}

	// Resize the image according to the specified mode
	var resized *image.NRGBA
	switch options.ResizeMode {
	case ResizeModeFit:
		resized = imaging.Fit(src, options.Width, options.Height, imaging.Lanczos)
		// If padding is needed, create a new image with the target dimensions and paste the resized image in the center
		if resized.Bounds().Dx() < options.Width || resized.Bounds().Dy() < options.Height {
			bgColor := color.RGBA{
				R: options.PadColor[0],
				G: options.PadColor[1],
				B: options.PadColor[2],
				A: 255,
			}
			bg := imaging.New(options.Width, options.Height, bgColor)
			resized = imaging.PasteCenter(bg, resized)
		}
	case ResizeModeFill:
		resized = imaging.Fill(src, options.Width, options.Height, imaging.Center, imaging.Lanczos)
	case ResizeModeStretch:
		resized = imaging.Resize(src, options.Width, options.Height, imaging.Lanczos)
	default:
		return fmt.Errorf("unknown resize mode: %s", options.ResizeMode)
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Save the image in the specified format
	switch strings.ToLower(options.OutputFormat) {
	case "jpg", "jpeg":
		err = jpeg.Encode(out, resized, &jpeg.Options{Quality: options.Quality})
	case "png":
		err = png.Encode(out, resized)
	case "gif":
		err = gif.Encode(out, resized, nil)
	case "bmp":
		err = bmp.Encode(out, resized)
	case "tiff", "tif":
		err = tiff.Encode(out, resized, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case "webp":
		err = webp.Encode(out, resized, &webp.Options{Lossless: false, Quality: float32(options.Quality)})
	case "avif":
		err = avif.Encode(out, resized, avif.Options{Quality: options.Quality, Speed: 8})
	case "ico":
		err = ico.Encode(out, resized)
	case "icns":
		// Use the resized image directly for ICNS encoding
		err = icns.Encode(out, resized)
	case "heic", "heif":
		// The goheif library (github.com/jdeng/goheif) only supports decoding HEIC/HEIF images, not encoding
		// There is no Go library available that supports encoding to HEIC/HEIF format
		return fmt.Errorf("encoding to %s format is not supported: the goheif library only provides decoding capability", options.OutputFormat)
	case "jxl":
		// The jxl-go library (github.com/kpfaulkner/jxl-go) only supports decoding JXL images, not encoding
		return fmt.Errorf("encoding to JXL format is not supported: the jxl-go library only provides decoding capability")
	case "jp2":
		// There's no Go library for JP2 encoding
		return fmt.Errorf("encoding to JPEG 2000 format is not supported: no Go library available for JP2 encoding")
	default:
		return fmt.Errorf("unsupported output format: %s", options.OutputFormat)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}
