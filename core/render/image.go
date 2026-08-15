package render

import (
	"encoding/json"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// GetImagePixelWithIndex returns the color of a image at index.
func GetImagePixelWithIndex(i *software.Image, id int) *common.Color { return i.Colors[i.Mask[id]] }

// GetImageByName returns an loaded image in memory.
func GetImageByName(name string) *software.Image {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}
	return &software.Image{}
}

var is []*software.Image

func loadImages() error {
	files, err := loadFiles("images")
	if err != nil {
		return errors.Errorf("loading images files: %w", err)
	}

	for _, file := range files {
		i := &software.Image{}
		if err := json.Unmarshal(file.Data, i); err != nil {
			return errors.Errorf("unmarshaling image file %s: %w", file.Name, err)
		}
		is = append(is, i)
	}

	return nil
}
