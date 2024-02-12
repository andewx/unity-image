package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andewx/img2sphere/tool"
)

const CS_GREEN = "\033[32m"
const CS_RESET = "\033[0m"
const CS_RED = "\033[31m"

func help() {
	fmt.Printf("%sUsage: img2sphere <input.png> <output.png> <options>\n", CS_RED)
	fmt.Printf("%sThis tool maps a 2D image to a disk hemisphere projection accepts PNG files only, default mode is a linear based projection scheme which should map correctly for spheres\n", CS_RESET)
	fmt.Printf("%sOptions:\n-h\tShow help\n-s\t<float> set scale mapping\n-l\tLog uv mapping calculations to file\n-c\t cubic mode\n-q\t quadratic mode\n-ln logarithmic mode\n-x exponential mode\n", CS_RESET)
}

func main() {
	fmt.Printf("%s%s%s\n", CS_GREEN, "img2sphere\n", CS_RESET)
	if len(os.Args) < 3 {
		help()
		os.Exit(1)
	}

	filename := os.Args[1]
	output := os.Args[2]
	var scale = 1.0
	var err error
	var mode = tool.MODE_LINEAR

	if len(os.Args) > 3 && os.Args[3] != "" {

		for i := 3; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-h":
				help()
				os.Exit(0)
			case "-s":
				if i+1 < len(os.Args) {
					scale, err = strconv.ParseFloat(os.Args[i+1], 64)
					if err != nil {
						fmt.Printf("%sError: %s\n", CS_RED, err)
						os.Exit(1)
					}
				}
			case "-l":
				tool.EnableLog()
			case "-c":
				mode = tool.MODE_CUBIC
			case "-q":
				mode = tool.MODE_SQUARE
			case "-ln":
				mode = tool.MODE_LOG
			case "-x":
				mode = tool.MODE_EXP

			}

		}
	}

	image, err := tool.OpenImage(filename)

	if err != nil {
		fmt.Printf("%sError: %s\n", CS_RED, err)
		os.Exit(1)
	}

	fmt.Printf("Mapping %s to disk hemisphere projection\n", filename)
	hemi := image.MapToHemisphere(scale, mode)
	err = hemi.Save(output)

	if err != nil {
		fmt.Printf("%sError: %s\n", CS_RED, err)
		os.Exit(1)
	}

	fmt.Printf("Saved to %s\n", output)

}
