package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/i-rocky/pixr/pkg/image"
	"github.com/spf13/cobra"
)

var (
	inputFile    string
	outputFile   string
	width        int
	height       int
	size         string
	resizeMode   string
	quality      int
	outputFormat string
	padColor     string
)

var rootCmd = &cobra.Command{
	Use:   "pixr [input] [output]",
	Short: "Pixr is an image manipulation tool",
	Long: `Pixr is a cross-platform CLI tool for image manipulation.
It can resize, crop, pad, and convert images between formats.`,
	Example: `  pixr -i input.jpg -o output.png -w 800 -H 600
  pixr -i input.png -o output.jpg -s 1024x768 -q 90
  pixr -i input.gif -o output.webp -s 300x300 -m stretch -p "#FF0000"
  pixr input.jpg output.png -w 800 -H 600
  pixr input.jpg output.png`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvedInput, resolvedOutput, err := resolveInputOutput(args, inputFile, outputFile)
		if err != nil {
			return err
		}

		targetWidth, targetHeight, err := parseSizeArg(size, width, height)
		if err != nil {
			return err
		}
		if err := validateDimension("width", targetWidth); err != nil {
			return err
		}
		if err := validateDimension("height", targetHeight); err != nil {
			return err
		}
		if err := validateQuality(quality); err != nil {
			return err
		}

		mode, err := parseResizeMode(resizeMode)
		if err != nil {
			return err
		}

		padColorRGB, err := parsePadColor(padColor)
		if err != nil {
			return err
		}

		options := image.ProcessOptions{
			Width:        targetWidth,
			Height:       targetHeight,
			ResizeMode:   mode,
			Quality:      quality,
			OutputFormat: outputFormat,
			PadColor:     padColorRGB,
		}

		if err := image.ProcessImage(resolvedInput, resolvedOutput, options); err != nil {
			return err
		}

		fmt.Printf("Image processed successfully: %s -> %s\n", resolvedInput, resolvedOutput)
		return nil
	},
}

func resolveInputOutput(args []string, input, output string) (string, string, error) {
	switch len(args) {
	case 0:
	case 1:
		output = args[0]
	case 2:
		input = args[0]
		output = args[1]
	default:
		return "", "", fmt.Errorf("too many arguments: expected at most 2 arguments (input and output files)")
	}

	if input == "" {
		return "", "", fmt.Errorf("input file is required")
	}
	if output == "" {
		return "", "", fmt.Errorf("output file is required")
	}
	return input, output, nil
}

func parseSizeArg(sizeArg string, currentWidth, currentHeight int) (int, int, error) {
	normalized := strings.TrimSpace(strings.ToLower(sizeArg))
	if normalized == "" {
		return currentWidth, currentHeight, nil
	}

	wStr, hStr, ok := strings.Cut(normalized, "x")
	if !ok || strings.Contains(hStr, "x") {
		return 0, 0, fmt.Errorf("invalid size format: %s (expected WIDTHxHEIGHT)", sizeArg)
	}

	w, err := strconv.Atoi(strings.TrimSpace(wStr))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width in size: %s", wStr)
	}
	h, err := strconv.Atoi(strings.TrimSpace(hStr))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height in size: %s", hStr)
	}

	return w, h, nil
}

func parseResizeMode(modeArg string) (image.ResizeMode, error) {
	switch strings.ToLower(strings.TrimSpace(modeArg)) {
	case "fit":
		return image.ResizeModeFit, nil
	case "fill":
		return image.ResizeModeFill, nil
	case "stretch":
		return image.ResizeModeStretch, nil
	default:
		return "", fmt.Errorf("invalid resize mode: %s", modeArg)
	}
}

func parsePadColor(padColorArg string) ([3]uint8, error) {
	var zero [3]uint8
	normalized := strings.TrimSpace(strings.TrimPrefix(padColorArg, "#"))
	if normalized == "" {
		return [3]uint8{255, 255, 255}, nil
	}
	if len(normalized) != 6 {
		return zero, fmt.Errorf("invalid pad color format: %s (expected #RRGGBB)", padColorArg)
	}

	r, err := strconv.ParseUint(normalized[0:2], 16, 8)
	if err != nil {
		return zero, fmt.Errorf("invalid pad color: %s", padColorArg)
	}
	g, err := strconv.ParseUint(normalized[2:4], 16, 8)
	if err != nil {
		return zero, fmt.Errorf("invalid pad color: %s", padColorArg)
	}
	b, err := strconv.ParseUint(normalized[4:6], 16, 8)
	if err != nil {
		return zero, fmt.Errorf("invalid pad color: %s", padColorArg)
	}

	return [3]uint8{uint8(r), uint8(g), uint8(b)}, nil
}

func validateDimension(name string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0", name)
	}
	return nil
}

func validateQuality(value int) error {
	if value < 1 || value > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable the built-in help flag
	rootCmd.PersistentFlags().BoolP("help", "", false, "Help for pixr")
	rootCmd.Flags().BoolP("help", "?", false, "Help for pixr")

	// Define flags and bind them to variables
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file")
	rootCmd.Flags().IntVarP(&width, "width", "w", 800, "Target width")
	rootCmd.Flags().IntVarP(&height, "height", "H", 512, "Target height")
	rootCmd.Flags().StringVarP(&size, "size", "s", "", "Target size in format WIDTHxHEIGHT (e.g., 512x512)")
	rootCmd.Flags().StringVarP(&resizeMode, "mode", "m", "fit", "Resize mode (fit, fill, stretch)")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 85, "Output quality (1-100)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Output format (jpg, png, gif, etc.)")
	rootCmd.Flags().StringVarP(&padColor, "pad-color", "p", "#FFFFFF", "Padding color in hex format (#RRGGBB)")
}
