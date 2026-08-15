package render

import (
	"encoding/json"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// GetColorFromLocalThemeByName returns an loaded theme's color in memory.
func GetColorFromLocalThemeByName(themeName, colorName string) *common.Color {
	for _, t := range ts {
		if t.Name == themeName {
			for k, c := range t.Colors {
				if k == colorName {
					return c
				}
			}
		}
	}
	return &common.Color{}
}

var ts []*software.Theme

func loadThemes() error {
	files, err := loadFiles("themes")
	if err != nil {
		return errors.Errorf("loading themes files: %w", err)
	}

	for _, file := range files {
		t := &software.Theme{}
		if err := json.Unmarshal(file.Data, t); err != nil {
			return errors.Errorf("unmarshaling theme file %s: %w", file.Name, err)
		}
		ts = append(ts, t)
	}

	return nil
}
