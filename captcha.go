package gocaptcha

import (
	"bytes"
	"errors"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	dpi        = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	r          = rand.New(rand.NewSource(time.Now().UnixNano()))
	fontFamily = make([]string, 0)
)

const txtChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

const (
	
	ImageFormatPng ImageFormat = iota
	ImageFormatJpeg
	ImageFormatGif
)
const (

	CaptchaComplexLower CaptchaComplex = iota
	CaptchaComplexMedium
	CaptchaComplexHigh
)

type ImageFormat int

type CaptchaComplex int

type CaptchaImage struct {
	nrgba   *image.NRGBA
	width   int
	height  int
	Complex int
	Error   error
}

func NewCaptchaImage(width int, height int, bgColor color.RGBA) *CaptchaImage {
	return New(width, height, bgColor)
}


func New(width int, height int, bgColor color.RGBA) *CaptchaImage {
	m := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	return &CaptchaImage{
		nrgba:  m,
		height: height,
		width:  width,
	}
}

func (captcha *CaptchaImage) SaveImage(w io.Writer, imageFormat ImageFormat) error {

	if imageFormat == ImageFormatPng {
		return png.Encode(w, captcha.nrgba)
	}
	if imageFormat == ImageFormatJpeg {
		return jpeg.Encode(w, captcha.nrgba, &jpeg.Options{Quality: 100})
	}
	if imageFormat == ImageFormatGif {
		return gif.Encode(w, captcha.nrgba, &gif.Options{NumColors: 256})
	}

	return errors.New("not supported image format")
}

func (captcha *CaptchaImage) DrawHollowLine() *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	first := captcha.width / 20
	end := first * 19

	lineColor := color.RGBA{R: 245, G: 250, B: 251, A: 255}

	x1 := float64(r.Intn(first))

	x2 := float64(r.Intn(first) + end)

	multiple := float64(r.Intn(5)+3) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple = multiple * -1.0
	}

	w := captcha.height / 20

	for ; x1 < x2; x1++ {

		y := math.Sin(x1*math.Pi*multiple/float64(captcha.width)) * float64(captcha.height/3)

		if multiple < 0 {
			y = y + float64(captcha.height/2)
		}
		captcha.nrgba.Set(int(x1), int(y), lineColor)

		for i := 0; i <= w; i++ {
			captcha.nrgba.Set(int(x1), int(y)+i, lineColor)
		}
	}

	return captcha
}

func (captcha *CaptchaImage) DrawSineLine() *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	px := 0
	var py float64 = 0

	
	a := r.Intn(captcha.height / 2)

	
	b := Random(int64(-captcha.height/4), int64(captcha.height/4))


	f := Random(int64(-captcha.height/4), int64(captcha.height/4))
	
	var t float64 = 0
	if captcha.height > captcha.width/2 {
		t = Random(int64(captcha.width/2), int64(captcha.height))
	} else {
		t = Random(int64(captcha.height), int64(captcha.width/2))
	}
	w := (2 * math.Pi) / t

	px1 := 0
	px2 := int(Random(int64(float64(captcha.width)*0.8), int64(captcha.width)))

	c := color.RGBA{R: uint8(r.Intn(150)), G: uint8(r.Intn(150)), B: uint8(r.Intn(150)), A: uint8(255)}

	for px = px1; px < px2; px++ {
		if w != 0 {
			py = float64(a)*math.Sin(w*float64(px)+f) + b + (float64(captcha.width) / float64(5))
			i := captcha.height / 5
			for i > 0 {
				captcha.nrgba.Set(px+i, int(py), c)
				i--
			}
		}
	}

	return captcha
}

func (captcha *CaptchaImage) DrawLine(num int) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	first := captcha.width / 10
	end := first * 9

	y := captcha.height / 3

	for i := 0; i < num; i++ {

		point1 := Point{X: r.Intn(first), Y: r.Intn(y)}
		point2 := Point{X: r.Intn(first) + end, Y: r.Intn(y)}

		if i%2 == 0 {
			point1.Y = r.Intn(y) + y*2
			point2.Y = r.Intn(y)
		} else {
			point1.Y = r.Intn(y) + y*(i%2)
			point2.Y = r.Intn(y) + y*2
		}

		captcha.drawBeeline(point1, point2, randDeepColor())

	}
	return captcha
}

func (captcha *CaptchaImage) drawBeeline(point1 Point, point2 Point, lineColor color.RGBA) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	dx := math.Abs(float64(point1.X - point2.X))

	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy
	
	for {
		captcha.nrgba.Set(point1.X, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+2, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-2, point1.Y, lineColor)
		if point1.X == point2.X && point1.Y == point2.Y {
			return captcha
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

func (captcha *CaptchaImage) DrawBorder(borderColor color.RGBA) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	for x := 0; x < captcha.width; x++ {
		captcha.nrgba.Set(x, 0, borderColor)
		captcha.nrgba.Set(x, captcha.height-1, borderColor)
	}
	for y := 0; y < captcha.height; y++ {
		captcha.nrgba.Set(0, y, borderColor)
		captcha.nrgba.Set(captcha.width-1, y, borderColor)
	}
	return captcha
}

func (captcha *CaptchaImage) DrawNoise(complex CaptchaComplex) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	density := 18
	if complex == CaptchaComplexLower {
		density = 28
	} else if complex == CaptchaComplexMedium {
		density = 18
	} else if complex == CaptchaComplexHigh {
		density = 8
	}
	maxSize := (captcha.height * captcha.width) / density

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		captcha.nrgba.Set(rw, rh, randColor())
		size := r.Intn(maxSize)
		if size%3 == 0 {
			captcha.nrgba.Set(rw+1, rh+1, randColor())
		}
	}
	return captcha
}

func (captcha *CaptchaImage) DrawTextNoise(complex CaptchaComplex) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	density := 1500
	if complex == CaptchaComplexLower {
		density = 2000
	} else if complex == CaptchaComplexMedium {
		density = 1500
	} else if complex == CaptchaComplexHigh {
		density = 1000
	}

	maxSize := (captcha.height * captcha.width) / density

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)
	rawFontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(10))

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		text := RandText(1)
		fontSize := rawFontSize/2 + float64(r.Intn(5))

		c.SetSrc(image.NewUniform(RandLightColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			captcha.Error = err
			return captcha
		}
		c.SetFont(f)
		pt := freetype.Pt(rw, rh)

		_, err = c.DrawString(text, pt)
		if err != nil {
			captcha.Error = err
			return captcha
		}
	}
	return captcha
}

func (captcha *CaptchaImage) DrawText(text string) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)

	fontWidth := captcha.width / len(text)

	for i, s := range text {

		fontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(9))

		c.SetSrc(image.NewUniform(randDeepColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			captcha.Error = err
			return captcha
		}
		c.SetFont(f)

		x := (fontWidth)*i + (fontWidth)/int(fontSize)

		y := 5 + r.Intn(captcha.height/2) + int(fontSize/2)

		pt := freetype.Pt(x, y)

		_, err = c.DrawString(string(s), pt)
		if err != nil {
			captcha.Error = err
			return captcha
		}
	}
	return captcha

}

func (captcha *CaptchaImage) Save() (io io.Reader, err error) {
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, captcha.nrgba, nil)
	if err != nil {
		return
	}
	io = bytes.NewReader(buf.Bytes())
	return
}
