package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	_ "golang.org/x/image/bmp"
)

func main() {
	var (
		bmpFile = flag.String("i", "test.bmp", "input bitmap file name")
	)
	flag.Parse()

	f, err := os.Open(*bmpFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
		return
	}

	of, err := os.Create("out.yuv")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer of.Close()
	w := bufio.NewWriter(of)

	buf := new(bytes.Buffer)
	yuv := make([]uint8, 4, 4)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x += 2 {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, _ := img.At(x+1, y).RGBA()
			//fmt.Println(y, x, r, g, b, a)
			y1, u, v := color.RGBToYCbCr(uint8(r1), uint8(g1), uint8(b1))
			y2, _, _ := color.RGBToYCbCr(uint8(r2), uint8(g2), uint8(b2))
			yuv[0] = u
			yuv[1] = y1
			yuv[2] = v
			yuv[3] = y2
			err = binary.Write(buf, binary.LittleEndian, yuv)
			if err != nil {
				log.Fatalln(err)
				return
			}
		}
	}
	fmt.Println(img.Bounds().Min.X, img.Bounds().Min.Y)
	fmt.Println(img.Bounds().Max.X, img.Bounds().Max.Y)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Fatalln(err)
		return
	}
	w.Flush()
}
