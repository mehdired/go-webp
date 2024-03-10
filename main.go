package main

import (
	"fmt"
	"go-webp/converter"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Image to WebP Converter"
	app.Usage = "Converts images or all images in a folder to WebP format with optional resizing"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "input, i",
			Usage: "Input image file path or folder path",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output folder path",
		},
		cli.IntFlag{
			Name:  "width, w",
			Usage: "Width for resizing (maintains aspect ratio)",
		},
		cli.IntFlag{
			Name:  "height, hght",
			Usage: "Height for resizing (maintains aspect ratio)",
		},
	}

	app.Action = func(c *cli.Context) error {
		inputs := c.StringSlice("input")
		outputFolder := c.String("output")
		resizeWidth := c.Int("width")
		resizeHeight := c.Int("height")

		if len(inputs) == 0 || outputFolder == "" {
			return cli.NewExitError("please specify both input and output paths.", 1)
		}

		// Create the output folder if it does not exist
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			if err := os.MkdirAll(outputFolder, 0755); err != nil {
				return cli.NewExitError(fmt.Sprintf("error creating output folder: %v", err), 1)
			}
		}

		for _, input := range inputs {
			// Check if the input is a file or a directory
			fileInfo, err := os.Stat(input)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("error accessing the input %s: %v", input, err), 1)
			}

			if fileInfo.IsDir() {
				// If it's a directory, convert all images in the directory
				err := converter.ConvertImagesInFolder(input, outputFolder, resizeWidth, resizeHeight)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("error converting images in the folder %s: %v", input, err), 1)
				}
			} else {
				// If it's a file, convert the individual image
				err := converter.ConvertToWebP(input, outputFolder, resizeWidth, resizeHeight)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("error converting the image %s: %v", input, err), 1)
				}
			}
		}

		fmt.Println("conversion completed successfully.")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
