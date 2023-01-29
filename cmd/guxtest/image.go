package main

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"unsafe"

	"github.com/leonsal/gux/gb"
	"github.com/leonsal/gux/window"
)

func init() {

	registerTest("image", 8, newTestImage)
}

type imageInfo struct {
	texID  gb.TextureID
	width  float32
	height float32
}

type testImage struct {
	img1     imageInfo
	img2     imageInfo
	img3     imageInfo
	imgScale float32
	delta    float64
}

func newTestImage(w *window.Window) ITest {

	t := new(testImage)
	var err error
	t.img1, err = createImageTexture(w, "assets/tux.png")
	if err != nil {
		log.Fatal(err)
	}
	t.img2, err = createImageTexture(w, "assets/tux.jpg")
	if err != nil {
		log.Fatal(err)
	}
	t.img3, err = createImageTexture(w, "assets/compression.jpg")
	if err != nil {
		log.Fatal(err)
	}
	t.imgScale = 1.0
	return t
}

func (t *testImage) draw(w *window.Window) {

	dl := w.DrawList()

	pmin := gb.Vec2{0, 0}
	pmax := gb.Vec2{t.img1.width, t.img1.height}
	w.AddImage(dl, t.img1.texID, pmin, pmax)

	pmin.X = pmax.X
	pmax.X += t.img2.width
	pmax.Y = t.img2.height
	w.AddImage(dl, t.img2.texID, pmin, pmax)

	pmin.X = pmax.X
	pmax.X += t.img3.width
	pmax.Y = t.img3.height
	w.AddImage(dl, t.img3.texID, pmin, pmax)

	pmin.X = 0
	pmin.Y = t.img1.height
	pmax.X = t.img1.width * t.imgScale
	pmax.Y = t.img1.height + t.img1.height*t.imgScale
	w.AddImage(dl, t.img1.texID, pmin, pmax)

	t.imgScale = float32(1.0+math.Sin(t.delta)) / 2
	t.delta += 0.01
}

func (t *testImage) destroy(w *window.Window) {

	w.DeleteTexture(t.img1.texID)
	w.DeleteTexture(t.img2.texID)
	w.DeleteTexture(t.img3.texID)
}

func createImageTexture(w *window.Window, path string) (imageInfo, error) {

	// Opens image file
	file, err := embedfs.Open(path)
	if err != nil {
		return imageInfo{}, err
	}
	defer file.Close()

	// Decodes the image from the specified reader
	img, _, err := image.Decode(file)
	if err != nil {
		return imageInfo{}, err
	}

	// Converts decoded image to RGBA
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	// Creates backend texture to store the image and transfer the image
	texID := w.CreateTexture(width, height, (*gb.RGBA)(unsafe.Pointer(&rgba.Pix[0])))
	return imageInfo{texID: texID, width: float32(width), height: float32(height)}, nil
}
