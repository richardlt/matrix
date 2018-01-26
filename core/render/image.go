package render

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// GetImagePixelWithIndex returns the color of a image at index.
func GetImagePixelWithIndex(i software.Image, id int) common.Color { return *i.Colors[i.Mask[id]] }

// GetImageByName returns an loaded image in memory.
func GetImageByName(name string) software.Image {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}
	return software.Image{}
}

var is []software.Image

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

		var i software.Image
		if err := json.Unmarshal(file, &i); err != nil {
			return errors.WithStack(err)
		}

		is = append(is, i)
	}

	return nil
}
