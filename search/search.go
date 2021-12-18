package main

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var serial []string

func is_image(s string) bool {
	temp := len(s) - 1

	return s[temp] == 'g' || s[temp] == 'G'
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	serial = append(serial, s)

	return nil
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	filepath.WalkDir("/home/tgallaher/Downloads/athletics carnival photos/Athletics carnival Ollie/", walk)

	for i := 0; i < len(serial); i++ {
		if is_image(serial[i]) {
			mw.ReadImage(serial[i])
			mw.ResizeImage(50, 50, imagick.FILTER_LANCZOS)
			var name string = "output"
			name += fmt.Sprint(i)
			name += ".png"
			mw.WriteImage(name)
		}
	}
}
