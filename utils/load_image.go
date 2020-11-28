package utils

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
)

func LoadImage(filename string) (image.Image, error) {
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("could not load file")
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
