package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andewx/imgTool/tool"
)

const CS_GREEN = "\033[32m"
const CS_RESET = "\033[0m"
const CS_RED = "\033[31m"

func help() {
	fmt.Printf("%sUsage: imgTool <options>\n", CS_RED)
	fmt.Printf("%sThis tool maps a 2D image to a disk hemisphere projection accepts PNG files only, default mode is a linear based projection scheme which should map correctly for spheres\n", CS_RESET)
	fmt.Printf("%sOptions:\n-h\tShow help\n-s\t<float> set scale mapping\n-l\tLog uv mapping calculations to file\n-c\t cubic mode\n-q\t quadratic mode\n-ln logarithmic mode\n-x exponential mode\n-umask unity pbr mask file from metallic (r), ambient (g), detail (b), smoothness (a)", CS_RESET)
}

func main() {
	fmt.Printf("%s%s%s\n", CS_GREEN, "img2sphere\n", CS_RESET)
	if len(os.Args) < 3 {
		help()
		os.Exit(1)
	}

	var filename string
	var output string
	var scale = 1.0
	var err error
	var mode = tool.MODE_LINEAR
	var metallic string
	var ambient string
	var detail string
	var smoothness string

	for i := 0; i < len(os.Args); i++ {
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
				i = i + 1
			}
		case "-l":
			tool.EnableLog()
		case "-c":
			mode = tool.MODE_CUBIC
			if i+2 < len(os.Args) {
				filename = os.Args[i+1]
				output = os.Args[i+2]
				i = i + 2
			}
		case "-q":
			mode = tool.MODE_SQUARE
			if i+2 < len(os.Args) {
				filename = os.Args[i+1]
				output = os.Args[i+2]
				i = i + 2
			}
		case "-ln":
			mode = tool.MODE_LOG
			if i+2 < len(os.Args) {
				filename = os.Args[i+1]
				output = os.Args[i+2]
				i = i + 2
			}
		case "-x":
			mode = tool.MODE_EXP
			if i+2 < len(os.Args) {
				filename = os.Args[i+1]
				output = os.Args[i+2]
				i = i + 2
			}
		case "-umask":
			mode = tool.MODE_UNITY_MASK
			if i+5 < len(os.Args) {
				metallic = os.Args[i+1]
				ambient = os.Args[i+2]
				detail = os.Args[i+3]
				smoothness = os.Args[i+4]
				output = os.Args[i+5]
				i = i + 5
			}
		}
	}

	if mode == tool.MODE_CUBIC || mode == tool.MODE_SQUARE || mode == tool.MODE_LOG || mode == tool.MODE_EXP {
		if filename == "" || output == "" {
			fmt.Printf("%sError: filename and output must be provided\n", CS_RED)
			os.Exit(1)
		} else {
			tool.MapToHemisphere(scale, mode, filename, output)
		}
	}

	if mode == tool.MODE_UNITY_MASK {
		if metallic == "" || ambient == "" || detail == "" || smoothness == "" || output == "" {
			fmt.Printf("%sError: metallic, ambient, detail, smoothness and output must be provided\n", CS_RED)
			os.Exit(1)
		} else {
			tool.CreateUnityDetailMask(metallic, ambient, detail, smoothness, output)
		}
	}

	fmt.Printf("Saved to %s\n", output)

}
