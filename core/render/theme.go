package render

import (
	"encoding/json"
	"image/color"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Theme contains a set of colors.
type Theme struct {
	Name   string                `json:"name"`
	Colors map[string]color.RGBA `json:"colors"`
}

// GetColorFromThemeByName returns an loaded theme's color in memory.
func GetColorFromThemeByName(themeName, colorName string) color.RGBA {
	for _, t := range ts {
		if t.Name == themeName {
			for k, c := range t.Colors {
				if k == colorName {
					return c
				}
			}
		}
	}

	return color.RGBA{}
}

var ts []Theme

func loadThemes() error {
	files, err := ioutil.ReadDir("./themes")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		file, err := ioutil.ReadFile("./themes/" + file.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		var t Theme
		if err := json.Unmarshal(file, &t); err != nil {
			return errors.WithStack(err)
		}

		ts = append(ts, t)
	}

	return nil
}
