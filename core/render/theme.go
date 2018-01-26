package render

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// GetColorFromThemeByName returns an loaded theme's color in memory.
func GetColorFromThemeByName(themeName, colorName string) common.Color {
	for _, t := range ts {
		if t.Name == themeName {
			for k, c := range t.Colors {
				if k == colorName {
					return *c
				}
			}
		}
	}

	return common.Color{}
}

var ts []software.Theme

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

		var t software.Theme
		if err := json.Unmarshal(file, &t); err != nil {
			return errors.WithStack(err)
		}

		ts = append(ts, t)
	}

	return nil
}
