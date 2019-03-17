package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"git.sr.ht/~sircmpwn/go-libvterm"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func drawChar(img *image.RGBA, x, y int, c color.Color, text string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

func main() {
	vt := vterm.New(25, 80)
	defer vt.Close()

	vt.SetUTF8(true)

	screen := vt.ObtainScreen()
	screen.Reset(true)
	state := vt.ObtainState()

	_, err := vt.Write([]byte("\033[31mHello \033[32mGolang\033[0m"))
	if err != nil {
		log.Fatal(err)
	}
	screen.Flush()

	rows, cols := vt.Size()
	img := image.NewRGBA(image.Rect(0, 0, cols*7, rows*13))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell, err := screen.GetCellAt(row, col)
			if err != nil {
				log.Fatal(err)
			}
			chars := cell.Chars()
			if len(chars) > 0 && chars[0] != 0 {
				drawChar(img, (col+1)*7, (row+1)*13, state.ConvertVTermColorToRGB(cell.Fg()), string(chars))
			}
		}
	}
	f, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
