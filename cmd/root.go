package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"nim/pkg/image"
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
	Use:   "nim [input] [output]",
	Short: "Nim is an image manipulation tool",
	Long: `Nim is a cross-platform CLI tool for image manipulation.
It can resize, crop, pad, and convert images between formats.`,
	Example: `  nim -i input.jpg -o output.png -w 800 -H 600
  nim -i input.png -o output.jpg -s 1024x768 -q 90
  nim -i input.gif -o output.webp -s 300x300 -m stretch -p "#FF0000"
  nim input.jpg output.png -w 800 -H 600
  nim input.jpg output.png`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle positional arguments
		if len(args) > 2 {
			return fmt.Errorf("too many arguments: expected at most 2 arguments (input and output files)")
		} else if len(args) == 2 {
			// Two args: input and output
			inputFile = args[0]
			outputFile = args[1]
		} else if len(args) == 1 {
			// One arg: output only
			outputFile = args[0]
		}

		// Check if input and output files are provided
		if inputFile == "" {
			return fmt.Errorf("input file is required")
		}
		if outputFile == "" {
			return fmt.Errorf("output file is required")
		}

		// Parse size if provided
		if size != "" {
			parts := strings.Split(size, "x")
			if len(parts) != 2 {
				return fmt.Errorf("invalid size format: %s (expected WIDTHxHEIGHT)", size)
			}

			w, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("invalid width in size: %s", parts[0])
			}

			h, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid height in size: %s", parts[1])
			}

			width = w
			height = h
		}

		// Parse resize mode
		var mode image.ResizeMode
		switch strings.ToLower(resizeMode) {
		case "fit":
			mode = image.ResizeModeFit
		case "fill":
			mode = image.ResizeModeFill
		case "stretch":
			mode = image.ResizeModeStretch
		default:
			return fmt.Errorf("invalid resize mode: %s", resizeMode)
		}

		// Parse pad color
		var padColorRGB [3]uint8
		if padColor != "" {
			// Remove # if present
			padColor = strings.TrimPrefix(padColor, "#")

			// Parse hex color
			if len(padColor) == 6 {
				r, err := strconv.ParseUint(padColor[0:2], 16, 8)
				if err != nil {
					return fmt.Errorf("invalid pad color: %s", padColor)
				}
				g, err := strconv.ParseUint(padColor[2:4], 16, 8)
				if err != nil {
					return fmt.Errorf("invalid pad color: %s", padColor)
				}
				b, err := strconv.ParseUint(padColor[4:6], 16, 8)
				if err != nil {
					return fmt.Errorf("invalid pad color: %s", padColor)
				}
				padColorRGB = [3]uint8{uint8(r), uint8(g), uint8(b)}
			} else {
				return fmt.Errorf("invalid pad color format: %s (expected #RRGGBB)", padColor)
			}
		} else {
			// Default to white
			padColorRGB = [3]uint8{255, 255, 255}
		}

		// Create options
		options := image.ProcessOptions{
			Width:        width,
			Height:       height,
			ResizeMode:   mode,
			Quality:      quality,
			OutputFormat: outputFormat,
			PadColor:     padColorRGB,
		}

		// Process the image
		if err := image.ProcessImage(inputFile, outputFile, options); err != nil {
			return err
		}

		fmt.Printf("Image processed successfully: %s -> %s\n", inputFile, outputFile)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable the built-in help flag
	rootCmd.PersistentFlags().BoolP("help", "", false, "Help for nim")
	rootCmd.Flags().BoolP("help", "?", false, "Help for nim")

	// Define flags and bind them to variables
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file")
	rootCmd.Flags().IntVarP(&width, "width", "w", 800, "Target width")
	rootCmd.Flags().IntVarP(&height, "height", "H", 512, "Target height")
	rootCmd.Flags().StringVarP(&size, "size", "s", "", "Target size in format WIDTHxHEIGHT (e.g., 512x512)")
	rootCmd.Flags().StringVarP(&resizeMode, "mode", "m", "fit", "Resize mode (fit, fill, stretch)")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 85, "Output quality (1-100, only for JPEG)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Output format (jpg, png, gif, etc.)")
	rootCmd.Flags().StringVarP(&padColor, "pad-color", "p", "#FFFFFF", "Padding color in hex format (#RRGGBB)")
}
