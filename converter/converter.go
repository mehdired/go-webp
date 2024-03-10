package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/schollz/progressbar/v3"
)

// ConvertToWebP converts an image to WebP format with optional resizing
func ConvertToWebP(inputPath, outputFolder string, width int, height int) error {
	// Read the image into memory
	img, err := bimg.Read(inputPath)
	if err != nil {
		return fmt.Errorf("error reading the image: %v", err)
	}

	size, _ := bimg.Size(img)
	isLandscape := size.Width > size.Height

	// Resize the image with the specified width or height based on landscape/portrait
	var imgOptions bimg.Options
	if isLandscape {
		imgOptions = bimg.Options{Width: width, Height: 0, Embed: true}
	} else {
		imgOptions = bimg.Options{Width: 0, Height: height, Embed: true}
	}

	// Resize the image with the specified width while maintaining the aspect ratio
	img, err = bimg.Resize(img, imgOptions)
	if err != nil {
		return fmt.Errorf("error resizing the image: %v", err)
	}

	// Build the full path for the WebP output file
	outputFileName := fmt.Sprintf("%s.webp", strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)))
	outputPath := filepath.Join(outputFolder, outputFileName)

	// Convert the image to WebP format
	webpImg, err := bimg.NewImage(img).Convert(bimg.WEBP)
	if err != nil {
		return fmt.Errorf("error converting the image to WebP: %v", err)
	}

	// Create the WebP output file in the specified output folder
	err = bimg.Write(outputPath, webpImg)
	if err != nil {
		return fmt.Errorf("error writing the WebP output file: %v", err)
	}

	return nil
}

// ConvertImages converts specified images or all images in a folder to WebP format with optional resizing
func ConvertImages(inputs []string, outputFolder string, width int, height int) error {
	for _, input := range inputs {
		// Check if the input is a file or a directory
		fileInfo, err := os.Stat(input)
		if err != nil {
			return fmt.Errorf("error accessing the input %s: %v", input, err)
		}

		if fileInfo.IsDir() {
			// If it's a directory, convert all images in the directory
			err := ConvertImagesInFolder(input, outputFolder, width, height)
			if err != nil {
				return fmt.Errorf("error converting images in the folder %s: %v", input, err)
			}
		} else {
			// If it's a file, convert the individual image
			err := ConvertToWebP(input, outputFolder, width, height)
			if err != nil {
				return fmt.Errorf("error converting the image %s: %v", input, err)
			}
		}
	}

	return nil
}

// ConvertImagesInFolder converts all images in a folder to WebP format with optional resizing
func ConvertImagesInFolder(inputFolder, outputFolder string, width int, height int) error {
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

		// Check if the file has a valid image extension
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
