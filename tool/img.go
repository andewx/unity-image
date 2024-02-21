package tool

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type Image struct {
	Image    image.Image
	filename string
	log      *log.Logger
	file     *os.File
}

const PI = 3.14159265358979323846
const MODE_SQUARE = 0
const MODE_CUBIC = 1
const MODE_LINEAR = 2
const MODE_LOG = 3
const MODE_EXP = 4
const MODE_UNITY_MASK = 5
const MODE_FLIPBOOK = 6

var logEnabled = false

// Open PNG image and store in an image object
func OpenImage(path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	logFile, errLog := os.OpenFile("img2disk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logFile.Truncate(0)
	if errLog != nil {
		return nil, errLog
	}

	logger := log.New(logFile, "img2disk", log.LstdFlags)

	return &Image{Image: img, filename: path, log: logger, file: logFile}, nil
}

func EnableLog() {
	logEnabled = true
}

// Save image to file
func (img *Image) Save(output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()
	err = png.Encode(file, img.Image)
	if err != nil {
		return err
	}
	return nil
}

// Map 2D image to disk hemishphere projection , here we preserve the angle
func MapToHemisphere(w float64, mode int, file string, output string) {

	img, err := OpenImage(file)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	hemi := image.NewRGBA(image.Rect(0, 0, width, height))
	epsilon := 0.000001

	//We want (0,0) to represent the origin of the image while the domain is (-0.5,0.5)
	//We will map the pixel (i,j) to the point (x,y) in the domain via conformal mapping function
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			x := float64(i)/float64(width) - 0.5
			y := float64(j)/float64(height) - 0.5
			r_sq := math.Sqrt(x*x + y*y)

			if r_sq > 0.5 {
				continue
			}

			target_radius := 0.5

			r := target_radius * 2 * math.Atan(math.Sqrt(x*x+y*y)/2)
			theta := math.Atan2(y, x)

			if r_sq < epsilon {
				r = 0
				theta = 0
			}

			//Map the point (r,theta) to the point (u,v) in the disk hemisphere
			u := r * math.Cos(theta)
			v := r * math.Sin(theta)

			//Map the point (u,v) to the point (x',y') in the image
			xp := int((u + 0.5) * float64(width))
			yp := int((v + 0.5) * float64(height))

			//Map down the original x,y uv coordinates by magnitude scale factor
			mag := math.Sqrt(x*x + y*y)
			if mag > 0 {

				sc := 1.0
				if mode == MODE_SQUARE {
					sc = mag * mag
				} else if mode == MODE_CUBIC {
					sc = mag * mag * mag
				} else if mode == MODE_LINEAR {
					sc = mag
				} else if mode == MODE_LOG {
					sc = math.Log(mag + 1)
				} else if mode == MODE_EXP {
					sc = math.Exp(mag) - 1
				}

				x -= x * sc * w
				y -= y * sc * w
			}

			x = repeat(x, -0.5, 0.5)
			y = repeat(y, -0.5, 0.5)

			//Remap x,y to image coordinates
			i2 := int((x + 0.5) * float64(width))
			j2 := int((y + 0.5) * float64(height))

			pixel := img.Image.At(i2, j2)

			hemi.Set(xp, yp, pixel)

			//Log the mapping
			if logEnabled {
				img.log.Printf("Mapping (%d,%d) to (%d,%d) -> (%f,%f) to (%d,%d) -> (%f,%f)\n", i, j, xp, yp, x, y, i2, j2, u, v)
			}
		}
	}

	defer img.file.Close()

	out := &Image{Image: hemi, filename: img.filename}
	err = out.Save(output)

	if err != nil {
		log.Fatal(err)
	}

}

func clamp(x, min, max float64) float64 {
	if x < min {
		x = min
	}
	if x > max {
		x = max
	}
	return x
}

func repeat(x, min, max float64) float64 {
	return min + math.Mod(x-min, max-min)
}

func CreateFlipbookTextures(rows int, cols int, height int, width int, files []string, output string) {
	//Open all images provided
	images := make([]Image, len(files))
	for i, file := range files {
		img, err := OpenImage(file)
		if err != nil {
			log.Fatal(err)
		}
		/*
			if img.Image.Bounds().Max.X != width || img.Image.Bounds().Max.Y != height {
				log.Fatal("All images must have the same bounds")
			}
		*/
		images[i] = *img
	}

	//Create new combined image
	flipbook := image.NewRGBA(image.Rect(0, 0, width*cols, height*rows))

	//Iterate over all images and combine them into the flipbook
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			img := images[i*cols+j]
			bounds := img.Image.Bounds()
			for x := 0; x < bounds.Max.X; x++ {
				for y := 0; y < bounds.Max.Y; y++ {
					flipbook.Set(x+j*width, y+i*height, img.Image.At(x, y))
				}
			}
		}
	}

	img := Image{Image: flipbook, filename: "flipbook.png"}

	//Save the combined image
	err := img.Save(output)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUnityDetailMask(metallic_file string, ambient_file string, detail_file string, smoothness_file string, output_file string) {
	//Open all images provided
	metallic, err := OpenImage(metallic_file)
	if err != nil {
		log.Fatal(err)
	}
	ambient, err := OpenImage(ambient_file)
	if err != nil {
		log.Fatal(err)
	}
	detail, err := OpenImage(detail_file)
	if err != nil {
		log.Fatal(err)
	}
	smoothness, err := OpenImage(smoothness_file)
	if err != nil {
		log.Fatal(err)
	}

	//Check all images have same bounds
	if metallic.Image.Bounds() != ambient.Image.Bounds() || metallic.Image.Bounds() != detail.Image.Bounds() || metallic.Image.Bounds() != smoothness.Image.Bounds() {
		log.Fatal("All images must have the same bounds")
	}

	//Create new combined image
	bounds := metallic.Image.Bounds()

	unity_detail := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))

	//Iterate over all pixels and combine the images
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			//Get pixel luminance from each image
			metallic_pixel := Luminance(metallic.Image.At(i, j))
			ambient_pixel := Luminance(ambient.Image.At(i, j))
			detail_pixel := Luminance(detail.Image.At(i, j))
			smoothness_pixel := Luminance(smoothness.Image.At(i, j))

			//Create combined pixel from the luminance values
			pixel := color.RGBA{uint8(metallic_pixel * 255), uint8(ambient_pixel * 255), uint8(smoothness_pixel * 255), uint8(detail_pixel * 255)}

			//Combine the images
			unity_detail.Set(i, j, pixel)
		}
	}

	img := Image{Image: unity_detail, filename: output_file}

	//Save the combined image
	err = img.Save(output_file)
	if err != nil {
		log.Fatal(err)
	}

}

func Luminance(c color.Color) float64 {
	//Assume that the color is in grayscale space
	pixel := color.Gray16Model.Convert(c)

	//Get the luminance value
	l := float64(pixel.(color.Gray16).Y) / 65535.0

	return l
}
