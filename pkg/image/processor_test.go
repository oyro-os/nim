package image

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage creates a test image with the given dimensions and color
func createTestImage(width, height int, c color.RGBA) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
	return img, nil
}

// saveTestImage saves a test image to a temporary file
func saveTestImage(img *image.RGBA, format string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-image-*."+format)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Save the image in the specified format
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(tmpFile, img)
	default:
		return "", err
	}

	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestProcessImage(t *testing.T) {
	// Create a test image (100x100 red square)
	img, err := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Save the test image as a JPEG
	inputPath, err := saveTestImage(img, "jpg")
	if err != nil {
		t.Fatalf("Failed to save test image: %v", err)
	}
	defer os.Remove(inputPath) // Clean up

	// Create a temporary output file
	outputPath := filepath.Join(os.TempDir(), "output.png")
	defer os.Remove(outputPath) // Clean up

	// Test cases
	testCases := []struct {
		name    string
		options ProcessOptions
	}{
		{
			name: "Resize with fit mode",
			options: ProcessOptions{
				Width:        200,
				Height:       200,
				ResizeMode:   ResizeModeFit,
				Quality:      90,
				OutputFormat: "png",
				PadColor:     [3]uint8{255, 255, 255},
			},
		},
		{
			name: "Resize with fill mode",
			options: ProcessOptions{
				Width:        50,
				Height:       50,
				ResizeMode:   ResizeModeFill,
				Quality:      90,
				OutputFormat: "jpg",
				PadColor:     [3]uint8{0, 0, 0},
			},
		},
		{
			name: "Resize with stretch mode",
			options: ProcessOptions{
				Width:        150,
				Height:       75,
				ResizeMode:   ResizeModeStretch,
				Quality:      90,
				OutputFormat: "png",
				PadColor:     [3]uint8{0, 255, 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process the image
			err := ProcessImage(inputPath, outputPath, tc.options)
			if err != nil {
				t.Fatalf("ProcessImage failed: %v", err)
			}

			// Check if the output file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Fatalf("Output file was not created")
			}

			// Open the output file to verify it's a valid image
			_, err = os.Open(outputPath)
			if err != nil {
				t.Fatalf("Failed to open output file: %v", err)
			}
		})
	}
}