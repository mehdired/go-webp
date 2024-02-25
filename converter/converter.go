package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
)

// ConvertToWebP converts an image to WebP format with optional resizing
func ConvertToWebP(inputPath, outputFolder string, width int) error {
	// Read the image into memory
	img, err := bimg.Read(inputPath)
	if err != nil {
		return fmt.Errorf("error reading the image: %v", err)
	}

	// Resize the image with the specified width while maintaining the aspect ratio
	img, err = bimg.Resize(img, bimg.Options{Width: width, Embed: true})
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
func ConvertImages(inputs []string, outputFolder string, width int) error {
	for _, input := range inputs {
		// Check if the input is a file or a directory
		fileInfo, err := os.Stat(input)
		if err != nil {
			return fmt.Errorf("error accessing the input %s: %v", input, err)
		}

		if fileInfo.IsDir() {
			// If it's a directory, convert all images in the directory
			err := ConvertImagesInFolder(input, outputFolder, width)
			if err != nil {
				return fmt.Errorf("error converting images in the folder %s: %v", input, err)
			}
		} else {
			// If it's a file, convert the individual image
			err := ConvertToWebP(input, outputFolder, width)
			if err != nil {
				return fmt.Errorf("error converting the image %s: %v", input, err)
			}
		}
	}

	return nil
}

// ConvertImagesInFolder converts all images in a folder to WebP format with optional resizing
func ConvertImagesInFolder(inputFolder, outputFolder string, width int) error {
	// Read the list of files in the input folder
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return fmt.Errorf("error reading the input folder: %v", err)
	}

	// Iterate through all files in the folder
	for _, file := range files {
		// Ignore folders, only process image files
		if file.IsDir() {
			continue
		}

		// Build the full path of the input file
		inputPath := filepath.Join(inputFolder, file.Name())

		// Convert the image to WebP
		err := ConvertToWebP(inputPath, outputFolder, width)
		if err != nil {
			fmt.Printf("error converting the image %s: %v\n", file.Name(), err)
		}
	}

	return nil
}
