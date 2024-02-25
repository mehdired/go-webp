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
	}

	app.Action = func(c *cli.Context) error {
		inputs := c.StringSlice("input")
		outputFolder := c.String("output")
		resizeWidth := c.Int("width")

		if len(inputs) == 0 || outputFolder == "" {
			return cli.NewExitError("please specify both input and output paths.", 1)
		}

		err := converter.ConvertImages(inputs, outputFolder, resizeWidth)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("error converting images: %v", err), 1)
		}

		fmt.Println("conversion completed successfully.")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
