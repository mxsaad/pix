package internal

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// ConvertFormat converts an image from one format to another.
// Supported formats: "png", "jpeg", "gif".
func ConvertFormat(inputPath, outputPath, outputFormat string) error {
	inputFormat := getFileFormat(inputPath)

	if inputFormat == outputFormat {
		return nil // No conversion needed if formats match
	}

	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Decode the input image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Encode and save the output image in the desired format
	switch outputFormat {
	case "png":
		return png.Encode(outputFile, img)
	case "jpeg", "jpg":
		return jpeg.Encode(outputFile, img, nil)
	case "gif":
		if err := gif.Encode(outputFile, img, nil); err != nil {
			return err
		}
	default:
		return errors.New("Unsupported output format: " + outputFormat)
	}

	return nil
}

// getFileFormat returns the file format (extension) of a given file path.
func getFileFormat(filePath string) string {
	parts := strings.Split(filePath, ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return ""
}
