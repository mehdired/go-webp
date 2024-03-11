package converter

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"

	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
)

// ConvertImagesInFolder converts all images in a folder to WebP format with optional resizing
func ConvertImagesInFolder(inputFolder, outputFolder string, width, height int) error {
	// Read the list of files in the input folder
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return fmt.Errorf("error reading the input folder: %v", err)
	}

	bar := progressbar.Default(int64(len(files)))

	// Iterate through all files in the folder
	for _, file := range files {
		// Ignore folders, only process image files
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if !isValidImageExtension(ext) {
			fmt.Printf("Skipping non-image file: %s\n", file.Name())
			continue
		}

		// Build the full path of the input file
		inputPath := filepath.Join(inputFolder, file.Name())

		// Convert the image to WebP
		err := ConvertToWebP(inputPath, outputFolder, width, height)
		if err != nil {
			fmt.Printf("error converting the image %s: %v\n", file.Name(), err)
		}

		bar.Add(1)
	}

	return nil
}

// ConvertToWebP converts an image to WebP format with optional resizing
func ConvertToWebP(inputPath, outputFolder string, width, height int) error {
	// Open the input image file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening the input image: %v", err)
	}
	defer inputFile.Close()

	// Decode the image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding the image: %v", err)
	}

	// Determine the orientation of the image (landscape or portrait)
	isLandscape := img.Bounds().Dy() < img.Bounds().Dx()

    // Calculate the new dimensions based on the orientation
    var newWidth, newHeight int

    if isLandscape {
		newWidth = width
		newHeight = int(float64(width) / float64(img.Bounds().Dx()) * float64(img.Bounds().Dy()))
	} else {
		newHeight = height
		newWidth = int(float64(height) / float64(img.Bounds().Dy()) * float64(img.Bounds().Dx()))
	}

	img = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	// Build the output path for the WebP file
	outputPath := filepath.Join(outputFolder, strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))+".webp")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating the WebP output file: %v", err)
	}
	defer outputFile.Close()

	// Encode the image to WebP format
	err = webp.Encode(outputFile, img, &webp.Options{Quality: 80})
	if err != nil {
		return fmt.Errorf("error encoding the image to WebP: %v", err)
	}

	return nil
}

// isValidImageExtension checks if the given file extension corresponds to a common image format
func isValidImageExtension(ext string) bool {
	// List of common image extensions
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true, // Include webp as well since we're converting to it
	}

	// Check if the extension exists in the map
	return imageExts[strings.ToLower(ext)]
}
