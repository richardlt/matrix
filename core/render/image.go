package render

import (
	"encoding/json"
	"image/color"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Image is a color matrix of pixels.
type Image struct {
	Name   string       `json:"name"`
	Height uint64       `json:"height"`
	Width  uint64       `json:"width"`
	Colors []color.RGBA `json:"colors"`
	Mask   []uint64     `json:"mask"`
}

// GetWithIndex returns the color of a image at index.
func (i Image) GetWithIndex(id int) color.RGBA {
	return i.Colors[i.Mask[id]]
}

// GetImageByName returns an loaded image in memory.
func GetImageByName(name string) Image {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}
	return Image{}
}

var is []Image

func loadImages() error {
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		file, err := ioutil.ReadFile("./images/" + file.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		var i Image
		if err := json.Unmarshal(file, &i); err != nil {
			return errors.WithStack(err)
		}

		is = append(is, i)
	}

	return nil
}
