package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestImage creates a test image with the given dimensions and color.
func createTestImage(width, height int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
	return img
}

// saveTestImage saves a test image to a temporary file.
func saveTestImage(img *image.RGBA, format string) (string, error) {
	tmpFile, err := os.CreateTemp("", "test-image-*."+format)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(tmpFile, img)
	default:
		return "", fmt.Errorf("unsupported test format: %s", format)
	}

	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestProcessImage(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})

	inputPath, err := saveTestImage(img, "jpg")
	if err != nil {
		t.Fatalf("failed to save test image: %v", err)
	}
	defer os.Remove(inputPath)

	testCases := []struct {
		name         string
		options      ProcessOptions
		expectWidth  int
		expectHeight int
	}{
		{
			name: "resize with fit mode",
			options: ProcessOptions{
				Width:        200,
				Height:       200,
				ResizeMode:   ResizeModeFit,
				Quality:      90,
				OutputFormat: "png",
				PadColor:     [3]uint8{255, 255, 255},
			},
			expectWidth:  200,
			expectHeight: 200,
		},
		{
			name: "resize with fill mode",
			options: ProcessOptions{
				Width:        50,
				Height:       50,
				ResizeMode:   ResizeModeFill,
				Quality:      90,
				OutputFormat: "jpg",
				PadColor:     [3]uint8{0, 0, 0},
			},
			expectWidth:  50,
			expectHeight: 50,
		},
		{
			name: "resize with stretch mode",
			options: ProcessOptions{
				Width:        150,
				Height:       75,
				ResizeMode:   ResizeModeStretch,
				Quality:      90,
				OutputFormat: "png",
				PadColor:     [3]uint8{0, 255, 0},
			},
			expectWidth:  150,
			expectHeight: 75,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outputPath := filepath.Join(t.TempDir(), "output."+tc.options.OutputFormat)
			if err := ProcessImage(inputPath, outputPath, tc.options); err != nil {
				t.Fatalf("ProcessImage failed: %v", err)
			}

			outFile, err := os.Open(outputPath)
			if err != nil {
				t.Fatalf("failed to open output file: %v", err)
			}
			defer outFile.Close()

			cfg, _, err := image.DecodeConfig(outFile)
			if err != nil {
				t.Fatalf("failed to decode output image config: %v", err)
			}
			if cfg.Width != tc.expectWidth || cfg.Height != tc.expectHeight {
				t.Fatalf("unexpected output dimensions: got %dx%d, expected %dx%d", cfg.Width, cfg.Height, tc.expectWidth, tc.expectHeight)
			}
		})
	}
}

func TestOpenImageUnsupportedFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("not an image"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	_, err := OpenImage(path)
	if err == nil {
		t.Fatal("expected error for unsupported image format")
	}
	if !strings.Contains(err.Error(), "unsupported image format") {
		t.Fatalf("expected unsupported format error, got: %v", err)
	}
}

func TestProcessImageRejectsInvalidDimensions(t *testing.T) {
	err := ProcessImage("in.jpg", "out.jpg", ProcessOptions{
		Width:      0,
		Height:     100,
		ResizeMode: ResizeModeFit,
		Quality:    90,
	})
	if err == nil {
		t.Fatal("expected width validation error")
	}
	if !strings.Contains(err.Error(), "width must be greater than 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessImageRejectsInvalidQuality(t *testing.T) {
	err := ProcessImage("in.jpg", "out.jpg", ProcessOptions{
		Width:      100,
		Height:     100,
		ResizeMode: ResizeModeFit,
		Quality:    101,
	})
	if err == nil {
		t.Fatal("expected quality validation error")
	}
	if !strings.Contains(err.Error(), "quality must be between 1 and 100") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessImageRejectsUnsupportedEncoding(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	inputPath, err := saveTestImage(img, "jpg")
	if err != nil {
		t.Fatalf("failed to save test image: %v", err)
	}
	defer os.Remove(inputPath)

	err = ProcessImage(inputPath, filepath.Join(t.TempDir(), "out.jxl"), ProcessOptions{
		Width:        100,
		Height:       100,
		ResizeMode:   ResizeModeFit,
		Quality:      90,
		OutputFormat: "jxl",
	})
	if err == nil {
		t.Fatal("expected unsupported encoding error")
	}
	if !strings.Contains(err.Error(), "encoding to JXL format is not supported") {
		t.Fatalf("unexpected error: %v", err)
	}
}
