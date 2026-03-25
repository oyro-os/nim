# Pixr - Image Manipulation Tool

Pixr is a cross-platform CLI tool for image manipulation written in Go. It can resize, crop, pad, and convert images between formats.

## Features

- Resize images with different modes (fit, fill, stretch)
- Convert between common image formats (JPEG, PNG, GIF)
- Adjust output quality for lossy output formats
- Customize padding color
- Cross-platform support

## Installation

### Scoop
```powershell
scoop bucket add rocky https://github.com/i-rocky/bucket
scoop install nim
```

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/i-rocky/pixr.git
   cd pixr
   ```

2. Build the binary:
   ```
   go build -o pixr
   ```

3. (Optional) Move the binary to a directory in your PATH:
   ```
   # Linux/macOS
   sudo mv pixr /usr/local/bin/

   # Windows
   # Move pixr.exe to a directory in your PATH
   ```

### Go Install

```
go install github.com/i-rocky/pixr@latest
```

## Usage

Basic usage:

```
pixr -i input.jpg -o output.png -w 800 -H 600
```

or with combined size:

```
pixr -i input.jpg -o output.png -s 800x600
```

You can also use positional arguments:

```
pixr input.jpg output.png -w 800 -H 600
```

or even simpler:

```
pixr input.jpg output.png
```

### Options

- `--input`, `-i`: Input image file (required if not provided as positional argument)
- `--output`, `-o`: Output image file (required if not provided as positional argument)

Note: You can also provide input and output files as positional arguments:
- If you provide one positional argument, it will be used as the output file.
- If you provide two positional arguments, they will be used as input and output files respectively.
- If you provide more than two positional arguments, an error will be shown.
- `--width`, `-w`: Target width (default: 800)
- `--height`, `-H`: Target height (default: 512)
- `--size`, `-s`: Target size in format WIDTHxHEIGHT (e.g., 512x512)
- `--mode`, `-m`: Resize mode (fit, fill, stretch) (default: fit)
  - `fit`: Resize the image to fit within the specified dimensions while maintaining aspect ratio
  - `fill`: Resize the image to fill the specified dimensions while maintaining aspect ratio and crops any excess
  - `stretch`: Resize the image to the specified dimensions without maintaining aspect ratio
- `--quality`, `-q`: Output quality (1-100) (default: 85)
- `--format`, `-f`: Output format (jpg, png, gif, etc.) (default: determined from output filename)
- `--pad-color`, `-p`: Padding color in hex format (#RRGGBB) (default: #FFFFFF)

### Examples

Resize an image to fit within 800x600 pixels:
```
pixr -i input.jpg -o output.jpg -w 800 -H 600 -m fit
```

Resize and crop an image to exactly 300x300 pixels:
```
pixr -i input.png -o output.jpg -s 300x300 -m fill
```

Resize an image to 1024x768 pixels without maintaining aspect ratio:
```
pixr -i input.gif -o output.png -s 1024x768 -m stretch
```

Convert an image to JPEG with 90% quality:
```
pixr -i input.png -o output.jpg -q 90
```

Resize an image with a red background for padding:
```
pixr -i input.png -o output.png -w 800 -H 600 -m fit -p "#FF0000"
```

## Supported Image Formats

### Fully Supported (Read and Write)
- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- TIFF (.tiff, .tif)
- WebP (.webp)
- AVIF (.avif)
- ICO (.ico)
- ICNS (.icns)

### Partially Supported (Read Only)
- HEIC/HEIF (.heic, .heif)
- JPEG XL (.jxl)

Note: For partially supported formats, Pixr can read the source image, but you must write to a supported output format (for example PNG/JPEG/WebP). This is because:
- The goheif library (github.com/jdeng/goheif) only supports decoding HEIC/HEIF images, not encoding them
- The jxl-go library (github.com/kpfaulkner/jxl-go) only supports decoding JXL images, not encoding them

### Not Supported
- JPEG 2000 (.jp2) for both decoding and encoding

If you need to write to partially supported formats or work with unsupported formats, use another tool after processing with Pixr.

## License

MIT
