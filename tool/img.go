package tool

import (
	"image"
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
func (img *Image) MapToHemisphere(w float64, mode int) *Image {
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

	return &Image{Image: hemi, filename: img.filename}
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
