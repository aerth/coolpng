package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var lport = os.Getenv("PORT")

func init() {
	if lport == "" {
		lport = ":8080"
	}
	if !strings.Contains(lport, ":") {
		lport = ":" + lport
	}
}
func drawpng(s, s2 string) (image.Image, error) {
	if s == "" {
		s = "hello png"
	}
	s = strings.Split(s, ":")[0]
	s2 = strings.Split(s2, " ")[0]
	// 100x100 image
	dest := image.NewRGBA(image.Rect(0, 0, 100, 100))

	draw.Draw(dest, dest.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	// grab font
	fontBytes, err := ioutil.ReadFile("TerminusTTF-4.40.1.ttf")
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// font options
	opts := &truetype.Options{}
	opts.DPI = 96
	opts.Size = 10
	opts.Hinting = font.HintingNone

	// write white text on the (already) black background
	d := font.Drawer{}
	d.Dst = dest
	d.Src = image.White
	d.Face = truetype.NewFace(f, opts)
	d.Dot = freetype.Pt(10, 45)
	d.DrawString(s)

	// move the typewriter for line 2
	d.Dot = freetype.Pt(10, 55)
	d.DrawString(s2)
	return dest, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	iimg, err := drawpng(r.RemoteAddr, r.UserAgent())
	if err != nil {
		println(err.Error())
		w.Write([]byte("error\n"))
		return
	}
	w.Header().Add("Content-Type", "image/png")
	encoder := png.Encoder{}
	encoder.CompressionLevel = png.DefaultCompression
	encoder.Encode(w, iimg)
	return
}
func main() {
	println("listening: ", lport)
	http.Handle("/image.png", http.HandlerFunc(handler))
	http.Handle("/cat.png", http.HandlerFunc(handler))
	err := http.ListenAndServe(lport, nil)
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
}
