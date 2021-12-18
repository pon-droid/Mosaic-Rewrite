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

type magical struct {
	image *imagick.MagickWand
	order int
}

type src_list []fRGB

type magick_list []magical

const (
	tile_size uint = 25
	out_size  uint = 175
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

	sort_man()

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

func (e magick_list) Len() int { return len(e) }

func (e magick_list) Less(i, j int) bool { return e[i].order < e[j].order }

func (e magick_list) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func sort_man() {
	fmt.Println("Sorting...")

	mw := imagick.NewMagickWand()

	sum := (len(out) / 5)
	var channels []chan magical
	buffer := []magical{}

	/*
		output1 := make(chan magical)
		output2 := make(chan magical)
		output3 := make(chan magical)
		output4 := make(chan magical)
		output5 := make(chan magical)

		go sort_iter(sum*4, sum*(4+1), output5, 5)
		go sort_iter(sum*3, sum*(3+1), output4, 4)
		go sort_iter(sum*2, sum*(2+1), output3, 3)
		go sort_iter(sum*1, sum*(1+1), output2, 2)
		go sort_iter(sum*0, sum*(0+1), output1, 1)

		buffer := []magical{}

		for i := 0; i < 5; i++ {
			select {
			case t1 := <-output1:
				buffer = append(buffer, t1)
			case t2 := <-output2:
				buffer = append(buffer, t2)
			case t3 := <-output3:
				buffer = append(buffer, t3)
			case t4 := <-output4:
				buffer = append(buffer, t4)
			case t5 := <-output5:
				buffer = append(buffer, t5)
			}
		}*/

	for i := 0; i < 5; i++ {
		output := make(chan magical)
		channels = append(channels, output)
		go sort_iter(sum*i, sum*(i+1), output, i)
	}

	for i := 0; i < 5; i++ {
		temp := <-channels[i]

		buffer = append(buffer, temp)

	}

	sort.Sort(magick_list(buffer))

	for i := 0; i < 5; i++ {
		mw.AddImage(buffer[i].image)
		defer buffer[i].image.Destroy()
	}

	/*
		for i := 0; i < 5; i++ {
			temp := sort_iter(sum*i, sum*(i+1))
			mw.AddImage(temp)
			defer temp.Destroy()
		}
	*/
	fmt.Println("Writing...")
	end := imagick.NewDrawingWand()
	total_di := fmt.Sprintf("%vx%v+0+0", out_size, out_size)
	tile_di := fmt.Sprintf("%vx%v+0+0", tile_size, tile_size)
	fmt.Println(total_di)
	fmt.Println(tile_di)
	montage := mw.MontageImage(end, total_di, tile_di, imagick.MONTAGE_MODE_CONCATENATE, "0x0+0+0")

	montage.WriteImage("go_manager.jpg")
}

func sort_iter(start int, end int, ch chan magical, num int) {
	mw := imagick.NewMagickWand()
	var temp_src []fRGB
	for i := range src {
		temp_src = append(temp_src, src[i])
	}
	for i := start; i < end; i++ {
		for j := range temp_src {
			temp_src[j].dist = calc_colour_diff(temp_src[j].R, temp_src[j].G, temp_src[j].B, out[i].R, out[i].G, out[i].B)
		}
		sort.Sort(src_list(temp_src))

		temp := imagick.NewMagickWand()
		temp.ReadImage(temp_src[0].file)
		temp.ResizeImage(tile_size, tile_size, imagick.FILTER_LANCZOS)
		mw.AddImage(temp)
		defer temp.Destroy()
	}
	ch <- magical{mw, num}
}

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
