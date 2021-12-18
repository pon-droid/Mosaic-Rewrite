package main

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"path/filepath"
	"sort"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type RGB struct {
	R int
	G int
	B int
}

type fRGB struct {
	R    int
	G    int
	B    int
	file string
	dist float64
}

type src_list []fRGB

const (
	tile_size uint = 50
	out_size  uint = 250
)

var (
	src    []fRGB
	out    []RGB
	serial []string
)

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	serial = append(serial, s)

	return nil
}

func is_image(s string) bool {
	temp := len(s) - 1

	return s[temp] == 'g'
}

func avg_colour(mw *imagick.MagickWand, s string) {
	fmt.Println("Tile processing")
	var i uint
	var j uint

	var aR int
	var aG int
	var aB int

	var nR int = 0
	var nG int = 0
	var nB int = 0

	for i = 0; i < tile_size; i++ {
		for j = 0; j < tile_size; j++ {
			temp, err := mw.GetImagePixelColor(int(j), int(i))
			if err != nil {
				log.Fatal(err)
			}
			R := temp.GetRed() * 255
			G := temp.GetGreen() * 255
			B := temp.GetBlue() * 255

			nR += int(R)
			nG += int(G)
			nB += int(B)

		}
	}
	//Average
	aR = nR / int(tile_size*tile_size)
	aG = nG / int(tile_size*tile_size)
	aB = nB / int(tile_size*tile_size)

	src = append(src, fRGB{aR, aG, aB, s, 9999})

}

func colour(mw *imagick.MagickWand) {
	fmt.Println("Output processing")
	var i uint
	var j uint
	for i = 0; i < out_size; i++ {
		for j = 0; j < out_size; j++ {
			temp, err := mw.GetImagePixelColor(int(j), int(i))
			if err != nil {
				log.Fatal(err)
			}

			R := temp.GetRed() * 255
			G := temp.GetGreen() * 255
			B := temp.GetBlue() * 255

			nR := int(R)
			nG := int(G)
			nB := int(B)

			out = append(out, RGB{nR, nG, nB})
		}
	}
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	//	mw := imagick.NewMagickWand()

	ow := imagick.NewMagickWand()

	ow.ReadImage("/home/tgallaher/Downloads/IMG_2990.jpg")
	ow.ResizeImage(out_size, out_size, imagick.FILTER_LANCZOS)
	colour(ow)
	defer ow.Destroy()
	//ow.DisplayImage(os.Getenv("DISPLAY"))

	//mw.ReadImage("/home/tgallaher/Downloads/IMG_2990.jpg")

	//mw.DisplayImage(os.Getenv("DISPLAY"))
	filepath.WalkDir("/home/tgallaher/int-git/Recusive-Search/img_cells/", walk)

	for i := 0; i < len(serial); i++ {
		//println(serial[i])
		if is_image(serial[i]) {
			temp := imagick.NewMagickWand()
			temp.ReadImage(serial[i])
			avg_colour(temp, serial[i])
			defer temp.Destroy()
			//mw.DisplayImage(os.Getenv("DISPLAY"))
		}

	}

	sort_diff()

}

func calc_colour_diff(R, G, B, oR, oG, oB int) float64 {
	nR := float64((oR - R) * (oR - R))
	nG := float64((oG - G) * (oG - G))
	nB := float64((oB - B) * (oB - B))
	return math.Sqrt(nR + nG + nB)
}

//Template sort

func (e src_list) Len() int { return len(e) }

func (e src_list) Less(i, j int) bool { return e[i].dist < e[j].dist }

func (e src_list) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func sort_diff() {
	fmt.Println("Sorting...")

	mw := imagick.NewMagickWand()

	for i := range out {
		for j := range src {
			src[j].dist = calc_colour_diff(src[j].R, src[j].G, src[j].B, out[i].R, out[i].G, out[i].B)
		}

		sort.Sort(src_list(src))

		temp := imagick.NewMagickWand()
		temp.ReadImage(src[0].file)
		temp.ResizeImage(tile_size, tile_size, imagick.FILTER_LANCZOS)
		mw.AddImage(temp)
		defer temp.Destroy()

		fmt.Println("Iteration : ")
		fmt.Println(i)
	}
	fmt.Println("Writing...")
	end := imagick.NewDrawingWand()
	total_di := fmt.Sprintf("%vx%v+0+0", out_size, out_size)
	tile_di := fmt.Sprintf("%vx%v+0+0", tile_size, tile_size)

	montage := mw.MontageImage(end, total_di, tile_di, imagick.MONTAGE_MODE_CONCATENATE, "0x0+0+0")
	fmt.Println(len(src))
	montage.WriteImage("bigger8.png")

}

func sort_montage(mw *imagick.MagickWand) {
	mw.Clear()

	for i := range src {
		temp := imagick.NewMagickWand()
		temp.ReadImage(src[i].file)
		temp.ResizeImage(tile_size, tile_size, imagick.FILTER_LANCZOS)
		mw.AddImage(temp)
	}
	end := imagick.NewDrawingWand()
	total_di := fmt.Sprintf("%vx%v+0+0", out_size, out_size)
	tile_di := fmt.Sprintf("%vx%v+0+0", tile_size, tile_size)

	montage := mw.MontageImage(end, tile_di, total_di, imagick.MONTAGE_MODE_CONCATENATE, "0x0+0+0")
	fmt.Println(len(src))
	montage.WriteImage("o.png")
}
