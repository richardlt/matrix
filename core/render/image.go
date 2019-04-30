package render

import (
	"encoding/json"
	"fmt"

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
	files, err := loadFiles("images")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		var i software.Image
		if err := json.Unmarshal(file.Data, &i); err != nil {
			return fmt.Errorf("Can't unmarshal %s image file", file.Name)
		}
		is = append(is, i)
	}

	return nil
}
